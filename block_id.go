package obsidian

import (
	"regexp"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/powerman/goldmark-obsidian/ast"
)

var blockIDRegexp = regexp.MustCompile(`^\^[A-Za-z0-9-]+$`)

// BlockIDParser is an Obsidian block id parser.
type BlockIDParser struct{}

// NewBlockIDParser returns a new BlockIDParser.
func NewBlockIDParser() BlockIDParser {
	return BlockIDParser{}
}

// Trigger implements [parser.InlineParser].
func (BlockIDParser) Trigger() []byte {
	return []byte{'^'}
}

// Parse implements [parser.InlineParser].
func (BlockIDParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	// Block ID is always at the end of line, so either whole line is an ID or it's not ID.
	id, _ := block.PeekLine()
	if !blockIDRegexp.Match(id) {
		return nil
	}

	// Block ID must be space-delimited from preceding text, if any.
	// Ensure previous child is either not a Text or that Text ends with a line break or
	// that Text ends with a space (and trim right spaces in the last case).
	switch prev, ok := parent.LastChild().(*gast.Text); {
	case !ok:
		// Preceded by not a text: OK.
	case prev.SoftLineBreak() || prev.HardLineBreak():
		// Preceded by line break: OK.
	case segmentEndsWithSpace(prev.Segment, block.Source()):
		// Preceded by space(s): trim them.
		prev.Segment = prev.Segment.TrimRightSpace(block.Source())
		if prev.Segment.IsEmpty() {
			prev.Parent().RemoveChild(prev.Parent(), prev)
		}
	default:
		return nil
	}

	block.Advance(len(id))

	// ID attribute should be set for some parent node.
	target := parent
	switch root := target.Parent(); {
	// Set list item ID from it TextBlock child.
	case target.Kind() == gast.KindTextBlock && root.Kind() == gast.KindListItem:
		target = root
	// Set list item/blockquote ID from a last paragraph of a list item/blockquote.
	case target.Kind() == gast.KindParagraph &&
		(root.Kind() == gast.KindListItem || root.Kind() == gast.KindBlockquote):
		if root.LastChild() != target {
			return ast.NewInvalidBlockID(id)
		}
		target = root
	}

	// Obsidian do not set ID for child elements inside a blockquote.
	for p := target.Parent(); p != nil; p = p.Parent() {
		if p.Kind() == gast.KindBlockquote {
			return ast.NewInvalidBlockID(id)
		}
	}
	// Obsidian do not set ID for empty paragraph containing only ^block-id.
	if target.Kind() == gast.KindParagraph && len(target.Lines().Value(block.Source())) == len(id) {
		return ast.NewInvalidBlockID(id)
	}

	target.SetAttribute(attrNameID, id)
	pc.IDs().Put(id)

	return ast.NewBlockID(id)
}

// BlockIDHTMLRenderer is an HTML renderer for Obsidian block id.
//
// Current implementation does not render block id at all (like Obsidian's "Reading" mode).
type BlockIDHTMLRenderer struct{}

// NewBlockIDHTMLRenderer returns a new BlockIDHTMLRenderer.
func NewBlockIDHTMLRenderer() BlockIDHTMLRenderer {
	return BlockIDHTMLRenderer{}
}

// RegisterFuncs implements [renderer.NodeRenderer].
func (BlockIDHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindBlockID, nil)
}

// BlockID is an extension that helps to setup Obsidian [block id] parser and HTML renderer.
//
// [block id]: https://help.obsidian.md/Linking+notes+and+files/Internal+links#Link+to+a+block+in+a+note
type BlockID struct{}

// NewBlockID returns a new BlockID extension.
func NewBlockID() BlockID {
	return BlockID{}
}

// Extend implements [goldmark.Extender].
func (BlockID) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewBlockIDParser(), prioHighest),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewBlockIDHTMLRenderer(), prioHTMLRenderer),
	))
}
