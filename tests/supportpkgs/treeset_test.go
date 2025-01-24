package supportpkgs_test

import (
	"gurms/internal/supportpkgs/datastructures/treeset"
	"testing"
)

func TestHeadSet(t *testing.T) {
	set := treeset.New[string](treeset.StringComparator)

	for i := 'a'; i < 'g'; i++ {
		set.Put(string(i))
	}

	t.Log(set.Keys())
	t.Log(set.HeadSet("i"))
	t.Log(set.HeadSetSize("i"))
	t.Log(set.HeadSet("c"))
	t.Log(set.HeadSetSize("c"))

}
