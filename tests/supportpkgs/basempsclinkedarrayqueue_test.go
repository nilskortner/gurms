package supportpkgs_test

import (
	mpsc "gurms/internal/supportpkgs/datastructures/mpscunboundedarrayqueue"
	"strconv"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	queue := mpsc.NewBaseMpscLinkedArrayQueue[string](1024)

	t.Log(*queue.BUFFER_CONSUMED)
	t.Log(*queue.JUMP)
	t.Log(queue.MovingBuffer)
	t.Log(queue.Capacity)
	t.Log(queue.GetBuffer())

	t.Errorf("test")
}

func TestOffer(t *testing.T) {
	queue := mpsc.NewBaseMpscLinkedArrayQueue[string](1024)

	for i := 0; i < 550; i++ {
		queue.Offer((strconv.Itoa(i)))
		t.Log(queue.LvProducerIndex())
	}

	buf := queue.GetCBuffer()

	t.Log(buf)
	for _, ptr := range buf {
		if ptr != nil {
			t.Log(*ptr.Load())
		}
	}

	queue.Offer("124235")

	buf = queue.GetBuffer()
	for _, ptr := range buf {
		if ptr != nil {
			t.Log(*ptr.Load())
		}
	}

	for i := 0; i < 512; i++ {
		queue.Offer("124235")
	}

	t.Log(queue.GetBuffer())

	t.Log(queue.GetMask())

	t.Errorf("test")
}

func TestMultiOffer(t *testing.T) {
	queue := mpsc.NewBaseMpscLinkedArrayQueue[string](1024)

	go func() {
		for i := 0; i < 500; i++ {
			queue.Offer("go2: " + strconv.Itoa(i))
			t.Log("go2: " + strconv.Itoa(i))
		}
	}()

	go func() {
		for i := 0; i < 500; i++ {
			queue.Offer("go1: " + strconv.Itoa(i))
			t.Log("go1: " + strconv.Itoa(i))
		}
	}()

	time.Sleep(2 * time.Second)
	buf := queue.GetBuffer()
	t.Log(buf)

	t.Errorf("test")
}
