package utils

import "testing"

func TestLockFreeQueue(t *testing.T) {
	q := NewLockFreeQueue[int]()

	q.Enqueue(1)
	q.Enqueue(2)

	t.Log(q.Dequeue())
	t.Log(q.Dequeue())
	t.Log(q.Dequeue())
	t.Log(q.Dequeue())

	q2 := NewLockFreeQueue[string]()
	q2.Enqueue("a")
	q2.Enqueue("b")

	t.Log(q2.Dequeue())
	t.Log(q2.Dequeue())
	t.Log(q2.Dequeue())
}
