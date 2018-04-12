package Sem

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
	t.Entries = append(t.Entries, entry)
}
