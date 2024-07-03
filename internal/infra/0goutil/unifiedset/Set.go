package goset

type Set map[string]struct{}

func NewSet() Set {
	return make(Set)
}

func (set Set) Add(item string) {
	set[item] = struct{}{}
}

func (set Set) Remove(item string) {
	delete(set, item)
}

func (set Set) Contains(item string) bool {
	_, found := set[item]
	return found
}

func (set Set) Size() int {
	return len(set)
}
