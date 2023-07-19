package lockfree

import "sync"

// SliceQueue is an unbounded queue which uses a slice as underlying.
type SliceQueue[T any] struct {
	data []interface{}
	mu   sync.Mutex
}

// NewSliceQueue returns an empty queue.
// You can give a initial capacity.
func NewSliceQueue[T any](n int) (q *SliceQueue[T]) {
	return &SliceQueue[T]{data: make([]interface{}, 0, n)}
}

// Enqueue puts the given value v at the tail of the queue.
func (q *SliceQueue[T]) Enqueue(v T) {
	q.mu.Lock()
	q.data = append(q.data, v)
	q.mu.Unlock()
}

func (q *SliceQueue[T]) Dequeue() interface{} {
	return q.dequeue()
}

func (q *SliceQueue[T]) DequeueAsType() (T, bool) {
	v := q.dequeue()
	if v == nil {
		var t T
		return t, false
	}
	return v.(T), true
}

// Dequeue removes and returns the value at the head of the queue.
// It returns nil if the queue is empty.
func (q *SliceQueue[T]) dequeue() interface{} {
	q.mu.Lock()
	if len(q.data) == 0 {
		q.mu.Unlock()
		return nil
	}
	v := q.data[0]
	q.data = q.data[1:]
	q.mu.Unlock()
	return v
}
