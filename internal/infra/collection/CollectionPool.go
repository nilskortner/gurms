package collection

import (
	"sync"
)

type Set map[interface{}]struct{}

var setPool = sync.Pool{
	New: func() interface{} {
		return make(Set, 64)
	},
}

func GetSet() Set {
	return setPool.Get().(Set)
}

func PutSet(s Set) {
	for k := range s {
		delete(s, k)
	}
	setPool.Put(s)
}
