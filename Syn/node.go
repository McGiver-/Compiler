package Syn

import (
	"github.com/McGiver-/Compiler/Lex"
)

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

func (n *Node) makeSiblings(y ...*Node) *Node {
	n.LeftMostSibling = n
	xRight := n
	yRight := y[0]

	for xRight.RightSibling != nil {
		xRight = xRight.RightSibling
	}

	for i := 0; i < len(y)-1; i++ {
		y[i].RightSibling = y[i+1]
	}

	xRight.RightSibling = yRight
	yRight.LeftMostSibling = n.LeftMostSibling
	yRight.Parent = xRight.Parent

	for yRight.RightSibling != nil {
		yRight = yRight.RightSibling
		yRight.LeftMostSibling = xRight.LeftMostSibling
		yRight.Parent = n.Parent
	}
	return n
}

func (n *Node) adoptChildren(y *Node) *Node {
	if n.LeftMostChild != nil {
		n.LeftMostChild.makeSiblings(y)
	} else {
		y.LeftMostSibling.Parent = n
		n.LeftMostChild = y.LeftMostSibling
		ysib := y.LeftMostSibling
		for ysib != nil {
			ysib.Parent = n
			ysib = ysib.RightSibling
		}
	}
	return n
}

func makeFamily(nodeType, lexeme string, op *Lex.Token, kids ...*Node) *Node {
	return makeNode(nodeType, lexeme, op).adoptChildren(kids[0].makeSiblings(kids[1:]...))
}

func (n *Node) set(token *Lex.Token) {
	n.Token = token
}

func makeNode(s, lexeme string, t *Lex.Token) *Node {
	node := &Node{
		s,
		lexeme,
		t, nil, nil, nil, nil,
	}
	node.LeftMostSibling = node
	return node
}
