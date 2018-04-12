package ast

import (
	"fmt"

	"github.com/McGiver-/Compiler/refactor/Lex"
	"github.com/McGiver-/Compiler/refactor/Syn/Sem"
)

type Node struct {
	Table           *Sem.Table
	Entry           *Sem.Entry
	Type            string
	Value           string
	Token           *Lex.Token
	Parent          *Node
	LeftMostChild   *Node
	LeftMostSibling *Node
	RightSibling    *Node
}

func (n *Node) makeSiblings(y ...*Node) {
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

	for yRight.RightSibling != nil {
		yRight.LeftMostSibling = xRight.LeftMostSibling
		yRight.Parent = n.Parent
		yRight = yRight.RightSibling
	}
}

func (n *Node) adoptChildren(y *Node) {
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
}

func (n *Node) MakeFamily(kids ...*Node) {
	if len(kids) > 1 {
		kids[0].makeSiblings(kids[1:]...)
		n.adoptChildren(kids[0])
		n.Token = n.LeftMostChild.Token
		return
	}
	n.adoptChildren(kids[0])
	n.Token = n.LeftMostChild.Token
}

func (n *Node) Set(token *Lex.Token) {
	n.Token = token
}

func MakeNode(s, lexeme string, t *Lex.Token) *Node {
	node := &Node{
		Type:  s,
		Value: lexeme,
		Token: t,
	}
	node.LeftMostSibling = node
	return node
}

func (n *Node) GetChild(name string) *Node {
	child := n.LeftMostChild
	for child != nil {
		if child.Type == name {
			return child
		}
		child = child.RightSibling
	}
	return &Node{}
}

func (n *Node) PrintChildren() string {
	toPrint := ""
	child := n.LeftMostChild
	for child != nil {
		if child.Type == "EPSILON" || child.Token == nil {
			toPrint += fmt.Sprintf("child : %s, ", child.Type)
		} else {
			toPrint += fmt.Sprintf("child : %s, ", child.Type)
		}
		child = child.RightSibling
	}
	return toPrint
}

func (n *Node) GetChildren() []*Node {
	var nodes []*Node
	child := n.LeftMostChild
	for child != nil {
		nodes = append(nodes, child)
		child = child.RightSibling
	}
	return nodes
}
