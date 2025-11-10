package redis

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLocker(t *testing.T) {
	rd := New("127.0.0.1:6379")
	l := NewLocker(rd)

	_, ok := l.Lock("aabb")
	assert.Equal(t, ok, true)

	_, ok = l.Lock("aabb")
	assert.Equal(t, ok, false)

	time.Sleep(time.Second * 30)

	r, ok := l.Lock("aabb")
	assert.Equal(t, ok, true)

	r.Unlock()
	_, ok = l.Lock("aabb")
	assert.Equal(t, ok, true)
}