package goback

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAfter(t *testing.T) {
	b := &SimpleBackoff{
		Factor:      2,
		Min:         100 * time.Millisecond,
		Max:         2 * time.Second,
		MaxAttempts: 1,
	}
	var t1, t2, t3 time.Time
	var err1, err2 error
	t1 = time.Now()
	select {
	case err1 = <-After(b):
		t2 = time.Now()
	case <-time.After(101 * time.Millisecond):
		t.Error("Not executed on time")
	}
	select {
	case err2 = <-After(b):
	case <-time.After(time.Millisecond):
		t.Error("Not executed on time")
	}
	t3 = time.Now()
	assert.WithinDuration(t, t1.Add(100*time.Millisecond), t2, time.Millisecond)
	assert.WithinDuration(t, t2, t3, time.Millisecond)
	assert.Nil(t, err1)
	assert.NotNil(t, err2)

}

func TestWait(t *testing.T) {
	b := &SimpleBackoff{
		Factor:      2,
		Min:         100 * time.Millisecond,
		Max:         2 * time.Second,
		MaxAttempts: 1,
	}
	t1 := time.Now()
	err1 := Wait(b)
	t2 := time.Now()
	err2 := Wait(b)
	t3 := time.Now()
	assert.WithinDuration(t, t1.Add(100*time.Millisecond), t2, time.Millisecond)
	assert.WithinDuration(t, t2, t3, time.Millisecond)
	assert.Nil(t, err1)
	assert.NotNil(t, err2)

}

func TestSimple(t *testing.T) {
	b := &SimpleBackoff{
		Factor:      2,
		Min:         100 * time.Millisecond,
		Max:         2 * time.Second,
		MaxAttempts: 6,
	}
	next, err := b.NextAttempt()
	assert.Equal(t, 100*time.Millisecond, next)
	assert.Nil(t, err)
	next, err = b.NextAttempt()
	assert.Equal(t, 200*time.Millisecond, next)
	assert.Nil(t, err)
	next, err = b.NextAttempt()
	assert.Equal(t, 400*time.Millisecond, next)
	assert.Nil(t, err)
	next, err = b.NextAttempt()
	assert.Equal(t, 800*time.Millisecond, next)
	assert.Nil(t, err)
	next, err = b.NextAttempt()
	assert.Equal(t, 1600*time.Millisecond, next)
	assert.Nil(t, err)
	next, err = b.NextAttempt()
	assert.Equal(t, 2*time.Second, next)
	assert.Nil(t, err)
	next, err = b.NextAttempt()
	assert.NotNil(t, err)
	b.Reset()
	next, err = b.NextAttempt()
	assert.Equal(t, 100*time.Millisecond, next)
	assert.Nil(t, err)
}

func TestJitter(t *testing.T) {
	min := 100 * time.Millisecond
	b := &JitterBackoff{
		Factor:      2,
		Min:         min,
		Max:         10 * time.Second,
		MaxAttempts: 3,
	}
	next, err := b.NextAttempt()
	between(t, next, 100*time.Millisecond, min)
	assert.Nil(t, err)
	next, err = b.NextAttempt()
	between(t, next, 200*time.Millisecond, min)
	assert.Nil(t, err)
	next, err = b.NextAttempt()
	between(t, next, 400*time.Millisecond, min)
	assert.Nil(t, err)
	next, err = b.NextAttempt()
	assert.NotNil(t, err)
}

func between(t *testing.T, next, expected, offset time.Duration) {
	assert.True(t, next >= expected-offset, fmt.Sprintf("%v %v %v", next, expected, offset))
	assert.True(t, next <= expected+offset, fmt.Sprintf("%v %v %v", next, expected, offset))

}
