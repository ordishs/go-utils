package utils

import (
	"time"
)

// Batcher is a utility that batches items together and then invokes the provided function
// on that whenever it reaches the specified size or the timeout is reached.
type Batcher[T any] struct {
	fn         func([]*T)
	size       int
	timeout    time.Duration
	batch      []*T
	ch         chan *T
	background bool
}

// New creates a new Batcher that will invoke the provided function when the batch size is reached.
// The size is the maximum number of items that can be batched before processing the batch.
// The timeout is the duration that will be waited before processing the batch.
func New[T any](size int, timeout time.Duration, fn func(batch []*T), background bool) *Batcher[T] {
	b := &Batcher[T]{
		fn:         fn,
		size:       size,
		timeout:    timeout,
		ch:         make(chan *T),
		background: background,
	}

	go b.worker()

	return b
}

// Put adds an item to the batch. If the batch is full, or the timeout is reached
// the batch will be processed.
func (b *Batcher[T]) Put(item *T) {
	b.ch <- item
}

func (b *Batcher[T]) worker() {
	for {
		expire := time.After(b.timeout)
		for {
			select {
			case item := <-b.ch:
				b.batch = append(b.batch, item)

				if len(b.batch) == b.size {
					goto saveBatch
				}

			case <-expire:
				goto saveBatch
			}
		}
	saveBatch:
		if len(b.batch) > 0 {
			if b.background {
				// var batch []*T
				// copy(batch, b.batch)
				go b.fn(b.batch)
			} else {
				b.fn(b.batch)
			}

			// De-reference the batch
			b.batch = nil
		}
	}
}
