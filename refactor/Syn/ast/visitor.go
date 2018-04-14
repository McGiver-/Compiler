package ast

import (
	"fmt"

	"github.com/McGiver-/Compiler/refactor/Syn/Sem"
)

type Visitor interface {
	visit(*Node) []error
}

type TableCreationVisitor struct {
}

func (visitor *TableCreationVisitor) visit(node *Node) (errors []error) {

	switch node.Type {
	case "Prog":
		return append(errors, node.visitProg()...)
	case "ClassDecl":
		return append(errors, node.visitClassDecl()...)
	case "VarDecl":
		return append(errors, node.visitVarDecl()...)
	case "StatBlock":
		return append(errors, node.visitStatBlock()...)
	case "FuncDecl":
		return append(errors, node.visitFuncDecl()...)
	case "Fparam":
		return append(errors, node.visitFParam()...)
	case "FuncDef":
		return append(errors, node.visitFuncDef()...)
	case "ifStat":
		return append(errors, node.visitIfStat()...)
	// case "forStat":
	// 	return append(errors, node.visitFuncDef()...)
	default:
		return append(errors, node.visitNone()...)
	}
}

func (n *Node) Accept(visitor Visitor) []error {
	switch n.Type {
	// case "Prog":
	// 	n.acceptProg(visitor)
	// case "ClassDecl":
	// 	n.acceptClassDecl(visitor)
	default:
		return n.acceptGeneric(visitor)
	}
}

func (n *Node) visitNone() (errors []error) {
	return errors
}

func (n *Node) visitProg() (errors []error) {
	n.Table = &Sem.Table{Name: "global"}
	nameMap := map[string]bool{}
	for _, v := range n.GetChildren()[0].GetChildren() {
		ok := nameMap[v.Entry.Name]
		if ok {
			errors = append(errors, fmt.Errorf("multiply declared <%s> at:%s", v.Entry.Name, v.Token.Position))
			continue
		}
		nameMap[v.Entry.Name] = true
		n.Table.AddEntry(v.Entry)
	}

	nameMap = map[string]bool{}
	for _, v := range n.GetChildren()[1].GetChildren() {
		ok := nameMap[v.Entry.Name]
		if ok {
			continue
		}
		class := classesContainFunc(n.Table.Entries, v.GetChildren()[1].Token.Lit)
		if class != nil {
			matchedEntry := entryContainsFunc(class.Child.Entries, v.Entry.Name)
			if matchedEntry != nil {
				if matchedEntry.Typ != v.Entry.Typ {
					errors = append(errors, fmt.Errorf("<%s> method signature does not match declaration at:%s\n", matchedEntry.Name, v.Token.Position))
				}
				matchedEntry.Child = v.Table
				v.Table.Parent = class.Child
				nameMap[v.Entry.Name] = true
				continue
			}
		}
		nameMap[v.Entry.Name] = true
		n.Table.AddEntry(v.Entry)
	}

	n.GetChildren()[2].Table = &Sem.Table{Name: "program"}
	// for _, v := range n.GetChildren()[2].GetChildren() {
	// 	if v == nil || v.Entry == nil {
	// 		continue
	// 	}
	// 	n.GetChildren()[2].Table.AddEntry(v.Entry)
	// }
	n.Table.AddEntry(Sem.NewEntry("program", "function", "", n.GetChildren()[2].Table))
	return errors
}

func (n *Node) visitStatBlock() (errors []error) {
	if n == nil {
		return errors
	}
	n.Entry = &Sem.Entry{n.Value, "StatBlock", "", n.Table}
	n.Table = &Sem.Table{Name: n.Value}
	for _, v := range n.GetChildren() {
		if v == nil || v.Entry == nil {
			continue
		}
		n.Table.AddEntry(v.Entry)
	}
	return errors
}

func (n *Node) visitIfStat() (errors []error) {
	n.Entry = &Sem.Entry{n.Value, "ifstat", "", n.Table}
	n.Table = &Sem.Table{Name: n.Value}
	n.Table.MergeTable(n.GetChildren()[0].Table)
	n.Table.AddEntry(n.GetChildren()[1].Entry)
	n.Table.AddEntry(n.GetChildren()[2].Entry)
	return errors
}

func (n *Node) visitExpr() (errors []error) {
	n.Table = &Sem.Table{Name: n.Value}
	n.Table.MergeTable(n.GetChildren()[0].Table)
	return errors
}

func (n *Node) visitRelExpr() (errors []error) {
	n.Table = &Sem.Table{Name: n.Value}
	n.Table.MergeTable(n.GetChildren()[1].Table)
	n.Table.MergeTable(n.GetChildren()[2].Table)
	return errors
}

func (n *Node) visitArithExpr() (errors []error) {
	n.Table = &Sem.Table{Name: n.Value}
	n.Table.MergeTable(n.GetChildren()[0].Table)
	if len(n.GetChildren()) > 1 {
		n.Table.MergeTable(n.GetChildren()[1].Table)
	}
	return errors
}

func (n *Node) visitTerm() (errors []error) {
	n.Table = &Sem.Table{Name: n.Value}
	n.Table.MergeTable(n.GetChildren()[0].Table)
	return errors
}

func (n *Node) visitFactor() (errors []error) {
	n.Table = &Sem.Table{Name: n.Value}
	n.Table.MergeTable(n.GetChildren()[0].Table)
	return errors
}

func (n *Node) visitVar() (errors []error) {
	n.Table = &Sem.Table{Name: n.Value}
	n.Table.MergeTable(n.GetChildren()[0].Table)
	return errors
}

func (n *Node) visitStat() (errors []error) {
	n.Entry = n.GetChildren()[0].Entry
	return errors
}

func (n *Node) visitClassDecl() (errors []error) {
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
	nameMap := map[string]bool{}
	for _, v := range n.GetChildren()[2].GetChildren() {
		ok := nameMap[v.Entry.Name]
		if ok {
			errors = append(errors, fmt.Errorf("multiply declared <%s> at:%s", v.Entry.Name, v.Token.Position))
			continue
		}
		nameMap[v.Entry.Name] = true
		n.Table.AddEntry(v.Entry)
	}
	n.Entry = Sem.NewEntry(className, "class", list, n.Table)
	return errors
}

func (n *Node) visitFuncDef() (errors []error) {
	n.Table = &Sem.Table{Name: n.GetChildren()[2].Token.Lit}
	retVal := n.GetChildren()[0].GetChildren()[0].Value + " "
	nameMap := map[string]bool{}
	for _, fparam := range n.GetChildren()[3].GetChildren() {
		ok := nameMap[fparam.Entry.Name]
		if ok {
			errors = append(errors, fmt.Errorf("multiply declared <%s> at:%s", fparam.Entry.Name, fparam.Token.Position))
			continue
		}
		nameMap[fparam.Entry.Name] = true
		retVal += fparam.GetChildren()[0].Value + " "
		for _, dim := range fparam.GetChildren()[2].GetChildren() {
			retVal += dim.Value + " "
		}
		n.Table.AddEntry(fparam.Entry)
	}
	nameMap = map[string]bool{}
	for _, stat := range n.GetChildren()[4].GetChildren() {
		ok := nameMap[stat.Entry.Name]
		if ok {
			errors = append(errors, fmt.Errorf("multiply declared <%s> at:%s", stat.Entry.Name, stat.Token.Position))
			continue
		}
		nameMap[stat.Entry.Name] = true
		n.Table.AddEntry(stat.Entry)
	}
	n.Entry = Sem.NewEntry(n.Table.Name, "function", retVal, n.Table)
	return errors
}

func (n *Node) visitFuncDecl() (errors []error) {
	n.Table = &Sem.Table{Name: n.GetChildren()[1].Token.Lit}
	retVal := n.GetChildren()[0].GetChildren()[0].Value + " "
	nameMap := map[string]bool{}
	for _, fparam := range n.GetChildren()[2].GetChildren() {
		ok := nameMap[fparam.Entry.Name]
		if ok {
			errors = append(errors, fmt.Errorf("multiply declared <%s> at:%s", fparam.Entry.Name, fparam.Token.Position))
			continue
		}
		retVal += fparam.GetChildren()[0].Value + " "
		for _, dim := range fparam.GetChildren()[2].GetChildren() {
			retVal += dim.Value + " "
		}
		nameMap[fparam.Entry.Name] = true
		n.Table.AddEntry(fparam.Entry)
	}
	n.Entry = Sem.NewEntry(n.Table.Name, "function", retVal, n.Table)
	return errors
}

func (n *Node) visitFParam() (errors []error) {
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
	return errors
}

func (n *Node) visitVarDecl() (errors []error) {
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
	return errors
}

func (n *Node) acceptGeneric(visitor Visitor) (errors []error) {
	for _, child := range n.GetChildren() {
		errors = append(errors, child.Accept(visitor)...)
	}
	errors = append(errors, visitor.visit(n)...)
	return errors
}

func (n *Node) acceptProg(visitor Visitor) {
	n.acceptGeneric(visitor)
}

func (n *Node) acceptClassDecl(visitor Visitor) {
	n.acceptGeneric(visitor)
}

/*------------------------------------------HELPERS--------------------------------*/

func classesContainFunc(classes []*Sem.Entry, f string) *Sem.Entry {
	for _, class := range classes {
		if class.Name == f {
			return class
		}
	}
	return nil
}

func entryContainsFunc(entries []*Sem.Entry, f string) *Sem.Entry {
	for _, entry := range entries {
		if entry.Name == f {
			return entry
		}
	}
	return nil
}

func matchFuncSignature(e1 []*Sem.Entry, e2 []*Sem.Entry) bool {
	params1 := []*Sem.Entry{}
	params2 := []*Sem.Entry{}
	for _, v := range e1 {
		if v.Kind == "parameter" {
			params1 = append(params1, v)
		}
	}
	for _, v := range e2 {
		if v.Kind == "parameter" {
			params2 = append(params2, v)
		}
	}
	if len(params1) != len(params2) {
		return false
	}

	for _, v := range params1 {
		found := false
		for _, k := range params2 {
			if k.Name == v.Name && k.Typ == v.Typ {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
