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
