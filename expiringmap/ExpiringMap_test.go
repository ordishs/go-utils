// cSpell: ignore expiringmap,Equalf,Nilf
package expiringmap

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringExpiringMap(t *testing.T) {
	m := New[string, string](100 * time.Millisecond)

	m.Set("foo", "bar")

	v, ok := m.Get("foo")
	assert.Equalf(t, ok, true, "expected to find key")
	assert.Equalf(t, v, "bar", "expected to find value")

	time.Sleep(200 * time.Millisecond)

	v, ok = m.Get("foo")
	assert.Equalf(t, ok, false, "expected to not find key")
	assert.Equalf(t, v, "", "expected v to be empty string")
}

func TestIntExpiringMap(t *testing.T) {
	m := New[string, int](100 * time.Millisecond)

	m.Set("foo", 1)

	v, ok := m.Get("foo")
	assert.Equalf(t, ok, true, "expected to find key")
	assert.Equalf(t, v, 1, "expected to find value")

	time.Sleep(200 * time.Millisecond)

	v, ok = m.Get("foo")
	assert.Equalf(t, ok, false, "expected to not find key")
	assert.Equalf(t, v, 0, "expected v to be empty string")
}

func TestStructExpiringMap(t *testing.T) {
	type Foo struct {
		Bar string
	}

	m := New[string, *Foo](100 * time.Millisecond)

	m.Set("foo", &Foo{Bar: "bar"})

	v, ok := m.Get("foo")
	assert.Equalf(t, ok, true, "expected to find key")
	assert.Equalf(t, v.Bar, "bar", "expected to find value")

	time.Sleep(200 * time.Millisecond)

	v, ok = m.Get("foo")
	assert.Equalf(t, ok, false, "expected to not find key")
	assert.Nilf(t, v, "expected v to be nil")
}

func TestLenExpiringMap(t *testing.T) {
	m := New[string, string](100 * time.Millisecond)

	m.Set("foo", "bar")

	l := m.Len()
	assert.Equalf(t, l, 1, "expected len to be 1")

	time.Sleep(200 * time.Millisecond)

	l = m.Len()
	assert.Equalf(t, l, 0, "expected len to be 0")
}

func TestItemsExpiringMap(t *testing.T) {
	m := New[string, string](100 * time.Millisecond)

	m.Set("foo", "bar")

	mm := m.Items()
	assert.Equalf(t, mm["foo"], "bar", "expected to find key")

	time.Sleep(200 * time.Millisecond)

	mm = m.Items()
	assert.Equalf(t, mm["foo"], "", "expected to find key")

}

func TestExpiringChannel(t *testing.T) {
	ch := make(chan []string, 1)

	m := New[string, string](100 * time.Millisecond).WithEvictionChannel(ch)

	m.Set("foo", "bar")

	mm := m.Items()
	assert.Equalf(t, mm["foo"], "bar", "expected to find key")

	time.Sleep(200 * time.Millisecond)

	mm = m.Items()
	assert.Equalf(t, mm["foo"], "", "expected to find key")

	items := <-ch
	require.Lenf(t, items, 1, "expected to find 1 expired item")
	assert.Equalf(t, items[0], "bar", "expected to find expired item")

}

func TestExpiringFunctionTrue(t *testing.T) {
	m := New[string, string](100 * time.Millisecond).WithEvictionFunction(func(string) bool {
		return true
	})

	m.Set("foo", "bar")

	mm := m.Items()
	assert.Equalf(t, mm["foo"], "bar", "expected to find key")

	time.Sleep(200 * time.Millisecond)

	mm = m.Items()
	assert.Equalf(t, mm["foo"], "", "expected to find key")

}

func TestExpiringFunctionFalse(t *testing.T) {
	m := New[string, string](100 * time.Millisecond).WithEvictionFunction(func(string) bool {
		return false
	})

	m.Set("foo", "bar")

	mm := m.Items()
	assert.Equalf(t, mm["foo"], "bar", "expected to find key")

	time.Sleep(200 * time.Millisecond)

	mm = m.Items()
	assert.Lenf(t, mm, 1, "expected to have 1 item")

}
