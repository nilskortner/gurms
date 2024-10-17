package supportpkgs_test

import (
	"gurms/internal/supportpkgs/datastructures/copyonwriteslice"
	"testing"
	"time"
)

func TestCOWS(t *testing.T) {
	slice := copyonwriteslice.NewCopyOnWriteSlice[int]()

	for i := 0; i < 100; i++ {
		slice.Add(i)
	}

	go func() {
		for i := 0; i < 100; i++ {
			slice.Add(i)
			t.Log(i)
			time.Sleep(1 * time.Microsecond)
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			slice.List()
			t.Log(slice)
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			slice.List()
			t.Log(slice)
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			slice.List()
			t.Log(slice)
		}
	}()

	time.Sleep(3 * time.Second)

	t.Error()
}
