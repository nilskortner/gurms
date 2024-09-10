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
		t.Log(queue.TestingLvProducerIndex())
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

func TestRelaxedPoll(t *testing.T) {
	queue := mpsc.NewBaseMpscLinkedArrayQueue[string](1024)

	for j := 0; j < 206; j++ {
		queue.Offer(("first"))
	}

	for i := 0; i < 600; i++ {
		queue.Offer((strconv.Itoa(i)))
		t.Log(queue.GetBuffer())
		abc, success := queue.RelaxedPoll()
		t.Log(queue.GetCBuffer())
		t.Log(abc)
		t.Log(success)
	}

	t.Errorf("test")
}

func TestMultiOfferRelaxedPoll(t *testing.T) {
	queue := mpsc.NewBaseMpscLinkedArrayQueue[string](1024)

	length := 0
	polllength := 0

	for j := 0; j < 11; j++ {
		//k := j
		go func() {
			for i := 0; i < 500; i++ {
				queue.Offer("go2: " + strconv.Itoa(i))
				//t.Log("go" + strconv.Itoa(k) + ": " + strconv.Itoa(i))
				//length++
			}
		}()
	}
	go func() {
		for i := 0; i < 1000; i++ {
			time.Sleep(5 * time.Millisecond)
			abc, success := queue.RelaxedPoll()
			t.Log("poll: " + abc)
			t.Log(success)
			t.Log(queue.TestingLvConsumerIndex())
			t.Log(queue.GetCBuffer())
			polllength++
		}
	}()

	time.Sleep(5 * time.Second)
	t.Log(length)
	t.Log(polllength)
	t.Errorf("testmultiofferrelaxed")
}
