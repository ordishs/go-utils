package lockfree

import (
	"sync"
)

// CQueue is a concurrent unbounded queue which uses two-Lock concurrent queue algorithm.
type CQueue[T any] struct {
	head     *cnode
	tail     *cnode
	headLock sync.Mutex
	tailLock sync.Mutex
}

type cnode struct {
	value interface{}
	next  *cnode
}

// NewCQueue returns an empty CQueue.
func NewCQueue[T any]() *CQueue[T] {
	n := &cnode{}
	return &CQueue[T]{head: n, tail: n}
}

// Enqueue puts the given value v at the tail of the queue.
func (q *CQueue[T]) Enqueue(v T) {
	n := &cnode{value: v}
	q.tailLock.Lock()
	q.tail.next = n // Link node at the end of the linked list
	q.tail = n      // Swing Tail to node
	q.tailLock.Unlock()
}

func (q *CQueue[T]) Dequeue() interface{} {
	return q.dequeue()
}

func (q *CQueue[T]) DequeueAsType() (T, bool) {
	v := q.dequeue()
	if v == nil {
		var t T
		return t, false
	}
	return v.(T), true
}

// Dequeue removes and returns the value at the head of the queue.
// It returns nil if the queue is empty.
func (q *CQueue[T]) dequeue() interface{} {
	q.headLock.Lock()
	n := q.head
	newHead := n.next
	if newHead == nil {
		q.headLock.Unlock()
		return nil
	}
	v := newHead.value
	newHead.value = nil
	q.head = newHead
	q.headLock.Unlock()
	return v
}
