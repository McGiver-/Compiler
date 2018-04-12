package ast

import (
	"fmt"

	"github.com/McGiver-/Compiler/refactor/Syn/Sem"
)

type Visitor interface {
	visit(*Node)
}

type TableCreationVisitor struct {
}

func (visitor *TableCreationVisitor) visit(node *Node) {

	switch node.Type {
	case "Prog":
		node.visitProg()
	case "ClassDecl":
		node.visitClassDecl()
	default:
		node.visitNone()
	}
}

func (n *Node) Accept(visitor Visitor) {
	switch n.Type {
	case "Prog":
		n.acceptProg(visitor)
	case "ClassDecl":
		n.acceptClassDecl(visitor)
	default:
		n.acceptGeneric(visitor)
	}
}

func (n *Node) visitNone() {
}

func (n *Node) visitProg() {
	n.Table = &Sem.Table{Name: "global"}

	for _, v := range n.GetChildren()[0].GetChildren() {
		n.Table.AddEntry(v.Entry)
	}

	for _, v := range n.GetChildren()[1].GetChildren() {
		n.Table.AddEntry(v.Entry)
	}

	n.Table.AddEntry(Sem.NewEntry("program", "function", "", n.GetChildren()[2].Table))
}

func (n *Node) visitClassDecl() {
	className := n.GetChildren()[0].Value
	n.Table = &Sem.Table{Name: className}
	list := ""
	inh := n.GetChildren()[1].GetChildren()
	for i := 0; i < len(inh); i++ {
		if i == len(inh)-1 {
			list += fmt.Sprintf("%s", inh[i].Value)
		} else {
			list += fmt.Sprintf("%s:", inh[i].Value)
		}
	}
	for _, v := range n.GetChildren()[2].GetChildren() {
		n.Table.AddEntry(v.Entry)
	}
	n.Entry = Sem.NewEntry(className, "class", list, n.Table)
}

func (n *Node) acceptGeneric(visitor Visitor) {
	for _, child := range n.GetChildren() {
		child.Accept(visitor)
	}
	visitor.visit(n)
}

func (n *Node) acceptProg(visitor Visitor) {
	n.acceptGeneric(visitor)
}

func (n *Node) acceptClassDecl(visitor Visitor) {
	n.acceptGeneric(visitor)
}
