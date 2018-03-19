package Sem

import (
	"fmt"
	"strings"

	"github.com/McGiver-/Compiler/Syn"
	"github.com/davecgh/go-spew/spew"
)

type Analyzer struct {
	rootNode *Syn.Node
	table    *table
}

func CreateAnalyzer(rootNode *Syn.Node) *Analyzer {
	return &Analyzer{rootNode, &table{}}
}

type table struct {
	parent  *table
	entries []*entry
}

type entry struct {
	name  string
	kind  string
	_type string
	child *table
}

func (t *table) findTable(entryName string) *table {
	if t == nil {
		return nil
	}
	for _, entry := range t.getEntries() {
		if entry.name == entryName {
			return entry.child
		} else {
			table := entry.child.findTable(entryName)
			if table != nil {
				return table
			}
		}
	}
	return nil
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

func (a *Analyzer) CreateTables() (errors []error) {
	a.createClasses()
	a.addFreeFuncs()
	a.addFuncs()
	errors = append(errors, a.inheritMembers()...)
	errors = append(errors, a.table.findDuplicatesInTables()...)
	spew.Dump(a.table)
	return removeDuplicates(errors)
}

func (t *table) findDuplicatesInTables() (errors []error) {
	currentTable := t
	if currentTable == nil {
		return nil
	}
	m := make(map[string]bool)
	for _, k := range currentTable.getEntries() {
		if m[k.name] {
			var parent *entry
			if currentTable.parent == nil {
				errors = append(errors, fmt.Errorf("%s: %s declared twice in global scope", k.kind, k.name))
				continue
			}
			for _, parentEntry := range currentTable.parent.getEntries() {
				if parentEntry.child == currentTable {
					parent = parentEntry
				}
			}
			errors = append(errors, fmt.Errorf("%s: %s declared twice in %s: %s", k.kind, k.name, parent.kind, parent.name))
		}
		m[k.name] = true
		errors = append(errors, k.child.findDuplicatesInTables()...)
	}
	return
}

func (a *Analyzer) inheritMembers() (errors []error) {
	for _, entry := range a.table.getEntries() {
		if entry.kind != "class" {
			continue
		}
		classtable := a.table.findTable(entry.name)
		for _, parentClass := range strings.Split(strings.TrimSpace(entry._type), " ") {
			if parentClass == "" {
				continue
			}
			if a.table.findTable(parentClass) == nil {
				errors = append(errors, fmt.Errorf("class :%s inheritied and not declared", parentClass))
				continue
			}
			for _, parentEntry := range a.table.findTable(parentClass).getEntries() {
				if hasDuplicateEntry(parentEntry, classtable.getEntries()) {
					errors = append(errors, fmt.Errorf("Warning shadowing %s from class %v in class %s", parentEntry.name, parentClass, entry.name))
					continue
				}
				classtable.entries = append(classtable.entries, parentEntry)
			}
		}
	}
	return
}

func hasDuplicateEntry(e *entry, entries []*entry) bool {
	for _, k := range entries {
		if e.name == k.name {
			return true
		}
	}
	return false
}

func (a *Analyzer) addFuncs() (errors []error) {
	funcDefList := a.rootNode.GetChildLink("FuncDef")
	paramNames := make([]string, 0)
	types := make([]string, 0)
	if len(funcDefList) == 0 {
		return append(errors, fmt.Errorf("could not find funcDefList"))
	}

	for _, funcDef := range funcDefList {
		retType := funcDef.Token.Lexeme
		scope := funcDef.GetChild("Scope")
		if scope.Value == "EPSILON" {
			continue
		}
		classBelongs := scope.Token.Lexeme
		classTable := a.table.findTable(classBelongs)

		id := funcDef.GetChild("id")
		funcName := id.Token.Lexeme

		fpMembers := funcDef.GetChild("FparamList").GetChildLink("FparamMember")
		if len(fpMembers) == 0 {
			classTable.addEntry(funcName, "function", retType, true)
			classTable.populateTableWithFuncVars(funcDef, funcName, paramNames, types)
			continue
		}
		for _, member := range fpMembers {
			t := member.Token.Type + " "
			dims := member.GetChild("DimList").GetChildLink("intNum")
			for _, dim := range dims {
				t += dim.Value + " "
			}
			paramNames = append(paramNames, member.GetChild("id").Token.Lexeme)
			types = append(types, t)
		}
		classTable.addEntry(funcName, "function", retType, true)
		classTable.populateTableWithFuncVars(funcDef, funcName, paramNames, types)
	}
	return nil
}

func (t *table) populateTableWithFuncVars(funcDef *Syn.Node, funcName string, paramNames, types []string) {
	for _, x := range t.getEntries() {
		if x.name == funcName {
			x.child.populateTable(paramNames, "parameter", types)
			paramNames, types, err := funcDef.GetFuncVars()
			if err != nil {
				continue
			}
			x.child.populateTable(paramNames, "variable", types)
		}
	}
}

func (t *table) populateTable(names []string, kind string, types []string) error {
	for k := range names {
		if len(types) < k {
			t.addEntry(names[k], kind, "", true)
		} else {
			t.addEntry(names[k], kind, types[k], true)
		}
	}
	return nil
}

func (a *Analyzer) addFreeFuncs() (errors []error) {
	funcDefList := a.rootNode.GetChildLink("FuncDef")
	if len(funcDefList) == 0 {
		return nil
	}

	for _, funcDef := range funcDefList {
		scope := funcDef.GetChild("Scope")
		if scope.Value != "EPSILON" {
			continue
		}
		name := scope.Token.Lexeme
		retType := funcDef.GetChild("Type").Token.Type
		paramNames := make([]string, 0)
		types := make([]string, 0)

		fplist := funcDef.GetChild("FparamList")
		if fplist.Token == nil {
			a.table.addEntry(name, "function", retType, true)
			for _, entry := range a.table.getEntries() {
				if entry.name == name {
					names, types, _ := funcDef.GetFuncVars()
					entry.child.populateTable(names, "variable", types)
				}
			}
			// a.populateFreeFuncTables(name, paramNames, paramNames, types)
			continue
		}

		fpMembers := fplist.GetChildLink("FparamMember")
		if len(fpMembers) == 0 {
			a.table.addEntry(name, "function", retType, true)
			for _, entry := range a.table.getEntries() {
				if entry.name == name {
					// entry.child.populateTable([]string{name}, "parameter", types)
					names, types, _ := funcDef.GetFuncVars()
					entry.child.populateTable(names, "variable", types)
				}
			}
			// a.populateFreeFuncTables(name, paramNames, paramNames, types)
			continue
		}

		for _, member := range fpMembers {
			t := member.Token.Type + " "
			dims := member.GetChild("DimList").GetChildLink("intNum")
			for _, dim := range dims {
				t += dim.Value + " "
			}
			paramNames = append(paramNames, member.GetChild("id").Token.Lexeme)
			types = append(types, t)
		}
		a.table.addEntry(name, "function", retType, true)
		for _, entry := range a.table.getEntries() {
			if entry.name == name {
				entry.child.populateTable([]string{name}, "parameter", types)
				names, types, _ := funcDef.GetFuncVars()
				entry.child.populateTable(names, "variable", types)
			}
		}
		// a.populateFreeFuncTables(name, paramNames, paramNames, types)

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

func (a *Analyzer) createClasses() (errors []error) {
	classList := a.rootNode.GetChildLink("ClassMember")
	if len(classList) == 0 {
		return errors
	}
	var _typeList []string
	var nameList []string
	// GetClasses
	for _, class := range classList {
		_type := ""
		name := class.Token.Lexeme
		list := class.GetChild("InheritList").GetChildLink("InheritListMember")
		for k := range list {
			if k == len(list)-1 {
				_type += list[k].Token.Lexeme
				continue
			}
			_type += list[k].Token.Lexeme + " "
		}
		_typeList = append(_typeList, _type)
		nameList = append(nameList, name)
		duplicateClasses := getDuplicates(nameList)
		for _, dup := range duplicateClasses {
			errors = append(errors, fmt.Errorf("Class %s declared twice", dup))
		}
	}
	//Add classes
	a.table.populateTable(nameList, "class", _typeList)

	//Get Variables
	for _, class := range classList {
		_typeList, nameList = []string{}, []string{}
		members := class.GetChildLink("MemberList")
		if len(members) == 0 {
			continue
		}
		for _, member := range members {
			vars := member.GetChildLink("VarDecl")
			if len(vars) == 0 {
				continue
			}
			for _, variable := range vars {
				_type := variable.GetChild("Type").Token.Lexeme + " "
				dims := variable.GetChild("DimList").GetChildLink("intNum")
				for _, dim := range dims {
					_type += dim.Value + " "
				}
				_typeList = append(_typeList, _type)
				nameList = append(nameList, variable.GetChild("id").Token.Lexeme)
				duplicateVariables := getDuplicates(nameList)
				for _, dup := range duplicateVariables {
					errors = append(errors, fmt.Errorf("variable %s declared twice in class %s", dup, class.Token.Lexeme))
				}
			}
		}

		a.table.findTable(class.Token.Lexeme).populateTable(nameList, "variable", _typeList)
	}
	return nil
}

func getDuplicates(a []string) (d []string) {
	m := make(map[string]bool)
	for _, v := range a {
		if m[v] {
			d = append(d, v)
			continue
		}
		m[v] = true
	}
	return
}

func removeDuplicates(a []error) (d []error) {
	m := make(map[string]bool)
	for _, v := range a {
		if m[v.Error()] {
			continue
		}
		d = append(d, v)
		m[v.Error()] = true
	}
	return
}
