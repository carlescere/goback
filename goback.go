// Package goback implements a simple exponential backoff
//
// An exponential backoff approach is typically used when treating with potentially
// faulty/slow systems. If a system fails quick retries may exacerbate the system
// specially when the system is dealing with several clients. In this case a backoff
// provides the faulty system enough room to recover.
//
// Simple example:
//  func main() {
//	    b := &goback.SimpleBackoff(
//		    Min:    100 * time.Millisecond,
//		    Max:    60 * time.Second,
//		    Factor: 2,
//	    )
//	    goback.Wait(b)           // sleeps 100ms
//	    goback.Wait(b)           // sleeps 200ms
//	    goback.Wait(b)           // sleeps 400ms
//	    fmt.Println(b.NextRun()) // prints 800ms
//	    b.Reset()                // resets the backoff
//	    goback.Wait(b)           // sleeps 100ms
//  }
//
// Furter examples can be found in the examples folder in the repository.
package goback

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

var (
	// ErrMaxAttemptsExceeded indicates that the maximum retries has been
	// excedeed. Usually to consider a service unreachable/unavailable.
	ErrMaxAttemptsExceeded = errors.New("maximum of attempts exceeded")
)

// Backoff is the interface that any Backoff strategy needs to implement.
type Backoff interface {
	// NextAttempt returns the duration to wait for the next retry.
	NextAttempt() (time.Duration, error)
	// Reset clears the number of tries. Next call to NextAttempt will return
	// the minimum backoff time (if there is no error).
	Reset()
}

// SimpleBackoff provides a simple strategy to backoff.
type SimpleBackoff struct {
	Attempts    int
	MaxAttempts int
	Factor      float64
	Min         time.Duration
	Max         time.Duration
}

// NextAttempt returns the duration to wait for the next retry.
func (b *SimpleBackoff) NextAttempt() (time.Duration, error) {
	if b.MaxAttempts > 0 && b.Attempts >= b.MaxAttempts {
		return 0, ErrMaxAttemptsExceeded
	}
	next := GetNextDuration(b.Min, b.Max, b.Factor, b.Attempts)
	b.Attempts++
	return next, nil
}

// Reset clears the number of tries. Next call to NextAttempt will return
// the minimum backoff time (if there is no error).
func (b *SimpleBackoff) Reset() {
	b.Attempts = 0
}

// JitterBackoff provides an strategy similar to SimpleBackoff but lightly randomises
// the duration to minimise collisions between contending clients.
type JitterBackoff SimpleBackoff

// NextAttempt returns the duration to wait for the next retry.
func (b *JitterBackoff) NextAttempt() (time.Duration, error) {
	if b.MaxAttempts > 0 && b.Attempts >= b.MaxAttempts {
		return 0, ErrMaxAttemptsExceeded
	}
	next := GetNextDuration(b.Min, b.Max, b.Factor, b.Attempts)
	next = addJitter(next, b.Min)
	b.Attempts++
	return next, nil
}

// GetNextDuration returns the duration for the strategies considering the minimum and
// maximum durations, the factor of increase and the number of attemtps tried.
func GetNextDuration(min, max time.Duration, factor float64, attempts int) time.Duration {
	d := time.Duration(float64(min) * math.Pow(factor, float64(attempts)))
	if d > max {
		return max
	}
	return d
}

// Wait sleeps for the duration of the time specified by the backoff strategy.
func Wait(b Backoff) error {
	next, err := b.NextAttempt()
	if err != nil {
		return err
	}
	time.Sleep(next)
	return nil
}

// After returns a channel that will be called after the time specified by the backoff
// strategy or will exit immediately with an error.
func After(b Backoff) <-chan error {
	c := make(chan error, 1)
	next, err := b.NextAttempt()
	if err != nil {
		c <- err
		return c
	}
	go func() {
		time.Sleep(next)
		c <- nil
	}()
	return c
}

// addJitter randomises the final duration
func addJitter(next, min time.Duration) time.Duration {
	return time.Duration(rand.Float64()*float64(2*min) + float64(next-min))
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
