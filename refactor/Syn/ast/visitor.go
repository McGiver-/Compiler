package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/McGiver-/Compiler/refactor/Syn/Sem"
)

type Visitor interface {
	visit(*Node) []error
}

type TableCreationVisitor struct {
}

func (node *Node) SetFcallType(t *Sem.Table) {
	if node == nil {
		return
	}
	if node.Type == "FCall" {
		id := node.GetChildren()[1].Value
		typ := getFunctionType(id, t)
		node.Entry.Typ = typ
	}
	for _, v := range node.GetChildren() {
		v.SetFcallType(t)
	}
}

func (node *Node) CheckReturnType() (errors []error) {
	if node == nil {
		return nil
	}
	if node.Type == "returnStat" {
		if node.Entry == nil || node.Entry.Typ == "" || node.Entry.Typ == node.Parent.Parent.Parent.Entry.Typ {
			goto startFor
		}
		entry := node.findEntry(node.Entry.Name)
		if entry != nil && entry.Typ == node.Parent.Parent.Parent.Entry.Typ {
			goto startFor
		} else {
			if node.Parent.Parent.Parent.Entry.Typ == "int " {
				_, err := strconv.Atoi(node.Entry.Typ)
				if err != nil {
					errors = append(errors, fmt.Errorf("%s returning invalid type", node.Token.Position))
				}
			} else {
				errors = append(errors, fmt.Errorf("%s returning invalid type", node.Token.Position))
			}
		}
	}
startFor:
	for _, v := range node.GetChildren() {
		errors = append(errors, v.CheckReturnType()...)
	}
	return errors
}

func (node *Node) CheckReturnParams() (errors []error) {
	if node == nil {
		return nil
	}
	if node.Type == "FCall" {
		aparams := len(node.GetChildren()[1].GetChildren())
		paramsNm := 0
		name := node.GetChildren()[0].Token.Lit
		entry := node.findEntry(name)
		if entry == nil || entry.Child == nil || entry.Child.Entries == nil {
			goto startFor
		}
		for _, v := range entry.Child.Entries {
			if v.Kind == "parameter" {
				paramsNm++
			}
		}
		if aparams != paramsNm {
			errors = append(errors, fmt.Errorf("%s incorrect number of parameter passed", node.Token.Position))
		}
	}
startFor:
	for _, v := range node.GetChildren() {
		errors = append(errors, v.CheckReturnParams()...)
	}
	return errors
}

func (node *Node) CheckCalledFuncDeclared() (errors []error) {
	if node == nil {
		return nil
	}
	if node.Type == "FCall" {
		idToken := node.GetChildren()[0]
		name := idToken.Token.Lit
		entry := node.findEntry(name)
		if entry == nil {
			errors = append(errors, fmt.Errorf("%s <%s> function called without having been declared", idToken.Token.Position, name))
		}
	}
	for _, v := range node.GetChildren() {
		errors = append(errors, v.CheckCalledFuncDeclared()...)
	}
	return errors
}

func (node *Node) CheckUndeclaredDataMemeber() (errors []error) {
	if node == nil {
		return nil
	}
	if node.Type == "forStat" {
		return errors
	}
	if node.Type == "DataMember" {
		idToken := node.GetChildren()[0]
		name := idToken.Token.Lit
		entry := node.findEntry(name)
		if entry == nil {
			errors = append(errors, fmt.Errorf("%s <%s> undeclared datamember", idToken.Token.Position, name))
		}
	}
	for _, v := range node.GetChildren() {
		errors = append(errors, v.CheckUndeclaredDataMemeber()...)
	}
	return errors
}

func (node *Node) CheckUndeclaredVariable() (errors []error) {
	if node == nil {
		return nil
	}
	if node.Type == "DataMember" {
		idToken := node.GetChildren()[0]
		name := idToken.Token.Lit
		entry := node.findEntry(name)
		if entry == nil {
			errors = append(errors, fmt.Errorf("%s <%s> undeclared datamember", idToken.Token.Position, name))
		}
	}
	for _, v := range node.GetChildren() {
		errors = append(errors, v.CheckUndeclaredDataMemeber()...)
	}
	return errors
}

func getFunctionType(id string, table *Sem.Table) string {
	for _, v := range table.Entries {
		if v.Kind == "class" {
			for _, e := range v.Child.Entries {
				if e.Name == id {
					return e.Typ
				}
			}
		}
	}
	return ""
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
	case "Stat":
		return append(errors, node.visitStat()...)
	case "ifStat":
		return append(errors, node.visitIfStat()...)
	case "Expr":
		return append(errors, node.visitExpr()...)
	case "RelExpr":
		return append(errors, node.visitRelExpr()...)
	case "RelOp":
		return append(errors, node.visitRelOp()...)
	case "ArithExpr":
		return append(errors, node.visitArithExpr()...)
	case "Term":
		return append(errors, node.visitTerm()...)
	case "returnStat":
		return append(errors, node.visitReturnStat()...)
	case "Factor":
		return append(errors, node.visitFactor()...)
	case "Var":
		return append(errors, node.visitVar()...)
	case "FCall":
		return append(errors, node.visitFCall()...)
	case "DataMember":
		return append(errors, node.visitDataMember()...)
	case "IndexList":
		return append(errors, node.visitIndexList()...)
	case "AddOp":
		return append(errors, node.visitAddOp()...)
	case "Num":
		return append(errors, node.visitNum()...)
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
			errors = append(errors, fmt.Errorf("%s multiply declared <%s>", v.Token.Position, v.Entry.Name))
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
					errors = append(errors, fmt.Errorf("%s <%s> method signature does not match declaration", v.Token.Position, matchedEntry.Name))
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

func (n *Node) visitStat() (errors []error) {
	if n == nil {
		return errors
	}
	n.Entry = n.GetChildren()[0].Entry
	n.Table = &Sem.Table{Name: n.Value}
	n.Table.MergeTable(n.GetChildren()[0].Table)
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

func (n *Node) visitReturnStat() (errors []error) {
	n.Table = &Sem.Table{Name: n.Value}
	n.Table.MergeTable(n.GetChildren()[0].Table)
	n.Entry = n.Table.Entries[0]
	return errors
}

func (n *Node) visitExpr() (errors []error) {
	n.Table = &Sem.Table{Name: n.Value}
	n.Table.MergeTable(n.GetChildren()[0].Table)
	n.Entry = n.Table.Entries[0]
	return errors
}

func (n *Node) visitRelExpr() (errors []error) {
	n.Table = &Sem.Table{Name: n.Value}
	n.Table.MergeTable(n.GetChildren()[0].Table)
	n.Entry = n.Table.Entries[0]
	return errors
}

func (n *Node) visitRelOp() (errors []error) {
	n.Table = &Sem.Table{Name: n.Value}
	n.Table.MergeTable(n.GetChildren()[1].Table)
	n.Table.MergeTable(n.GetChildren()[2].Table)
	return errors
}

func (n *Node) visitAddOp() (errors []error) {
	n.Table = &Sem.Table{Name: n.Value}
	n.Table.MergeTable(n.GetChildren()[1].Table)
	n.Table.MergeTable(n.GetChildren()[2].Table)
	n.Entry = n.Table.Entries[0]
	return errors
}

func (n *Node) visitArithExpr() (errors []error) {
	n.Table = &Sem.Table{Name: n.Value}

	n.Table.MergeTable(n.GetChildren()[0].Table)
	if len(n.GetChildren()) > 1 {
		n.Table.MergeTable(n.GetChildren()[1].Table)
	}
	n.Entry = n.Table.Entries[0]
	return errors
}

func (n *Node) visitIndexList() (errors []error) {
	list := ""
	inh := n.GetChildren()
	for i := 0; i < len(inh); i++ {
		if i == len(inh)-1 {
			if inh[i].Entry.Typ != "" {
				index := strings.Trim(inh[i].Entry.Typ, " ")
				_, err := strconv.Atoi(index)
				if err != nil {
					errors = append(errors, fmt.Errorf("%s <%s> index must be an integer", n.Token.Position, index))
				}
			}
			list += fmt.Sprintf("%s", inh[i].Entry.Typ)
		} else {
			list += fmt.Sprintf("%s:", inh[i].Entry.Typ)
		}
	}
	n.Entry = Sem.NewEntry("inheritList", "inheritList", list, nil)
	return errors
}

func (n *Node) visitDataMember() (errors []error) {
	list := n.GetChildren()[0].Token.Lit + " "
	inh := n.GetChildren()
	for i := 0; i < len(inh); i++ {
		if i == len(inh)-1 {
			list += fmt.Sprintf("%s", inh[i].Entry.Typ)
		} else {
			if inh[i].Entry == nil {
				continue
			}
			list += fmt.Sprintf("%s:", inh[i].Entry.Typ)
		}
	}
	n.Entry = Sem.NewEntry("DataMember", "DataMember", list, nil)
	return errors
}

func (n *Node) visitFCall() (errors []error) {
	n.Entry = Sem.NewEntry(n.GetChildren()[0].Token.Lit, "FCall", "", nil)
	return errors
}

func (n *Node) visitTerm() (errors []error) {
	n.Table = &Sem.Table{Name: n.Value}
	n.Table.MergeTable(n.GetChildren()[0].Table)
	n.Entry = n.Table.Entries[0]
	return errors
}

func (n *Node) visitFactor() (errors []error) {
	n.Table = &Sem.Table{Name: n.Value}
	n.Table.MergeTable(n.GetChildren()[0].Table)
	n.Entry = n.Table.Entries[0]
	return errors
}

func (n *Node) visitNum() (errors []error) {
	n.Entry = &Sem.Entry{n.Value, "num", n.GetChildren()[0].Token.Lit, n.Table}
	n.Table = &Sem.Table{Name: n.Value}
	n.Table.Entries = append(n.Table.Entries, n.Entry)
	return errors
}

func (n *Node) visitVar() (errors []error) {
	n.Table = &Sem.Table{Name: n.Value}
	typ := "int"
	name := ""
	for _, v := range n.GetChildren() {
		if v.Entry != nil {
			typ = v.Entry.Typ
			name = v.Entry.Name
		}
	}
	n.Entry = &Sem.Entry{name, typ, "", nil}
	n.Table.Entries = append(n.Table.Entries, n.Entry)
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
			errors = append(errors, fmt.Errorf("%s <%s> multiply declared", v.Token.Position, v.Entry.Name))
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
			errors = append(errors, fmt.Errorf("%s multiply declared <%s>", fparam.Token.Position, fparam.Entry.Name))
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
		if stat.Type == "Stat" {
			continue
		}
		ok := nameMap[stat.Entry.Name]
		if ok {
			errors = append(errors, fmt.Errorf("%s multiply declared <%s>", stat.Token.Position, stat.Entry.Name))
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
			errors = append(errors, fmt.Errorf("%s multiply declared <%s>", fparam.Token.Position, fparam.Entry.Name))
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

func (node *Node) findTable(name string) *Sem.Table {
	current := node.Parent
	for current != nil {
		current = current.Parent
	}
	return current.Table.FindTable(name)
}

func (node *Node) findEntry(name string) *Sem.Entry {
	current := node
	for current.Parent != nil {
		current = current.Parent
	}
	return current.Table.FindEntry(name)
}
