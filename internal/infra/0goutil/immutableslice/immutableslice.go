package immutableslice

type ImmutableSlice struct {
	islice []any
}

func NewImmutableSlice[K comparable](values []K) *ImmutableSlice {
	var i = &ImmutableSlice{}
	for count, value := range values {
		i.islice[count] = value
	}
	return i
}

func (i *ImmutableSlice) Size() int {
	return len(i.islice)
}

func (i *ImmutableSlice) ContainsValue(value any) bool {
	for count := range i.islice {
		if i.islice[count] == value {
			return true
		}
	}
	return false
}
