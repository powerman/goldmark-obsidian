package obsidian

import (
	"bytes"
	"slices"

	gast "github.com/yuin/goldmark/ast"
)

var (
	attrNameID    = []byte("id")    // Const.
	attrNameClass = []byte("class") // Const.
)

func appendClass(n gast.Node, class []byte) {
	if len(class) == 0 { // Make it easier to use configurable empty class names.
		return
	}

	val, found := n.Attribute(attrNameClass)
	if !found {
		n.SetAttribute(attrNameClass, bytes.Clone(class))
	} else {
		value := val.([]byte)
		for _, existing := range bytes.Fields(value) {
			if bytes.Equal(class, existing) {
				return
			}
		}
		value = slices.Grow(value, 1+len(class))
		value = append(append(value, ' '), class...)
		n.SetAttribute(attrNameClass, value)
	}
}
