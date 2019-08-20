package grimoire

// TODO: assoc changes
// Use Assoc fields in Changes?
// Table name not stored here, but handled by repo logic.
// TODO: handle deleteion
//	- Answer: Changes should be forward only operation, no delete change is supported (use changeset instead).
type Changes struct {
	Fields       map[string]int // TODO: not copy friendly
	Changes      []Change
	Assoc        map[string]int
	AssocChanges [][]Changes
}

func (c Changes) Empty() bool {
	return len(c.Changes) == 0
}

func (c Changes) Get(field string) (Change, bool) {
	if index, ok := c.Fields[field]; ok {
		return c.Changes[index], true
	}

	return Change{}, false
}

func (c *Changes) Set(ch Change) {
	if index, exist := c.Fields[ch.Field]; exist {
		c.Changes[index] = ch
	} else {
		c.Fields[ch.Field] = len(c.Changes)
		c.Changes = append(c.Changes, ch)
	}
}

func (c Changes) GetValue(field string) (interface{}, bool) {
	var (
		ch, ok = c.Get(field)
	)

	return ch.Value, ok
}

func (c *Changes) SetValue(field string, value interface{}) {
	c.Set(Set(field, value))
}

func (c Changes) GetAssoc(field string) ([]Changes, bool) {
	if index, ok := c.Assoc[field]; ok {
		return c.AssocChanges[index], true
	}

	return nil, false
}

func (c *Changes) SetAssoc(field string, chs ...Changes) {
	if index, exist := c.Assoc[field]; exist {
		c.AssocChanges[index] = chs
	} else {
		c.Assoc[field] = len(c.AssocChanges)
		c.AssocChanges = append(c.AssocChanges, chs)
	}
}

type ChangeOp int

const (
	ChangeSetOp ChangeOp = iota
	ChangeIncOp
	ChangeDecOp
	ChangeFragmentOp
)

type Change struct {
	Type  ChangeOp
	Field string
	Value interface{}
}

func (c Change) Build(changes *Changes) {
	changes.Set(c)
}

func Set(field string, value interface{}) Change {
	return Change{
		Type:  ChangeSetOp,
		Field: field,
		Value: value,
	}
}

func Inc(field string) Change {
	return IncBy(field, 1)
}

func IncBy(field string, n int) Change {
	return Change{
		Type:  ChangeIncOp,
		Field: field,
		Value: n,
	}
}

func Dec(field string) Change {
	return DecBy(field, 1)
}

func DecBy(field string, n int) Change {
	return Change{
		Type:  ChangeDecOp,
		Field: field,
		Value: n,
	}
}

func ChangeFragment(raw string, args ...interface{}) Change {
	return Change{
		Type:  ChangeFragmentOp,
		Field: raw,
		Value: args,
	}
}