package Syn

import "github.com/McGiver-/Compiler/Lex"

type Node struct {
	Type            string
	Value           string
	Token           *Lex.Token
	Parent          *Node
	LeftMostChild   *Node
	LeftMostSibling *Node
	RightSibling    *Node
}

type noder interface {
	makeSiblings(noder) noder
	adoptChildren(noder) noder
	makeFamily(string, ...noder) noder
	get()
	set()
}

func makeSiblings(x, y *Node) *Node {
	xRight := x
	yRight := y

	for xRight.RightSibling != nil {
		xRight = xRight.RightSibling
	}

	xRight.RightSibling = yRight
	yRight.LeftMostSibling = x.LeftMostSibling
	yRight.Parent = xRight.Parent

	for yRight.RightSibling != nil {
		yRight = yRight.RightSibling
		yRight.LeftMostSibling = yRight.RightSibling
		yRight.Parent = x.Parent
	}
	return x
}

func (n *Node) set(token *Lex.Token) {
	n.Token = token
}

func makeNode(s string, t *Lex.Token) *Node {
	switch s {
	case "id", "intNum", "floatNum":
		return &Node{
			t.Type,
			t.Lexeme,
			t,
			nil,
			nil,
			nil,
			nil,
		}
	}
	return Node{
		t,
		nil,
		nil,
		nil,
		nil,
	}
}

// func affect(node noder) {
// 	node.(*IdNode).set("this")
// }
