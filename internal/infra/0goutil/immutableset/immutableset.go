package immutableset

type ImmutableSet struct {
	iset map[any]struct{}
}

func NewImmutableSet[K comparable](values []K) *ImmutableSet {
	var i = &ImmutableSet{}
	for _, key := range values {
		i.iset[key] = struct{}{}
	}
	return i
}

func (i *ImmutableSet) Size() int {
	return len(i.iset)
}

func (i *ImmutableSet) ContainsKey(key any) bool {
	_, ok := i.iset[key]
	return ok
}
