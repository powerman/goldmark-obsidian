package obsidian

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/powerman/goldmark-obsidian/ast"
)

// Default attributes used by [NewPlugTasksParser].
const (
	DefaultPlugTasksListClass               = "contains-task-list"
	DefaultPlugTasksListItemNotCheckedClass = "task-list-item"
	DefaultPlugTasksListItemCheckedClass    = "task-list-item is-checked"
	DefaultPlugTasksListItemStatusAttr      = "data-task"
	DefaultPlugTasksCheckboxClass           = "task-list-item-checkbox"
)

// Task status symbols used by [NewPlugTasksParser] by default.
// Any other (unknown) symbols will have status TODO.
//
// NOTE: Upper case 'X' is unknown (i.e. has TODO status) in default Tasks configuration.
const (
	DefaultPlugTasksTODOSymbol       = ' '
	DefaultPlugTasksInProgressSymbol = '/'
	DefaultPlugTasksDoneSymbol       = 'x'
	DefaultPlugTasksCancelledSymbol  = '-'
)

type plugTasksConfig struct {
	StatusType              map[rune]ast.PlugTasksStatusType
	ListClass               []byte
	ListItemNotCheckedClass []byte
	ListItemCheckedClass    []byte
	ListItemStatusAttr      []byte
	CheckboxClass           []byte
}

func newPlugTasksConfig() plugTasksConfig {
	return plugTasksConfig{
		StatusType: map[rune]ast.PlugTasksStatusType{
			DefaultPlugTasksTODOSymbol:       ast.PlugTasksStatusTypeTODO,
			DefaultPlugTasksInProgressSymbol: ast.PlugTasksStatusTypeInProgress,
			DefaultPlugTasksDoneSymbol:       ast.PlugTasksStatusTypeDone,
			DefaultPlugTasksCancelledSymbol:  ast.PlugTasksStatusTypeCancelled,
		},
		ListClass:               []byte(DefaultPlugTasksListClass),
		ListItemNotCheckedClass: []byte(DefaultPlugTasksListItemNotCheckedClass),
		ListItemCheckedClass:    []byte(DefaultPlugTasksListItemCheckedClass),
		ListItemStatusAttr:      []byte(DefaultPlugTasksListItemStatusAttr),
		CheckboxClass:           []byte(DefaultPlugTasksCheckboxClass),
	}
}

// A PlugTasksOption configures a [PlugTasksParser].
type PlugTasksOption func(*plugTasksConfig)

// WithPlugTasksStatusType sets a statusType for a symbol.
// If you need to set many symbols then using [WithPlugTasksStatusTypes] may be more
// convenient than multiple calls to WithPlugTasksStatusType.
//
// To "disable" one of default symbols set it statusType to ast.PlugTasksStatusTypeTODO,
// which is used by default for all unknown symbols.
func WithPlugTasksStatusType(symbol rune, statusType ast.PlugTasksStatusType) PlugTasksOption {
	return func(config *plugTasksConfig) {
		config.StatusType[symbol] = statusType
	}
}

// WithPlugTasksStatusTypes sets a statusTypes for a symbols.
//
// To "disable" one of default symbols set it statusType to ast.PlugTasksStatusTypeTODO,
// which is used by default for all unknown symbols.
func WithPlugTasksStatusTypes(statusTypes map[rune]ast.PlugTasksStatusType) PlugTasksOption {
	return func(config *plugTasksConfig) {
		for symbol, statusType := range statusTypes {
			config.StatusType[symbol] = statusType
		}
	}
}

// WithPlugTasksListClass sets a class for a list containing task(s).
func WithPlugTasksListClass[T []byte | string](class T) PlugTasksOption {
	return func(config *plugTasksConfig) {
		config.ListClass = []byte(class)
	}
}

// WithPlugTasksListItemNotCheckedClass sets a class for a list item containing
// a task without status symbol: [ ].
func WithPlugTasksListItemNotCheckedClass[T []byte | string](class T) PlugTasksOption {
	return func(config *plugTasksConfig) {
		config.ListItemNotCheckedClass = []byte(class)
	}
}

// WithPlugTasksListItemCheckedClass sets a class for a list item containing
// a task with any status symbol (i.e. not a space).
func WithPlugTasksListItemCheckedClass[T []byte | string](class T) PlugTasksOption {
	return func(config *plugTasksConfig) {
		config.ListItemCheckedClass = []byte(class)
	}
}

// WithPlugTasksListItemStatusAttr sets a name for a list item attribute which will contain
// a task status symbol (one within [ ]).
// Attribute will be set to empty string if status symbol is a space.
func WithPlugTasksListItemStatusAttr[T []byte | string](name T) PlugTasksOption {
	return func(config *plugTasksConfig) {
		config.ListItemStatusAttr = []byte(name)
	}
}

// WithPlugTasksCheckboxClass sets a class for an <input type=checkbox>.
func WithPlugTasksCheckboxClass[T []byte | string](class T) PlugTasksOption {
	return func(config *plugTasksConfig) {
		config.CheckboxClass = []byte(class)
	}
}

// PlugTasksParser is an Obsidian plugin Tasks's task status parser.
//
// This parser must take precedence over the parser.LinkParser.
type PlugTasksParser struct {
	plugTasksConfig
}

// NewPlugTasksParser returns a new PlugTasksParser.
func NewPlugTasksParser(opts ...PlugTasksOption) *PlugTasksParser {
	s := &PlugTasksParser{
		plugTasksConfig: newPlugTasksConfig(),
	}
	for _, opt := range opts {
		opt(&s.plugTasksConfig)
	}
	return s
}

// Trigger implements [parser.InlineParser].
func (*PlugTasksParser) Trigger() []byte {
	return []byte{'['}
}

var plugTasksStatusRegexp = regexp.MustCompile(`^\[(.)\]\s*`)

// Parse implements [parser.InlineParser].
func (p *PlugTasksParser) Parse(parent gast.Node, block text.Reader, _ parser.Context) gast.Node {
	// Given AST structure must be like
	// - List
	//   - ListItem			: parent.Parent
	//     - TextBlock|Paragraph    : parent
	//       (current line)
	if parent.Parent() == nil || parent.Parent().FirstChild() != parent {
		return nil
	}
	if parent.HasChildren() {
		return nil
	}
	if _, ok := parent.Parent().(*gast.ListItem); !ok {
		return nil
	}

	line, _ := block.PeekLine()
	m := plugTasksStatusRegexp.FindSubmatchIndex(line)
	if m == nil {
		return nil
	}
	symbol, _ := utf8.DecodeRune(line[m[2]:m[3]])
	block.Advance(m[1])

	statusType, ok := p.StatusType[symbol]
	if !ok { // https://publish.obsidian.md/tasks/Getting+Started/Statuses#Unknown+Statuses
		statusType = ast.PlugTasksStatusTypeTODO
	}

	node := ast.NewPlugTasksStatus(symbol, statusType)
	itemNode := parent.Parent()
	listNode := itemNode.Parent()

	if symbol == ' ' {
		itemNode.SetAttribute(p.ListItemStatusAttr, []byte{})
	} else {
		itemNode.SetAttribute(p.ListItemStatusAttr, []byte(string(symbol)))
	}

	appendClass(node, p.CheckboxClass)
	if node.IsChecked() {
		appendClass(itemNode, p.ListItemCheckedClass)
	} else {
		appendClass(itemNode, p.ListItemNotCheckedClass)
	}
	appendClass(listNode, p.ListClass)

	return node
}

const (
	// Plugin Tasks has own definition for Obsidian tags, which differs from Obsidian one:
	// - it does not require at least one non-numeric symbol
	// - it accepts extra symbols: '+;=[\]`~
	plugTasksTagRe     = `#[-a-zA-Z0-9/_'+;=[\\\]` + "`" + `~\x{000080}-\x{10FFFF}]+(?:\s+|$)`
	plugTasksBlockIDRe = `\^[A-Za-z0-9-]+$`
	// Plugin Tasks parses any numbers in this format, but after parsing marks wrong
	// dates as "invalid".
	plugTasksDateRe = `\d\d\d\d-\d\d-\d\d`
	// Plugin Tasks parses any phrase where words consists of A-Za-z0-9!, and then
	// try to parse it as an instruction like "every week on Monday, Thuesday".
	plugTasksRecurringRe = `[A-Za-z0-9!,](?:[A-Za-z0-9!, ]*[A-Za-z0-9!,])?`
	plugTasksIDRe        = `[A-Za-z0-9_-]+`
	plugTasksIDsRe       = plugTasksIDRe + `(?:,` + plugTasksIDRe + `)*`
	// Plugin Tasks parses any word which consists only from letters, and then
	// try to parse it (case-insensitive) as one of actions "keep" or "delete".
	plugTasksOnCompletionActionRe = `[A-Za-z]+`
	plugTasksPrioLowestEmoji      = `⏬`
	plugTasksPrioLowEmoji         = `🔽`
	plugTasksPrioMediumEmoji      = `🔼`
	plugTasksPrioHighEmoji        = `⏫`
	plugTasksPrioHighestEmoji     = `🔺`
	plugTasksIDEmoji              = `🆔`
	plugTasksDependsOnEmoji       = `⛔`
	plugTasksDueEmoji             = `📅`
	plugTasksScheduledEmoji       = `⏳`
	plugTasksStartEmoji           = `🛫`
	plugTasksCreatedEmoji         = `➕`
	plugTasksDoneEmoji            = `✅`
	plugTasksCancelledEmoji       = `❌`
	plugTasksRecurringEmoji       = `🔁`
	plugTasksOnCompletionEmoji    = `🏁`
)

var plugTasksPropRegexp = regexp.MustCompile(`^\s*(` + // Group 1 is for block.Advance().
	// Starts with one of known emoji.
	`([` + // Group 2 is an emoji without any args.
	plugTasksPrioLowestEmoji +
	plugTasksPrioLowEmoji +
	plugTasksPrioMediumEmoji +
	plugTasksPrioHighEmoji +
	plugTasksPrioHighestEmoji +
	`])\s*|` +
	`([` + // Group 3 is an emoji with ID arg.
	plugTasksIDEmoji +
	`])\s*(` + plugTasksIDRe + `)\s*|` + // Group 4 is an ID.
	`([` + // Group 5 is an emoji with ID list arg.
	plugTasksDependsOnEmoji +
	`])\s*(` + plugTasksIDsRe + `)\s*|` + // Group 6 is an ID list.
	`([` + // Group 7 is an emoji with date arg.
	plugTasksDueEmoji +
	plugTasksScheduledEmoji +
	plugTasksStartEmoji +
	plugTasksCreatedEmoji +
	plugTasksDoneEmoji +
	plugTasksCancelledEmoji +
	`])\s*(` + plugTasksDateRe + `)\s*|` + // Group 8 is a date.
	`([` + // Group 9 is an emoji with recurring schedule arg.
	plugTasksRecurringEmoji +
	`])\s*(` + plugTasksRecurringRe + `)\s*|` + // Group 10 is a recurring schedule.
	`([` + // Group 11 is an emoji with an action arg.
	plugTasksOnCompletionEmoji +
	`])\s*(` + plugTasksOnCompletionActionRe + `)\s*` + // Group 12 is an action.
	`)` + // End of group 1.
	// Followed by any amount of same emoji or tag.
	`(?:` +
	plugTasksTagRe + `\s*|` +
	plugTasksPrioLowestEmoji + `\s*|` +
	plugTasksPrioLowEmoji + `\s*|` +
	plugTasksPrioMediumEmoji + `\s*|` +
	plugTasksPrioHighEmoji + `\s*|` +
	plugTasksPrioHighestEmoji + `\s*|` +
	plugTasksIDEmoji + `\s*` + plugTasksIDRe + `\s*|` +
	plugTasksDependsOnEmoji + `\s*` + plugTasksIDsRe + `\s*|` +
	plugTasksDueEmoji + `\s*` + plugTasksDateRe + `\s*|` +
	plugTasksScheduledEmoji + `\s*` + plugTasksDateRe + `\s*|` +
	plugTasksStartEmoji + `\s*` + plugTasksDateRe + `\s*|` +
	plugTasksCreatedEmoji + `\s*` + plugTasksDateRe + `\s*|` +
	plugTasksDoneEmoji + `\s*` + plugTasksDateRe + `\s*|` +
	plugTasksCancelledEmoji + `\s*` + plugTasksDateRe + `\s*|` +
	plugTasksRecurringEmoji + `\s*` + plugTasksRecurringRe + `\s*|` +
	plugTasksOnCompletionEmoji + `\s*` + plugTasksOnCompletionActionRe + `\s*` +
	`)*` +
	// Followed by optional block id and end of line.
	`(?:` + plugTasksBlockIDRe + `)?$`,
)

var plugTasksPropTrigger = func(emojis ...string) map[byte]bool { // Const.
	firstByte := make(map[byte]bool)
	for _, s := range emojis {
		firstByte[s[0]] = true
	}
	return firstByte
}(
	plugTasksIDEmoji,
	plugTasksPrioLowestEmoji,
	plugTasksPrioLowEmoji,
	plugTasksPrioMediumEmoji,
	plugTasksPrioHighEmoji,
	plugTasksPrioHighestEmoji,
	plugTasksRecurringEmoji,
	plugTasksDueEmoji,
	plugTasksScheduledEmoji,
	plugTasksStartEmoji,
	plugTasksDependsOnEmoji,
	plugTasksCreatedEmoji,
	plugTasksDoneEmoji,
	plugTasksCancelledEmoji,
	plugTasksOnCompletionEmoji,
)

// PlugTasksPropParser is an Obsidian plugin Tasks's task properties parser.
//
// Current implementation supports only properties defined in emoji style.
//
// This parser must have lowest precedence than all other inline parsers (e.g. 1000).
type PlugTasksPropParser struct{}

// NewPlugTasksPropParser returns a new PlugTasksPropParser.
func NewPlugTasksPropParser() PlugTasksPropParser {
	return PlugTasksPropParser{}
}

// Trigger implements [parser.InlineParser].
func (PlugTasksPropParser) Trigger() []byte {
	// Only space and punct bytes works here, so we can't just return first byte of emoji.
	// To make it trigger Parse() on emoji we have to trigger on everything possible
	// (to ensure we won't miss anything) and then in Parse() skip up to (but not including)
	// next emoji, to make that emoji a start of next "line".
	var puncts []byte
	for b := byte(0); b < math.MaxUint8; b++ {
		if util.IsPunct(b) {
			puncts = append(puncts, b)
		}
	}
	return append(puncts, '\t', ' ') // ' ' means space or start of line.
}

// Parse implements [parser.InlineParser].
func (PlugTasksPropParser) Parse(parent gast.Node, block text.Reader, _ parser.Context) gast.Node { //nolint:funlen,gocognit // Split will hurt readability.
	// Skip parsing task properties if it's not a task.
	preceding := parent.FirstChild()
	for preceding != nil && preceding.Kind() != ast.KindPlugTasksStatus {
		preceding = preceding.NextSibling()
	}
	if preceding == nil {
		return nil // Not a task: there is no PlugTasksStatus in preceding nodes.
	}

	line, seg := block.PeekLine()

	// Obsidian plugin Tasks parse properties only in first line of a task.
	// See https://publish.obsidian.md/tasks/Getting+Started/Getting+Started#Multi-line+checklist+items.
	if seg.Stop != parent.Lines().At(0).Stop {
		return nil
	}

	// Regexp will match first task property only if rest of line contains only task
	// properties and/or tags and/or block ID.
	// See https://publish.obsidian.md/tasks/Getting+Started/Getting+Started#Order+of+metadata/emojis.
	m := plugTasksPropRegexp.FindSubmatchIndex(line)
	if m == nil {
		// goldmark can call inline parser only at start of "line" or space or punct.
		// So, to make it call parser at emoji we must ensure emoji will be at start
		// of line. For this we need to move prefix of current line (up to but not
		// including next emoji) into own Text node.
		// To ensure this Text node won't overwrite other inline parsers we also avoid
		// including space and punct. Plus run this parser with lowest priority (1000).
		i := int(util.UTF8Len(line[0])) // Advance by at least 1 rune.
		for i < len(line) && !util.IsSpace(line[i]) && !util.IsPunct(line[i]) {
			if plugTasksPropTrigger[line[i]] { // Our emoji may start at this byte.
				block.Advance(i)
				return gast.NewTextSegment(seg.WithStop(seg.Start + i))
			}
			i += int(util.UTF8Len(line[i]))
		}
		return nil
	}

	// Start indexes in m for each group.
	// Only Advance and one of Emoji with it arg (if any) will be set.
	const (
		idxAdvance = iota*2 + 2
		idxEmojiWithoutArgs
		idxEmojiWithID
		idxID
		idxEmojiWithIDs
		idxIDs
		idxEmojiWithDate
		idxDate
		idxEmojiWithRecurring
		idxRecurring
		idxEmojiWithOnCompletionAction
		idxOnCompletionAction
	)
	block.Advance(m[idxAdvance+1])

	var node gast.Node

	switch {
	case m[idxEmojiWithoutArgs] >= 0:
		switch string(line[m[idxEmojiWithoutArgs]:m[idxEmojiWithoutArgs+1]]) {
		case plugTasksPrioLowestEmoji:
			node = ast.NewPlugTasksPrio(ast.PlugTasksPrioLowest)
		case plugTasksPrioLowEmoji:
			node = ast.NewPlugTasksPrio(ast.PlugTasksPrioLow)
		case plugTasksPrioMediumEmoji:
			node = ast.NewPlugTasksPrio(ast.PlugTasksPrioMedium)
		case plugTasksPrioHighEmoji:
			node = ast.NewPlugTasksPrio(ast.PlugTasksPrioHigh)
		case plugTasksPrioHighestEmoji:
			node = ast.NewPlugTasksPrio(ast.PlugTasksPrioHighest)
		default:
			panic("unknown emoji without args")
		}

	case m[idxEmojiWithID] >= 0:
		id := string(line[m[idxID]:m[idxID+1]])

		switch string(line[m[idxEmojiWithID]:m[idxEmojiWithID+1]]) {
		case plugTasksIDEmoji:
			node = ast.NewPlugTasksID(id)
		default:
			panic("unknown emoji with ID")
		}

	case m[idxEmojiWithIDs] >= 0:
		ids := strings.Split(string(line[m[idxIDs]:m[idxIDs+1]]), ",")

		switch string(line[m[idxEmojiWithIDs]:m[idxEmojiWithIDs+1]]) {
		case plugTasksDependsOnEmoji:
			node = ast.NewPlugTasksDependsOn(ids)
		default:
			panic("unknown emoji with IDs")
		}

	case m[idxEmojiWithDate] >= 0:
		date, _ := time.Parse(time.DateOnly, string(line[m[idxDate]:m[idxDate+1]]))

		switch string(line[m[idxEmojiWithDate]:m[idxEmojiWithDate+1]]) {
		case plugTasksDueEmoji:
			node = ast.NewPlugTasksDue(date)
		case plugTasksScheduledEmoji:
			node = ast.NewPlugTasksScheduled(date)
		case plugTasksStartEmoji:
			node = ast.NewPlugTasksStart(date)
		case plugTasksCreatedEmoji:
			node = ast.NewPlugTasksCreated(date)
		case plugTasksDoneEmoji:
			node = ast.NewPlugTasksDone(date)
		case plugTasksCancelledEmoji:
			node = ast.NewPlugTasksCancelled(date)
		default:
			panic("unknown emoji with date")
		}

	case m[idxEmojiWithRecurring] >= 0:
		rule := string(line[m[idxRecurring]:m[idxRecurring+1]])

		switch string(line[m[idxEmojiWithRecurring]:m[idxEmojiWithRecurring+1]]) {
		case plugTasksRecurringEmoji:
			node = ast.NewPlugTasksRecurring(rule)
		default:
			panic("unknown emoji with recurring schedule")
		}

	case m[idxEmojiWithOnCompletionAction] >= 0:
		var action ast.PlugTasksOnCompletionAction
		switch strings.ToLower(string(line[m[idxOnCompletionAction]:m[idxOnCompletionAction+1]])) {
		case "keep":
			action = ast.PlugTasksOnCompletionKeep
		case "delete":
			action = ast.PlugTasksOnCompletionDelete
		}

		switch string(line[m[idxEmojiWithOnCompletionAction]:m[idxEmojiWithOnCompletionAction+1]]) {
		case plugTasksOnCompletionEmoji:
			node = ast.NewPlugTasksOnCompletion(action)
		default:
			panic("unknown emoji with on completion action")
		}

	default:
		panic("unknown emoji")
	}

	// Only first property of each kind is used (even if it has invalid arg).
	for preceding != nil && preceding.Kind() != node.Kind() {
		preceding = preceding.NextSibling()
	}
	if preceding != nil { // Current property was already defined, so ignore following ones.
		return gast.NewText() // Must return non-nil to apply block.Advance() call.
	}

	return node
}

// PlugTasksHTMLRenderer is an HTML renderer for Obsidian plugin Tasks's tasks.
//
// It renders task status checkbox and task properties in format similar to produced by saving
// Obsidian plugin Tasks's modal window with task details (add missing spaces, drop invalid
// properties, etc.).
type PlugTasksHTMLRenderer struct {
	html.Config // Embed to implement renderer.SetOptioner and apply global html options.
}

// NewPlugTasksHTMLRenderer returns a new PlugTasksHTMLRenderer.
func NewPlugTasksHTMLRenderer() *PlugTasksHTMLRenderer {
	r := &PlugTasksHTMLRenderer{
		Config: html.NewConfig(),
	}
	return r
}

// RegisterFuncs implements [renderer.NodeRenderer].
func (r *PlugTasksHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindPlugTasksStatus, entering(r.renderStatus))
	reg.Register(ast.KindPlugTasksPrio, entering(r.renderPrio))
	reg.Register(ast.KindPlugTasksID, entering(r.renderID))
	reg.Register(ast.KindPlugTasksDependsOn, entering(r.renderDependsOn))
	reg.Register(ast.KindPlugTasksDue, entering(r.renderDate))
	reg.Register(ast.KindPlugTasksScheduled, entering(r.renderDate))
	reg.Register(ast.KindPlugTasksStart, entering(r.renderDate))
	reg.Register(ast.KindPlugTasksCreated, entering(r.renderDate))
	reg.Register(ast.KindPlugTasksDone, entering(r.renderDate))
	reg.Register(ast.KindPlugTasksCancelled, entering(r.renderDate))
	reg.Register(ast.KindPlugTasksRecurring, entering(r.renderRecurring))
	reg.Register(ast.KindPlugTasksOnCompletion, entering(r.renderOnCompletion))
}

// PlugTasksCheckboxAttributeFilter defines attribute names which <input type=checkbox>
// elements can have.
var PlugTasksCheckboxAttributeFilter = html.GlobalAttributeFilter //nolint:gochecknoglobals // By design.

func (r *PlugTasksHTMLRenderer) renderStatus(w util.BufWriter, _ []byte, node gast.Node) (gast.WalkStatus, error) {
	n := node.(*ast.PlugTasksStatus)

	if n.IsChecked() {
		_, _ = w.WriteString(`<input checked="" disabled="" type="checkbox"`)
	} else {
		_, _ = w.WriteString(`<input disabled="" type="checkbox"`)
	}

	if n.Attributes() != nil {
		html.RenderAttributes(w, n, PlugTasksCheckboxAttributeFilter)
	}

	if r.XHTML {
		_, _ = w.WriteString(` /> `)
	} else {
		_, _ = w.WriteString(`> `)
	}
	return gast.WalkContinue, nil
}

func (*PlugTasksHTMLRenderer) renderPrio(w util.BufWriter, _ []byte, node gast.Node) (gast.WalkStatus, error) {
	n := node.(*ast.PlugTasksPrio)

	if n.Prio == ast.PlugTasksPrioDefault {
		return gast.WalkContinue, nil
	}

	_ = w.WriteByte(' ')
	switch n.Prio {
	case ast.PlugTasksPrioLowest:
		_, _ = w.WriteString(plugTasksPrioLowestEmoji)
	case ast.PlugTasksPrioLow:
		_, _ = w.WriteString(plugTasksPrioLowEmoji)
	case ast.PlugTasksPrioDefault: // Make linter exhaustive happy without disabling it.
		panic("never here")
	case ast.PlugTasksPrioMedium:
		_, _ = w.WriteString(plugTasksPrioMediumEmoji)
	case ast.PlugTasksPrioHigh:
		_, _ = w.WriteString(plugTasksPrioHighEmoji)
	case ast.PlugTasksPrioHighest:
		_, _ = w.WriteString(plugTasksPrioHighestEmoji)
	default:
		panic(fmt.Sprintf("unknown node priority: %v", n.Prio))
	}

	return gast.WalkContinue, nil
}

func (*PlugTasksHTMLRenderer) renderID(w util.BufWriter, _ []byte, node gast.Node) (gast.WalkStatus, error) {
	n := node.(*ast.PlugTasksID)

	_ = w.WriteByte(' ')
	_, _ = w.WriteString(plugTasksIDEmoji)
	_ = w.WriteByte(' ')
	_, _ = w.WriteString(n.ID)

	return gast.WalkContinue, nil
}

func (*PlugTasksHTMLRenderer) renderDependsOn(w util.BufWriter, _ []byte, node gast.Node) (gast.WalkStatus, error) {
	n := node.(*ast.PlugTasksDependsOn)

	_ = w.WriteByte(' ')
	_, _ = w.WriteString(plugTasksDependsOnEmoji)
	_ = w.WriteByte(' ')
	_, _ = w.WriteString(strings.Join(n.IDs, ","))

	return gast.WalkContinue, nil
}

func (*PlugTasksHTMLRenderer) renderDate(w util.BufWriter, _ []byte, node gast.Node) (gast.WalkStatus, error) {
	// Do not output properties with invalid date.
	// This isn't actually match behaviour of Obsidian plugin Tasks's modal window
	// because it's impossible to save with invalid dates at all.
	// But this behaviour is consistent with handling all other errors.
	if !node.(interface{ IsValid() bool }).IsValid() {
		return gast.WalkContinue, nil
	}

	var date time.Time
	_ = w.WriteByte(' ')
	switch n := node.(type) {
	case *ast.PlugTasksDue:
		_, _ = w.WriteString(plugTasksDueEmoji)
		date = n.Date
	case *ast.PlugTasksScheduled:
		_, _ = w.WriteString(plugTasksScheduledEmoji)
		date = n.Date
	case *ast.PlugTasksStart:
		_, _ = w.WriteString(plugTasksStartEmoji)
		date = n.Date
	case *ast.PlugTasksCreated:
		_, _ = w.WriteString(plugTasksCreatedEmoji)
		date = n.Date
	case *ast.PlugTasksDone:
		_, _ = w.WriteString(plugTasksDoneEmoji)
		date = n.Date
	case *ast.PlugTasksCancelled:
		_, _ = w.WriteString(plugTasksCancelledEmoji)
		date = n.Date
	default:
		panic(fmt.Sprintf("unknown node with date: %T", node))
	}
	_ = w.WriteByte(' ')
	_, _ = w.WriteString(date.Format(time.DateOnly))

	return gast.WalkContinue, nil
}

func (*PlugTasksHTMLRenderer) renderRecurring(w util.BufWriter, _ []byte, node gast.Node) (gast.WalkStatus, error) {
	n := node.(*ast.PlugTasksRecurring)

	_ = w.WriteByte(' ')
	_, _ = w.WriteString(plugTasksRecurringEmoji)
	_ = w.WriteByte(' ')
	_, _ = w.WriteString(n.Rule)

	return gast.WalkContinue, nil
}

func (*PlugTasksHTMLRenderer) renderOnCompletion(w util.BufWriter, _ []byte, node gast.Node) (gast.WalkStatus, error) {
	n := node.(*ast.PlugTasksOnCompletion)

	if !n.IsValid() {
		return gast.WalkContinue, nil
	}

	_ = w.WriteByte(' ')
	_, _ = w.WriteString(plugTasksOnCompletionEmoji)
	_ = w.WriteByte(' ')
	_, _ = w.WriteString(n.Action.String())

	return gast.WalkContinue, nil
}

// PlugTasks is an extension that helps to setup Obsidian plugin [Tasks] parser and HTML renderer.
//
// [Tasks]: https://publish.obsidian.md/tasks/Introduction
type PlugTasks struct {
	opts []PlugTasksOption
}

// NewPlugTasks returns a new PlugTasks extension.
func NewPlugTasks(opts ...PlugTasksOption) PlugTasks {
	return PlugTasks{opts: opts}
}

// Extend implements [goldmark.Extender].
func (e PlugTasks) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		// Not sure why it's important to use so high priority, but
		// official extension.TaskList uses this priority.
		util.Prioritized(NewPlugTasksParser(e.opts...), prioHighest),
		// This parser creates Text nodes for everything up to next emoji,
		// so run it with lowest priority to make sure it won't turn into Text node
		// something what should be handled by some other parser.
		util.Prioritized(NewPlugTasksPropParser(), prioInlineParserLowest),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewPlugTasksHTMLRenderer(), prioHTMLRenderer),
	))
}
