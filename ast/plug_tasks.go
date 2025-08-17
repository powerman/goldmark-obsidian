package ast

import (
	"fmt"
	"strings"
	"time"

	gast "github.com/yuin/goldmark/ast"
)

// PlugTasksStatusType represents Obsidian plugin Tasks's [status type].
//
// [status type]: https://publish.obsidian.md/tasks/Getting+Started/Statuses/Status+Types
type PlugTasksStatusType int

// Obsidian plugin Tasks's status types.
const (
	PlugTasksStatusTypeTODO PlugTasksStatusType = iota + 1
	PlugTasksStatusTypeInProgress
	PlugTasksStatusTypeDone
	PlugTasksStatusTypeCancelled
	PlugTasksStatusTypeNonTask
)

// String implements [fmt.Stringer].
func (t PlugTasksStatusType) String() string {
	switch t {
	case PlugTasksStatusTypeTODO:
		return "TODO"
	case PlugTasksStatusTypeInProgress:
		return "InProgress"
	case PlugTasksStatusTypeDone:
		return "Done"
	case PlugTasksStatusTypeCancelled:
		return "Cancelled"
	case PlugTasksStatusTypeNonTask:
		return "NonTask"
	default:
		return "Unknown"
	}
}

// IsDone reports whether the task is considered "done" (has status Done, Cancelled or NonTask).
func (t PlugTasksStatusType) IsDone() bool {
	switch t {
	case PlugTasksStatusTypeTODO, PlugTasksStatusTypeInProgress:
		return false
	case PlugTasksStatusTypeDone, PlugTasksStatusTypeCancelled, PlugTasksStatusTypeNonTask:
		return true
	}
	panic(fmt.Sprintf("invalid PlugTasksStatusType: %v", t))
}

// A PlugTasksStatus represents a checkbox of an Obsidian plugin Tasks's task.
type PlugTasksStatus struct {
	gast.BaseInline

	Symbol     rune
	StatusType PlugTasksStatusType
}

// NewPlugTasksStatus returns a new PlugTasksStatus node.
func NewPlugTasksStatus(symbol rune, statusType PlugTasksStatusType) *PlugTasksStatus {
	return &PlugTasksStatus{
		Symbol:     symbol,
		StatusType: statusType,
	}
}

// IsChecked reports whether Symbol != ' '.
func (n *PlugTasksStatus) IsChecked() bool {
	return n.Symbol != ' '
}

// Dump implements [ast.Node].
func (n *PlugTasksStatus) Dump(source []byte, level int) {
	m := map[string]string{
		"Symbol":     fmt.Sprintf("%q", n.Symbol),
		"StatusType": n.StatusType.String(),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindPlugTasksStatus is a NodeKind of the PlugTasksStatus node.
var KindPlugTasksStatus = gast.NewNodeKind("PlugTasksStatus") // Const.

// Kind implements [ast.Node].
func (*PlugTasksStatus) Kind() gast.NodeKind {
	return KindPlugTasksStatus
}

// PlugTasksPriority is an Obsidian plugin Tasks's [priority].
//
// [priority]: https://publish.obsidian.md/tasks/Getting+Started/Priority
type PlugTasksPriority int

// Obsidian plugin Tasks's priorities.
const (
	PlugTasksPrioLowest PlugTasksPriority = iota - 2
	PlugTasksPrioLow
	PlugTasksPrioDefault
	PlugTasksPrioMedium
	PlugTasksPrioHigh
	PlugTasksPrioHighest
)

// String implements [fmt.Stringer].
func (t PlugTasksPriority) String() string {
	switch t {
	case PlugTasksPrioLowest:
		return "Lowest"
	case PlugTasksPrioLow:
		return "Low"
	case PlugTasksPrioDefault:
		return ""
	case PlugTasksPrioMedium:
		return "Medium"
	case PlugTasksPrioHigh:
		return "High"
	case PlugTasksPrioHighest:
		return "Highest"
	default:
		return "Unknown"
	}
}

// A PlugTasksPrio represents an Obsidian plugin Tasks's task priority.
type PlugTasksPrio struct {
	gast.BaseInline

	Prio PlugTasksPriority
}

// NewPlugTasksPrio returns a new PlugTasksPrio node.
func NewPlugTasksPrio(prio PlugTasksPriority) *PlugTasksPrio {
	return &PlugTasksPrio{
		Prio: prio,
	}
}

// Dump implements [ast.Node].
func (n *PlugTasksPrio) Dump(source []byte, level int) {
	m := map[string]string{
		"Prio": n.Prio.String(),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindPlugTasksPrio is a NodeKind of the PlugTasksPrio node.
var KindPlugTasksPrio = gast.NewNodeKind("PlugTasksPrio") // Const.

// Kind implements [ast.Node].
func (*PlugTasksPrio) Kind() gast.NodeKind {
	return KindPlugTasksPrio
}

// A PlugTasksID represents an Obsidian plugin Tasks's task ID.
type PlugTasksID struct {
	gast.BaseInline

	ID string
}

// NewPlugTasksID returns a new PlugTasksID node.
func NewPlugTasksID(id string) *PlugTasksID {
	return &PlugTasksID{
		ID: id,
	}
}

// Dump implements [ast.Node].
func (n *PlugTasksID) Dump(source []byte, level int) {
	m := map[string]string{
		"ID": n.ID,
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindPlugTasksID is a NodeKind of the PlugTasksID node.
var KindPlugTasksID = gast.NewNodeKind("PlugTasksID") // Const.

// Kind implements [ast.Node].
func (*PlugTasksID) Kind() gast.NodeKind {
	return KindPlugTasksID
}

// A PlugTasksDependsOn represents an Obsidian plugin Tasks's task dependencies.
type PlugTasksDependsOn struct {
	gast.BaseInline

	IDs []string
}

// NewPlugTasksDependsOn returns a new PlugTasksDependsOn node.
func NewPlugTasksDependsOn(ids []string) *PlugTasksDependsOn {
	return &PlugTasksDependsOn{
		IDs: ids,
	}
}

// Dump implements [ast.Node].
func (n *PlugTasksDependsOn) Dump(source []byte, level int) {
	m := map[string]string{
		"IDs": strings.Join(n.IDs, ","),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindPlugTasksDependsOn is a NodeKind of the PlugTasksDependsOn node.
var KindPlugTasksDependsOn = gast.NewNodeKind("PlugTasksDependsOn") // Const.

// Kind implements [ast.Node].
func (*PlugTasksDependsOn) Kind() gast.NodeKind {
	return KindPlugTasksDependsOn
}

// A PlugTasksDue represents an Obsidian plugin Tasks's task due date.
type PlugTasksDue struct {
	gast.BaseInline

	Date time.Time
}

// NewPlugTasksDue returns a new PlugTasksDue node.
func NewPlugTasksDue(date time.Time) *PlugTasksDue {
	return &PlugTasksDue{
		Date: date,
	}
}

// IsValid reports whether date is valid.
func (n *PlugTasksDue) IsValid() bool {
	return !n.Date.IsZero()
}

// Dump implements [ast.Node].
func (n *PlugTasksDue) Dump(source []byte, level int) {
	m := map[string]string{
		"Date": n.Date.Format(time.DateOnly),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindPlugTasksDue is a NodeKind of the PlugTasksDue node.
var KindPlugTasksDue = gast.NewNodeKind("PlugTasksDue") // Const.

// Kind implements [ast.Node].
func (*PlugTasksDue) Kind() gast.NodeKind {
	return KindPlugTasksDue
}

// A PlugTasksScheduled represents an Obsidian plugin Tasks's task scheduled date.
type PlugTasksScheduled struct {
	gast.BaseInline

	Date time.Time
}

// NewPlugTasksScheduled returns a new PlugTasksScheduled node.
func NewPlugTasksScheduled(date time.Time) *PlugTasksScheduled {
	return &PlugTasksScheduled{
		Date: date,
	}
}

// IsValid reports whether date is valid.
func (n *PlugTasksScheduled) IsValid() bool {
	return !n.Date.IsZero()
}

// Dump implements [ast.Node].
func (n *PlugTasksScheduled) Dump(source []byte, level int) {
	m := map[string]string{
		"Date": n.Date.Format(time.DateOnly),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindPlugTasksScheduled is a NodeKind of the PlugTasksScheduled node.
var KindPlugTasksScheduled = gast.NewNodeKind("PlugTasksScheduled") // Const.

// Kind implements [ast.Node].
func (*PlugTasksScheduled) Kind() gast.NodeKind {
	return KindPlugTasksScheduled
}

// A PlugTasksStart represents an Obsidian plugin Tasks's task start date.
type PlugTasksStart struct {
	gast.BaseInline

	Date time.Time
}

// NewPlugTasksStart returns a new PlugTasksStart node.
func NewPlugTasksStart(date time.Time) *PlugTasksStart {
	return &PlugTasksStart{
		Date: date,
	}
}

// IsValid reports whether date is valid.
func (n *PlugTasksStart) IsValid() bool {
	return !n.Date.IsZero()
}

// Dump implements [ast.Node].
func (n *PlugTasksStart) Dump(source []byte, level int) {
	m := map[string]string{
		"Date": n.Date.Format(time.DateOnly),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindPlugTasksStart is a NodeKind of the PlugTasksStart node.
var KindPlugTasksStart = gast.NewNodeKind("PlugTasksStart") // Const.

// Kind implements [ast.Node].
func (*PlugTasksStart) Kind() gast.NodeKind {
	return KindPlugTasksStart
}

// A PlugTasksCreated represents an Obsidian plugin Tasks's task created date.
type PlugTasksCreated struct {
	gast.BaseInline

	Date time.Time
}

// NewPlugTasksCreated returns a new PlugTasksCreated node.
func NewPlugTasksCreated(date time.Time) *PlugTasksCreated {
	return &PlugTasksCreated{
		Date: date,
	}
}

// IsValid reports whether date is valid.
func (n *PlugTasksCreated) IsValid() bool {
	return !n.Date.IsZero()
}

// Dump implements [ast.Node].
func (n *PlugTasksCreated) Dump(source []byte, level int) {
	m := map[string]string{
		"Date": n.Date.Format(time.DateOnly),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindPlugTasksCreated is a NodeKind of the PlugTasksCreated node.
var KindPlugTasksCreated = gast.NewNodeKind("PlugTasksCreated") // Const.

// Kind implements [ast.Node].
func (*PlugTasksCreated) Kind() gast.NodeKind {
	return KindPlugTasksCreated
}

// A PlugTasksDone represents an Obsidian plugin Tasks's task done date.
type PlugTasksDone struct {
	gast.BaseInline

	Date time.Time
}

// NewPlugTasksDone returns a new PlugTasksDone node.
func NewPlugTasksDone(date time.Time) *PlugTasksDone {
	return &PlugTasksDone{
		Date: date,
	}
}

// IsValid reports whether date is valid.
func (n *PlugTasksDone) IsValid() bool {
	return !n.Date.IsZero()
}

// Dump implements [ast.Node].
func (n *PlugTasksDone) Dump(source []byte, level int) {
	m := map[string]string{
		"Date": n.Date.Format(time.DateOnly),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindPlugTasksDone is a NodeKind of the PlugTasksDone node.
var KindPlugTasksDone = gast.NewNodeKind("PlugTasksDone") // Const.

// Kind implements [ast.Node].
func (*PlugTasksDone) Kind() gast.NodeKind {
	return KindPlugTasksDone
}

// A PlugTasksCancelled represents an Obsidian plugin Tasks's task cancelled date.
type PlugTasksCancelled struct {
	gast.BaseInline

	Date time.Time
}

// NewPlugTasksCancelled returns a new PlugTasksCancelled node.
func NewPlugTasksCancelled(date time.Time) *PlugTasksCancelled {
	return &PlugTasksCancelled{
		Date: date,
	}
}

// IsValid reports whether date is valid.
func (n *PlugTasksCancelled) IsValid() bool {
	return !n.Date.IsZero()
}

// Dump implements [ast.Node].
func (n *PlugTasksCancelled) Dump(source []byte, level int) {
	m := map[string]string{
		"Date": n.Date.Format(time.DateOnly),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindPlugTasksCancelled is a NodeKind of the PlugTasksCancelled node.
var KindPlugTasksCancelled = gast.NewNodeKind("PlugTasksCancelled") // Const.

// Kind implements [ast.Node].
func (*PlugTasksCancelled) Kind() gast.NodeKind {
	return KindPlugTasksCancelled
}

// A PlugTasksRecurring represents an Obsidian plugin Tasks's recurring task.
//
// Rule is not validated and may be invalid.
type PlugTasksRecurring struct {
	gast.BaseInline

	Rule string
}

// NewPlugTasksRecurring returns a new PlugTasksRecurring node.
func NewPlugTasksRecurring(rule string) *PlugTasksRecurring {
	return &PlugTasksRecurring{
		Rule: rule,
	}
}

// Dump implements [ast.Node].
func (n *PlugTasksRecurring) Dump(source []byte, level int) {
	m := map[string]string{
		"Rule": fmt.Sprintf("%q", n.Rule),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindPlugTasksRecurring is a NodeKind of the PlugTasksRecurring node.
var KindPlugTasksRecurring = gast.NewNodeKind("PlugTasksRecurring") // Const.

// Kind implements [ast.Node].
func (*PlugTasksRecurring) Kind() gast.NodeKind {
	return KindPlugTasksRecurring
}

// PlugTasksOnCompletionAction is an Obsidian plugin Tasks's action on completion.
// See https://publish.obsidian.md/tasks/Getting+Started/On+Completion.
type PlugTasksOnCompletionAction int

// Obsidian plugin Tasks's action on completion.
const (
	PlugTasksOnCompletionKeep PlugTasksOnCompletionAction = iota + 1
	PlugTasksOnCompletionDelete
)

// String implements [fmt.Stringer].
func (t PlugTasksOnCompletionAction) String() string {
	switch t {
	case PlugTasksOnCompletionKeep:
		return "keep"
	case PlugTasksOnCompletionDelete:
		return "delete"
	default:
		return "unknown"
	}
}

// A PlugTasksOnCompletion represents an Obsidian plugin Tasks's task action on completion.
type PlugTasksOnCompletion struct {
	gast.BaseInline

	Action PlugTasksOnCompletionAction
}

// NewPlugTasksOnCompletion returns a new PlugTasksOnCompletion node.
func NewPlugTasksOnCompletion(action PlugTasksOnCompletionAction) *PlugTasksOnCompletion {
	return &PlugTasksOnCompletion{
		Action: action,
	}
}

// IsValid reports whether action is valid.
func (n *PlugTasksOnCompletion) IsValid() bool {
	return n.Action == PlugTasksOnCompletionKeep || n.Action == PlugTasksOnCompletionDelete
}

// Dump implements [ast.Node].
func (n *PlugTasksOnCompletion) Dump(source []byte, level int) {
	m := map[string]string{
		"Action": n.Action.String(),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindPlugTasksOnCompletion is a NodeKind of the PlugTasksOnCompletion node.
var KindPlugTasksOnCompletion = gast.NewNodeKind("PlugTasksOnCompletion") // Const.

// Kind implements [ast.Node].
func (*PlugTasksOnCompletion) Kind() gast.NodeKind {
	return KindPlugTasksOnCompletion
}
