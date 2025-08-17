// Package obsidian provides [github.com/yuin/goldmark] markdown parser extensions for
// [Obsidian] Flavored Markdown and some 3rd-party Obsidian plugins.
//
// [Obsidian]: https://obsidian.md/
package obsidian

import (
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"go.abhg.dev/goldmark/hashtag"
	"go.abhg.dev/goldmark/mermaid"
	"go.abhg.dev/goldmark/wikilink"
)

// Obsidian is an extension that helps to setup parsers and HTML renderers for
// [Obsidian Flavored Markdown] and all other Obsidian features like [Footnotes], [Diagrams],
// [LaTeX], [Tags] and [Properties].
//
// It's recommended to add other extensions (e.g. Obsidian plugins) before this one to
// ensure correct parsing priorities. E.g.:
//
//	goldmark.WithExtensions(
//	    obsidian.NewPlugTasks(),
//	    obsidian.NewObsidian(),
//	)
//
// [Obsidian Flavored Markdown]: https://help.obsidian.md/Editing+and+formatting/Obsidian+Flavored+Markdown
// [Footnotes]: https://help.obsidian.md/Editing+and+formatting/Basic+formatting+syntax#Footnotes
// [Diagrams]: https://help.obsidian.md/Editing+and+formatting/Advanced+formatting+syntax#Diagram
// [LaTeX]: https://help.obsidian.md/Editing+and+formatting/Advanced+formatting+syntax#Math
// [Tags]: https://help.obsidian.md/Editing+and+formatting/Tags
// [Properties]: https://help.obsidian.md/Editing+and+formatting/Properties
type Obsidian struct {
	linkify  goldmark.Extender
	table    goldmark.Extender
	footnote goldmark.Extender
	meta     goldmark.Extender
	hashtag  hashtag.Extender
	wikilink wikilink.Extender
	mermaid  mermaid.Extender
	mathjax  goldmark.Extender
}

// NewObsidian returns a new Obsidian extension.
func NewObsidian() Obsidian {
	return Obsidian{
		linkify:  extension.Linkify,
		table:    extension.Table,
		footnote: extension.Footnote,
		meta:     meta.Meta,
		hashtag: hashtag.Extender{
			Variant: hashtag.ObsidianVariant,
		},
		wikilink: wikilink.Extender{},
		mermaid:  mermaid.Extender{},
		mathjax:  mathjax.MathJax,
	}
}

// WithLinkifyOptions returns a new extension modified to use given opts.
func (e Obsidian) WithLinkifyOptions(opts ...extension.LinkifyOption) Obsidian {
	e.linkify = extension.NewLinkify(opts...)
	return e
}

// WithTableOptions returns a new extension modified to use given opts.
func (e Obsidian) WithTableOptions(opts ...extension.TableOption) Obsidian {
	e.table = extension.NewTable(opts...)
	return e
}

// WithFootnoteOptions returns a new extension modified to use given opts.
func (e Obsidian) WithFootnoteOptions(opts ...extension.FootnoteOption) Obsidian {
	e.footnote = extension.NewFootnote(opts...)
	return e
}

// WithMetaOptions returns a new extension modified to use given opts.
func (e Obsidian) WithMetaOptions(opts ...meta.Option) Obsidian {
	e.meta = meta.New(opts...)
	return e
}

// WithHashtagResolver returns a new extension modified to use given resolver.
func (e Obsidian) WithHashtagResolver(resolver hashtag.Resolver) Obsidian {
	e.hashtag.Resolver = resolver
	return e
}

// WithWikilinkResolver returns a new extension modified to use given resolver.
func (e Obsidian) WithWikilinkResolver(resolver wikilink.Resolver) Obsidian {
	e.wikilink.Resolver = resolver
	return e
}

// WithMermaid returns a new extension modified to use given m.
func (e Obsidian) WithMermaid(m mermaid.Extender) Obsidian {
	e.mermaid = m
	return e
}

// WithMathJaxOptions returns a new extension modified to use given opts.
func (e Obsidian) WithMathJaxOptions(opts ...mathjax.Option) Obsidian {
	e.mathjax = mathjax.NewMathJax(opts...)
	return e
}

// Extend implements [goldmark.Extender].
func (e Obsidian) Extend(m goldmark.Markdown) {
	e.meta.Extend(m)
	NewBlockID().Extend(m)
	e.hashtag.Extend(m)
	e.wikilink.Extend(m)
	e.mermaid.Extend(m)
	e.mathjax.Extend(m)
	e.linkify.Extend(m)
	e.table.Extend(m)
	extension.Strikethrough.Extend(m)
	extension.TaskList.Extend(m)
	e.footnote.Extend(m)
}
