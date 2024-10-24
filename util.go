package obsidian

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type enterNodeRendererFunc func(writer util.BufWriter, source []byte, n ast.Node) (ast.WalkStatus, error)

func entering(f enterNodeRendererFunc) renderer.NodeRendererFunc {
	return func(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		return f(w, source, node)
	}
}

// SegmentsEndsWithSpace checks actual segment value (ignoring padding).
func segmentEndsWithSpace(t text.Segment, source []byte) bool {
	v := source[t.Start:t.Stop]
	return len(v) > 0 && util.IsSpace(v[len(v)-1])
}
