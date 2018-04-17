package Sem

import "fmt"

type Table struct {
	Name    string
	Parent  *Table
	Entries []*Entry
}

type Entry struct {
	Name  string
	Kind  string
	Typ   string
	Child *Table
}

func NewEntry(name, kind, typ string, child *Table) *Entry {
	return &Entry{
		Name:  name,
		Kind:  kind,
		Typ:   typ,
		Child: child,
	}
}

func (t *Table) AddEntry(entry *Entry) {
	if entry.Child != nil {
		entry.Child.Parent = t
	}
	t.Entries = append(t.Entries, entry)
}

func (t *Table) GetShadows() (errors []error) {
	for _, v := range t.Entries {
		if t.IsShadowing(v) {
			errors = append(errors, fmt.Errorf("<%s> is a shadowing member", v.Name))
		}
		if v.Child != nil {
			errors = append(errors, v.Child.GetShadows()...)
		}
	}
	return errors
}

func (t *Table) CheckMemberFuncNoDef() (errors []error) {
	classTables := []*Table{}
	for _, entry := range t.Entries {
		if entry.Kind == "class" {
			classTables = append(classTables, entry.Child)
		}
	}

	for _, classTable := range classTables {
		for _, entry := range classTable.Entries {
			if entry.Kind == "function" {
				hasVars := false
				for _, e := range entry.Child.Entries {
					if e.Kind == "variable" {
						hasVars = true
					}
				}
				if !hasVars {
					errors = append(errors, fmt.Errorf("<%s> member function has no definition", entry.Name))
				}
			}
		}
	}
	return errors
}

func (t *Table) IsShadowing(entry *Entry) bool {
	if t == nil {
		return false
	}
	parent := t.Parent
	for parent != nil {
		for _, v := range parent.Entries {
			if v.Name == entry.Name {
				return true
			}
		}
		parent = parent.Parent
	}
	return false
}

func (t *Table) MergeTable(t2 *Table) {
	if t2 == nil {
		return
	}
	for _, v := range t2.Entries {
		t.AddEntry(v)
	}
}

func (t *Table) FindTable(name string) *Table {
	if t == nil {
		return nil
	}
	table := t
	for _, v := range table.Entries {
		if v.Name == name {
			return v.Child
		}
		f := v.Child.FindTable(name)
		if f != nil {
			return f
		}
	}
	return nil
}

func (t *Table) FindEntry(name string) *Entry {
	if t == nil {
		return nil
	}
	table := t
	for _, v := range table.Entries {
		if v.Name == name {
			return v
		}
		f := v.Child.FindEntry(name)
		if f != nil {
			return f
		}
	}
	return nil
}
