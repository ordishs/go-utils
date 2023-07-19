package lockfree

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLockFreeQueue(t *testing.T) {
	q := NewQueue[int]()
	require.NotNil(t, q)
	require.IsType(t, (*Queue[int])(nil), q)

	q.Enqueue(1)
	q.Enqueue(2)

	a := q.Dequeue()
	assert.Equal(t, 1, a)
	assert.IsType(t, (int(0)), a)

	b, ok := q.DequeueAsType()
	assert.Equal(t, 2, b)
	assert.IsType(t, (int(0)), b)
	assert.True(t, ok)

	c := q.Dequeue()
	assert.Nil(t, c)
	assert.IsType(t, nil, c)

	d, ok := q.DequeueAsType()
	assert.Equal(t, 0, d)
	assert.IsType(t, (int(0)), d)
	assert.False(t, ok)

	q2 := NewQueue[string]()
	require.NotNil(t, q2)
	require.IsType(t, (*Queue[string])(nil), q2)

	q2.Enqueue("a")
	q2.Enqueue("b")

	e := q2.Dequeue()
	assert.Equal(t, "a", e)
	assert.IsType(t, (string("")), e)

	f, ok := q2.DequeueAsType()
	assert.Equal(t, "b", f)
	assert.IsType(t, (string("")), f)
	assert.True(t, ok)

	g := q2.Dequeue()
	assert.Nil(t, g)
	assert.IsType(t, nil, g)

	h, ok := q2.DequeueAsType()
	assert.Equal(t, "", h)
	assert.IsType(t, (string("")), h)
	assert.False(t, ok)
}
