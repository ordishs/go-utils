package lockfree

import (
	"sync/atomic"
	"unsafe"
)

// Queue is a lock-free unbounded queue.
type Queue[T any] struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

type node struct {
	value interface{}
	next  unsafe.Pointer
}

// NewQueue returns an empty queue.
func NewQueue[T any]() *Queue[T] {
	n := unsafe.Pointer(&node{})
	return &Queue[T]{head: n, tail: n}
}

// Enqueue puts the given value v at the tail of the queue.
func (q *Queue[T]) Enqueue(v T) {
	// Create a node for the new value
	n := &node{value: v}

	// Now spin until we can update the tail
	for {
		tail := load[T](&q.tail)
		next := load[T](&tail.next)

		if tail == load[T](&q.tail) { // are tail and next consistent?
			if next == nil {
				if compareAndSwap(&tail.next, next, n) {
					compareAndSwap(&q.tail, tail, n) // Enqueue is done.  try to swing tail to the inserted node
					return
				}
			} else { // tail was not pointing to the last node
				// try to swing Tail to the next node
				compareAndSwap(&q.tail, tail, next)
			}
		}
	}
}

func (q *Queue[T]) Dequeue() interface{} {
	return q.dequeue()
}

func (q *Queue[T]) DequeueAsType() (T, bool) {
	v := q.dequeue()
	if v == nil {
		var t T
		return t, false
	}
	return v.(T), true
}

// Dequeue removes and returns the value at the head of the queue.
// It returns a default value of T and false if the queue is empty.
func (q *Queue[T]) dequeue() interface{} {
	for {
		head := load[T](&q.head)
		tail := load[T](&q.tail)
		next := load[T](&head.next)

		if head == load[T](&q.head) { // are head, tail, and next consistent?
			if head == tail { // is queue empty or tail falling behind?
				if next == nil { // is queue empty?
					return nil // queue is empty, couldn't dequeue
				}
				// tail is falling behind.  try to advance it
				compareAndSwap(&q.tail, tail, next)
			} else {
				// read value before CAS otherwise another dequeue might free the next node
				v := next.value

				if compareAndSwap(&q.head, head, next) {
					return v // Dequeue is done.  return
				}
			}
		}
	}
}

func load[T any](p *unsafe.Pointer) (n *node) {
	return (*node)(atomic.LoadPointer(p))
}

// Perform an atomic compare-and-swap operation on a pointer.
// This operation is used in concurrent programming to ensure that
// a value is updated only if it has not been modified by another
// thread since it was last observed.
func compareAndSwap(p *unsafe.Pointer, old, new *node) (ok bool) {
	return atomic.CompareAndSwapPointer(
		p, unsafe.Pointer(old), unsafe.Pointer(new))
}
