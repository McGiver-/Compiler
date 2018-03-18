package Syn

import (
	"testing"

	"github.com/McGiver-/Compiler/Lex"
	"github.com/davecgh/go-spew/spew"
)

// func TestMakeNode(t *testing.T) {
// 	tokenX := &Lex.Token{
// 		"id",
// 		"name",
// 		"10:10",
// 	}
// 	nodeX := makeNode(tokenX)
// 	t.Log(spew.Sdump(nodeX))
// }

// func TestMakeSiblings(t *testing.T) {
// 	tokenX := &Lex.Token{
// 		"id",
// 		"name",
// 		"10:10",
// 	}

// 	tokenY := &Lex.Token{
// 		"intNum",
// 		"43",
// 		"10:10",
// 	}

// 	tokenY2 := &Lex.Token{
// 		"floatNum",
// 		"44.44",
// 		"10:10",
// 	}

// 	nodeX := makeNode(tokenX)
// 	nodeY := makeNode(tokenY)
// 	nodeY2 := makeNode(tokenY2)

// 	nodeX.makeSiblings(nodeY, nodeY2)

// 	t.Log(spew.Sdump(nodeX))

// }

// func TestAdoptChildren(t *testing.T) {
// 	tParent := &Lex.Token{
// 		"id",
// 		"Parent",
// 		"10:10",
// 	}
// 	tokenX := &Lex.Token{
// 		"id",
// 		"name",
// 		"10:10",
// 	}

// 	tokenY := &Lex.Token{
// 		"intNum",
// 		"43",
// 		"10:10",
// 	}

// 	tokenY2 := &Lex.Token{
// 		"floatNum",
// 		"44.44",
// 		"10:10",
// 	}

// 	nodeParent := makeNode(tParent)
// 	nodeX := makeNode(tokenX)
// 	nodeY := makeNode(tokenY)
// 	nodeY2 := makeNode(tokenY2)

// 	nodeX.makeSiblings(nodeY, nodeY2)
// 	nodeParent.adoptChildren(nodeX)
// 	t.Log(spew.Sdump(nodeParent))

// }

func TestMakeFamily(t *testing.T) {
	tParent := &Lex.Token{
		"id",
		"Parent",
		"10:10",
	}
	tokenX := &Lex.Token{
		"id",
		"name",
		"10:10",
	}

	// tokenY := &Lex.Token{
	// 	"intNum",
	// 	"43",
	// 	"10:10",
	// }

	// tokenY2 := &Lex.Token{
	// 	"floatNum",
	// 	"44.44",
	// 	"10:10",
	// }

	nodeX := makeNode("child", "child", tokenX)
	// nodeY := makeNode(tokenY)
	// nodeY2 := makeNode(tokenY2)

	t.Log(spew.Sdump(makeFamily("parent", "parent", tParent, nodeX)))

}
