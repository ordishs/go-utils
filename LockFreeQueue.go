package utils

import (
	"sync/atomic"
	"unsafe"
)

// LockFreeQueue is a lock-free unbounded queue.
type LockFreeQueue[T any] struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

type node[T any] struct {
	value T
	next  unsafe.Pointer
}

// NewLockFreeQueue returns an empty queue.
func NewLockFreeQueue[T any]() *LockFreeQueue[T] {
	n := unsafe.Pointer(&node[T]{})
	return &LockFreeQueue[T]{head: n, tail: n}
}

// Enqueue puts the given value v at the tail of the queue.
func (q *LockFreeQueue[T]) Enqueue(v T) {
	// Create a node for the new value
	n := &node[T]{value: v}

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

// Dequeue removes and returns the value at the head of the queue.
// It returns a default value of T and false if the queue is empty.
func (q *LockFreeQueue[T]) Dequeue() (T, bool) {
	for {
		head := load[T](&q.head)
		tail := load[T](&q.tail)
		next := load[T](&head.next)

		if head == load[T](&q.head) { // are head, tail, and next consistent?
			if head == tail { // is queue empty or tail falling behind?
				if next == nil { // is queue empty?
					var t T
					return t, false // queue is empty, couldn't dequeue
				}
				// tail is falling behind.  try to advance it
				compareAndSwap(&q.tail, tail, next)
			} else {
				// read value before CAS otherwise another dequeue might free the next node
				v := next.value

				if compareAndSwap[T](&q.head, head, next) {
					return v, true // Dequeue is done.  return
				}
			}
		}
	}
}

func load[T any](p *unsafe.Pointer) (n *node[T]) {
	return (*node[T])(atomic.LoadPointer(p))
}

// Perform an atomic compare-and-swap operation on a pointer.
// This operation is used in concurrent programming to ensure that
// a value is updated only if it has not been modified by another
// thread since it was last observed.
func compareAndSwap[T any](p *unsafe.Pointer, old, new *node[T]) (ok bool) {
	return atomic.CompareAndSwapPointer(
		p, unsafe.Pointer(old), unsafe.Pointer(new))
}
