package obsidian_test

import (
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/testutil"

	obsidian "github.com/powerman/goldmark-obsidian"
)

func TestBlockID(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(obsidian.NewBlockID()),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	testutil.DoTestCaseFile(markdown, "testdata/block_id.txt", t, testutil.ParseCliCaseArg()...)
}

func TestBlockIDWithAutoHeadingID(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(obsidian.NewBlockID()),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)
	testutil.DoTestCaseFile(markdown, "testdata/block_id_auto.txt", t, testutil.ParseCliCaseArg()...)
}
