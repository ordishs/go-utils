package batcher

import (
	"sync"
	"testing"
	"time"
)

func TestBatcher(t *testing.T) {
	var wg sync.WaitGroup

	fn := func(batch []*int) {
		t.Logf("Performing batch of size %d", len(batch))
		for range batch {
			wg.Done()
		}
	}

	b := New(200, 10*time.Millisecond, fn, true)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		b.Put(&i)
	}

	wg.Wait()
	t.Log("done")

	// time.Sleep(2 * time.Second)
}
