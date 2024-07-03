package entry

type Entry struct {
	key   any
	value any
}

func NewEntry(key, value any) Entry {
	e := Entry{
		key:   key,
		value: value,
	}
	return e
}

func (e Entry) GetKey() any {
	return e.key
}

func (e Entry) GetValue() any {
	return e.value
}

func (e *Entry) SetValue(value any) {
	e.value = value
}
