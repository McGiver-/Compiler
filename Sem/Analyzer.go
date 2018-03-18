package Sem

import (
	"fmt"

	"github.com/McGiver-/Compiler/Syn"
	"github.com/davecgh/go-spew/spew"
)

type Analyzer struct {
	rootNode *Syn.Node
	table    SymbolTable
}

func CreateAnalyzer(rootNode *Syn.Node) *Analyzer {
	return &Analyzer{rootNode, &table{}}
}

type SymbolTable interface {
	addEntry(string, string, string, bool)
	getEntries() []*entry
	populateFuncTable(names, kinds, types []string) error
}

type table struct {
	parent  SymbolTable
	entries []*entry
}

type entry struct {
	name  string
	kind  string
	_type string
	child SymbolTable
}

func (t *table) addEntry(name, kind, _type string, make bool) {
	if make {
		childTable := &table{parent: t}
		t.entries = append(t.entries, &entry{name, kind, _type, childTable})
	} else {
		t.entries = append(t.entries, &entry{name, kind, _type, nil})
	}
}

func (t *table) getEntries() []*entry {
	return t.entries
}

func (a *Analyzer) CreateTables() error {
	a.createTables()
	a.addFreeFuncs()
	a.addFuncs()
	spew.Dump(a.table)
	return nil
}

func (a *Analyzer) addFuncs() error {
	funcDefList, err := a.rootNode.GetChildLink("FuncDef")
	if err != nil {
		return fmt.Errorf("could not find funcDefList")
	}
	for _, funcDef := range funcDefList {
		id, err := funcDef.GetChild("id")
		funcName := id.Token.Lexeme
		scope, err := funcDef.GetChild("Scope")
		if err != nil {
			return fmt.Errorf("could not find scope")
		}
		classBelongs := scope.Token.Lexeme
		retType := funcDef.Token.Lexeme
		paramNames := make([]string, 0)
		types := make([]string, 0)

		for _, entry := range a.table.getEntries() {
			if entry.name != classBelongs || entry.kind != "class" {
				continue
			}
			classTable := entry.child
			fplist, err := funcDef.GetChild("FparamList")
			if err != nil {
				classTable.addEntry(funcName, "function", retType, true)
				continue
			}
			fpMembers, err := fplist.GetChildLink("FparamMember")
			if err != nil {
				classTable.addEntry(funcName, "function", retType, true)
				continue
			}
			for _, member := range fpMembers {
				t := member.Token.Type + " "
				dimlist, _ := member.GetChild("DimList")
				dims, _ := dimlist.GetChildLink("intNum")
				for _, dim := range dims {
					t += dim.Value + " "
				}
				n, _ := member.GetChild("id")
				paramNames = append(paramNames, n.Token.Lexeme)
				types = append(types, t)
			}
			classTable.addEntry(funcName, "function", retType, true)
			for _, x := range classTable.getEntries() {
				if x.name == funcName {
					x.child.populateFuncTable(paramNames, paramNames, types)
				}
			}
		}
	}
	return nil
}

func (t *table) populateFuncTable(names, kinds, types []string) error {
	for k := range names {
		t.addEntry(names[k], "parameter", types[k], false)
	}
	return nil
}

func (a *Analyzer) addFreeFuncs() error {
	funcDefList, err := a.rootNode.GetChildLink("FuncDef")
	if err != nil {
		return fmt.Errorf("could not find funcDef")
	}
	for _, funcDef := range funcDefList {
		scope, err := funcDef.GetChild("Scope")
		if err != nil {
			return fmt.Errorf("could not find scope")
		}
		if scope.Value != "EPSILON" {
			continue
		}
		kind := "function"
		name := scope.Token.Lexeme
		x, err := funcDef.GetChild("Type")
		if err != nil {
			return fmt.Errorf("could not find Type")
		}
		retType := x.Token.Type
		paramNames := make([]string, 0)
		types := make([]string, 0)

		fplist, err := funcDef.GetChild("FparamList")
		if err != nil {
			a.table.addEntry(name, kind, retType, true)
			a.populateFreeFuncTables(name, paramNames, paramNames, types)
			// add the func with no params and no new table TODO************

			return nil
		}

		fpMembers, err := fplist.GetChildLink("FparamMember")
		if err != nil {
			a.table.addEntry(name, kind, retType, true)
			a.populateFreeFuncTables(name, paramNames, paramNames, types)
			// add the func with no params and no new table TODO************
			return nil
		}

		for _, member := range fpMembers {
			t := member.Token.Type + " "
			dimlist, _ := member.GetChild("DimList")
			dims, _ := dimlist.GetChildLink("intNum")
			for _, dim := range dims {
				t += dim.Value + " "
			}
			n, _ := member.GetChild("id")
			paramNames = append(paramNames, n.Token.Lexeme)
			types = append(types, t)
		}
		a.table.addEntry(name, kind, retType, true)
		a.populateFreeFuncTables(name, paramNames, paramNames, types)
	}
	return nil
}

func (a *Analyzer) populateFreeFuncTables(funcName string, names, kind, types []string) {
	for _, entry := range a.table.getEntries() {
		if entry.kind != "function" || entry.name != funcName {
			continue
		}

		newT := entry.child
		if len(types) == 0 {
			continue
		}
		for k := range names {
			newT.addEntry(names[k], "parameter", types[k], false)
		}
	}
}

func (a *Analyzer) createTables() error {
	classList, err := a.rootNode.GetChildLink("ClassMember")
	if err == nil {
		for _, class := range classList {
			_type := ""
			name := class.Token.Lexeme
			InheritListParent, nil := class.GetChild("InheritList")
			list, err := InheritListParent.GetChildLink("InheritListMember")
			if err == nil {
				for k := range list {
					if k == len(list)-1 {
						_type += list[k].Token.Lexeme
						continue
					}
					_type += list[k].Token.Lexeme + " "
				}
			}
			a.table.addEntry(name, "class", _type, true)
		}
	}
	return nil
}
