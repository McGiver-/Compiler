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
	case "VarDecl":
		node.visitVarDecl()
	case "FuncDecl":
		node.visitFuncDecl()
	case "Fparam":
		node.visitFParam()
	case "FuncDef":
		node.visitFuncDef()
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

func (n *Node) visitFuncDef() {
	n.Table = &Sem.Table{Name: n.GetChildren()[2].Token.Lit}
	retVal := n.GetChildren()[0].GetChildren()[0].Value
	for _, fparam := range n.GetChildren()[3].GetChildren() {
		n.Table.AddEntry(fparam.Entry)
	}
	for _, stat := range n.GetChildren()[4].GetChildren() {
		n.Table.AddEntry(stat.Entry)
	}
	n.Entry = Sem.NewEntry(n.Table.Name, "function", retVal, n.Table)
}

func (n *Node) visitFuncDecl() {
	n.Table = &Sem.Table{Name: n.GetChildren()[1].Token.Lit}
	retVal := n.GetChildren()[0].GetChildren()[0].Value
	for _, fparam := range n.GetChildren()[2].GetChildren() {
		n.Table.AddEntry(fparam.Entry)
	}
	n.Entry = Sem.NewEntry(n.Table.Name, "function", retVal, n.Table)
}

func (n *Node) visitFParam() {
	list := n.GetChildren()[0].GetChildren()[0].Value + " "
	dims := n.GetChildren()[2].GetChildren()
	for i := 0; i < len(dims); i++ {
		if i == len(dims)-1 {
			list += fmt.Sprintf("%s", dims[i].Value)
		} else {
			list += fmt.Sprintf("%s:", dims[i].Value)
		}
	}
	n.Entry = Sem.NewEntry(n.GetChildren()[1].Token.Lit, "parameter", list, nil)
}

func (n *Node) visitVarDecl() {
	list := n.GetChildren()[0].GetChildren()[0].Value + " "
	dims := n.GetChildren()[2].GetChildren()
	for i := 0; i < len(dims); i++ {
		if i == len(dims)-1 {
			list += fmt.Sprintf("%s", dims[i].Value)
		} else {
			list += fmt.Sprintf("%s:", dims[i].Value)
		}
	}
	n.Entry = Sem.NewEntry(n.GetChildren()[1].Token.Lit, "variable", list, nil)
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
