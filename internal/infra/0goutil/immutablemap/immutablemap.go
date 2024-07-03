package immutablemap

import "gurms/internal/infra/0goutil/entry"

type ImmutableMap struct {
	imap map[any]any
}

func NewImmutableMap(values []entry.Entry) *ImmutableMap {
	var i = &ImmutableMap{}
	for _, entry := range values {
		i.imap[entry.GetKey()] = entry.GetValue()
	}
	return i
}

func (i *ImmutableMap) Size() int {
	return len(i.imap)
}

func (i *ImmutableMap) GetValue(key any) any {
	return i.imap[key]
}

func (i *ImmutableMap) ContainsKey(key any) bool {
	_, ok := i.imap[key]
	return ok
}

func (i *ImmutableMap) GetCopyAsMap() map[any]any {
	copy := make(map[any]any)
	for key, value := range i.imap {
		copy[key] = value
	}
	return copy
}

func GetValuesAsSlice() {

}
