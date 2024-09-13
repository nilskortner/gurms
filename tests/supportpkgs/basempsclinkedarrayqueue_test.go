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
	queue := mpsc.NewBaseMpscLinkedArrayQueue[string](512)

	length := 0
	polllength := 0

	early := queue.TestingGetMovingBuffer()
	earlyBuf := queue.TestingGetMovingBufferData(early)

	for j := 0; j < 500; j++ {
		k := j
		go func() {
			for i := 0; i < 1000; i++ {
				queue.Offer("go" + strconv.Itoa(k) + ": " + strconv.Itoa(i))
				length++
			}
		}()
	}
	time.Sleep(10 * time.Second)

	t.Log("first buffer")

	count := 0
	countf := 0
	go func() {
		for i := 0; i < 500000; i++ {
			_, success := queue.RelaxedPoll()
			// t.Log("poll: " + abc)
			if success == true {
				count++
			} else {
				countf++
			}
			if i < 3000 {

				// cData := queue.GetCBuffer()
				// values := make([]string, len(cData))
				// for i, ptr := range cData {
				// 	if ptr != nil {
				// 		values[i] = *ptr.Load()
				// 	} else {
				// 		values[i] = "nil"
				// 	}
				// }
				// t.Log(values)

				polllength++
			}
		}
	}()

	time.Sleep(3 * time.Second)

	// data, data2, data3, data4, data5 := queue.TestingGetMultiMovingBufferData()
	// values := make([]string, len(data))
	// for i, ptr := range data {
	// 	if ptr != nil {
	// 		values[i] = *ptr.Load()
	// 	} else {
	// 		values[i] = "nil"
	// 	}
	// }
	// values2 := make([]string, len(data2))
	// for i, ptr := range data2 {
	// 	if ptr != nil {
	// 		values2[i] = *ptr.Load()
	// 	} else {
	// 		values2[i] = "nil"
	// 	}
	// }
	// values3 := make([]string, len(data3))
	// for i, ptr := range data3 {
	// 	if ptr != nil {
	// 		values3[i] = *ptr.Load()
	// 	} else {
	// 		values3[i] = "nil"
	// 	}
	// }
	// values4 := make([]string, len(data4))
	// for i, ptr := range data4 {
	// 	if ptr != nil {
	// 		values4[i] = *ptr.Load()
	// 	} else {
	// 		values4[i] = "nil"
	// 	}
	// }
	// values5 := make([]string, len(data5))
	// for i, ptr := range data5 {
	// 	if ptr != nil {
	// 		values5[i] = *ptr.Load()
	// 	} else {
	// 		values5[i] = "nil"
	// 	}
	// }
	valuesEarly := make([]string, len(earlyBuf))
	for i, ptr := range earlyBuf {
		if ptr != nil {
			valuesEarly[i] = *ptr.Load()
		} else {
			valuesEarly[i] = "nil"
		}
	}
	// t.Log(values)
	// t.Log(len(values))
	// t.Log(values2)
	// t.Log(len(values2))
	// t.Log(values3)
	// t.Log(len(values3))
	// t.Log(values4)
	// t.Log(len(values4))
	// t.Log(values5)
	// t.Log(len(values5))
	t.Log(valuesEarly)

	time.Sleep(1 * time.Second)
	//t.Log(queue.TestingGetMovingBuffer())
	t.Log(length)
	t.Log(polllength)
	t.Log("successful: " + strconv.Itoa(count))
	t.Log("unsuccessful: " + strconv.Itoa(countf))
	t.Errorf("testmultiofferrelaxed")
}
