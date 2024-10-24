package ast

import (
	"fmt"

	gast "github.com/yuin/goldmark/ast"
)

// A BlockID struct represents an Obsidian block id.
// https://help.obsidian.md/Linking+notes+and+files/Internal+links#Link+to+a+block+in+a+note
type BlockID struct {
	gast.BaseInline
	ID []byte
	// Invalid means this block id can't be used as a reference target.
	// This is how Obsidian works, so we are just following it behaviour.
	Invalid bool
}

// Dump implements Node.Dump.
func (n *BlockID) Dump(source []byte, level int) {
	m := map[string]string{
		"ID":      string(n.ID),
		"Invalid": fmt.Sprintf("%v", n.Invalid),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindBlockID is a NodeKind of the BlockID node.
var KindBlockID = gast.NewNodeKind("BlockID") // Const.

// Kind implements Node.Kind.
func (*BlockID) Kind() gast.NodeKind {
	return KindBlockID
}

// NewBlockID returns a new valid BlockID node.
func NewBlockID(id []byte) *BlockID {
	return &BlockID{ID: id, Invalid: false}
}

// NewInvalidBlockID returns a new invalid BlockID node.
func NewInvalidBlockID(id []byte) *BlockID {
	return &BlockID{ID: id, Invalid: true}
}
