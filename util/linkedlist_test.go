package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedList(t *testing.T) {
	l := LinkedList{}
	_, ok := l.Remove()
	assert.Equal(t, 0, l.Size())
	assert.Equal(t, false, ok)

	l.Add([]byte{'a'})
	assert.Equal(t, 1, l.Size())
	r, ok := l.Remove()
	assert.Equal(t, 0, l.Size())
	assert.Equal(t, []byte{'a'}, r)
	assert.Equal(t, true, ok)

	_, ok = l.Remove()
	assert.Equal(t, l.Size(), 0)
	assert.Equal(t, false, ok)
	_, ok = l.Remove()
	assert.Equal(t, l.Size(), 0)
	assert.Equal(t, false, ok)

	l.Add([]byte{'a'})
	l.Add([]byte{'b'})
	assert.Equal(t, l.Size(), 2)
	r, ok = l.Remove()
	assert.Equal(t, l.Size(), 1)
	assert.Equal(t, []byte{'a'}, r)
	assert.Equal(t, true, ok)
	r, ok = l.Remove()
	assert.Equal(t, l.Size(), 0)
	assert.Equal(t, []byte{'b'}, r)
	assert.Equal(t, true, ok)
	_, ok = l.Remove()
	assert.Equal(t, false, ok)

	for c := 'a'; c <= 'c'; c++ {
		l.Add([]byte{byte(c)})
	}
	for c := 'a'; c <= 'c'; c++ {
		r, ok := l.Remove()
		assert.Equal(t, []byte{byte(c)}, r)
		assert.Equal(t, true, ok)
	}
	_, ok = l.Remove()
	assert.Equal(t, false, ok)
}
