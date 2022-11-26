package utils

import (
	"time"
)

type Batcher[T any] struct {
	fn         func([]*T)
	size       int
	timeout    time.Duration
	batch      []*T
	ch         chan *T
	background bool
}

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
		if b.background {
			// var batch []*T
			// copy(batch, b.batch)
			go b.fn(b.batch)
		} else {
			b.fn(b.batch)
		}
		b.batch = b.batch[:0] // Clear batch but keep the allocated memory
		// b.batch = nil // Clear batch and clear allocated memory
	}
}
