# goback
[![GoDoc](https://godoc.org/github.com/carlescere/goback?status.svg)](https://godoc.org/github.com/carlescere/goback)
[![Build Status](https://travis-ci.org/carlescere/goback.svg)](https://travis-ci.org/carlescere/goback)
[![Coverage Status](https://coveralls.io/repos/carlescere/goback/badge.svg)](https://coveralls.io/r/carlescere/goback)


Goback implements a simple exponential backoff.

An exponential backoff approach is typically used when treating with potentially faulty/slow systems. If a system fails quick retries may exacerbate the system specially when the system is dealing with several clients. In this case a backoff provides the faulty system enough room to recover.

## How to use
```go
func main() {
        b := &goback.SimpleBackoff(
                Min:    100 * time.Millisecond,
                Max:    60 * time.Second,
                Factor: 2,
        )
        goback.Wait(b)           // sleeps 100ms
        goback.Wait(b)           // sleeps 200ms
        goback.Wait(b)           // sleeps 400ms
        fmt.Println(b.NextRun()) // prints 800ms
        b.Reset()                // resets the backoff
        goback.Wait(b)           // sleeps 100ms
}
```

Furter examples can be found in the examples folder.

## Strategies
At the moment there are two backoff strategies implemented.

### Simple Backoff
It starts with a minumum duration and multiplies it by the factor until the maximum waiting time is reached. In that case it will return `Max`.

The optional `MaxAttempts` will limit the maximum number of retries and will return an error when is exceeded.

### Jitter Backoff
The Jitter strategy is based on the simple backoff but adds a light randomisation to minimise collisions between contending clients.

The result of the 'NextDuration()' method will be a random duration between `[d-min, d+min]` where `d` is the expected duration without jitter and `min` is the minimum duration.

### Extensibility
By creating structs that implement the methods of the `Backoff` interface you will be able to use them as a backoff strategy.

A naive example of this is:
```go
type NaiveBackoff struct{}

func (b *NaiveBackoff) NextAttempt() (time.Duration, error) { return time.Second, nil }
func (b *NaiveBackoff) Reset() { }
```
This will return always a 1s duration.

## Credits
This package is inspired in https://github.com/jpillora/backoff

## License
Distributed under MIT license. See `LICENSE` file for more information.
