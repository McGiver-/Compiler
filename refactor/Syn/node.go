package Syn

import (
	"fmt"

	"github.com/McGiver-/Compiler/refactor/Lex"
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

func (n *Node) makeFamily(kids ...*Node) {
	if len(kids) > 1 {
		kids[0].makeSiblings(kids[1:]...)
		n.adoptChildren(kids[0])
		n.Token = n.LeftMostChild.Token
		return
	}
	n.adoptChildren(kids[0])
	n.Token = n.LeftMostChild.Token
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

func (n *Node) GetChildLink(names ...string) []*Node {
	var children []*Node
	child := n.LeftMostChild
	for child != nil {
		if match(child.Type, names) {
			children = append(children, child)
		}
		child = child.RightSibling
	}
	return children
}

func match(s string, ss []string) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}
	return false
}

//AsumeFuncDefNode
func (n *Node) GetFuncVars() ([]string, []string, error) {
	var names []string
	var types []string
	statBlock := n.GetChild("StatBlock")
	vars := statBlock.GetChildLink("VarDecl")

	if len(vars) == 0 {
		return nil, nil, fmt.Errorf("could not find vars")
	}

	for _, variable := range vars {
		varname := variable.GetChild("id")
		t := variable.GetChild("Type")
		_type := t.Token.Lit + " "
		dimlist := variable.GetChild("DimList")
		dims := dimlist.GetChildLink("intNum")
		for _, dim := range dims {
			_type += dim.Value + " "
		}
		types = append(types, _type)
		names = append(names, varname.Token.Lit)
	}

	if len(names) == 0 {
		return nil, nil, fmt.Errorf("could not find vars")
	}

	return names, types, nil
}

func (n *Node) PrintChildren() string {
	toPrint := ""
	child := n.LeftMostChild
	for child != nil {
		if child.Type == "EPSILON" || child.Token == nil {
			toPrint += fmt.Sprintf("child : %s, ", child.Type)
		} else {
			toPrint += fmt.Sprintf("child : %s, ", child.Token.Lit)
		}
		child = child.RightSibling
	}
	return toPrint
}
