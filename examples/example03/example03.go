package main

import (
	"fmt"
	"log"
	"time"

	"github.com/carlescere/goback"
)

// Creates a function that will fail 6 times before connecting
func retryGenerator() func(chan bool) {
	var failedAttempts = 6
	return func(done chan bool) {
		if failedAttempts > 0 {
			failedAttempts--
			return
		}
		done <- true
	}
}

func connect(retry func(chan bool), b goback.Backoff) {
	done := make(chan bool, 1)
	for {
		retry(done)
		select {
		case err := <-goback.After(b):
			if err != nil {
				log.Fatalf("Error connecting: %v", err)
			}
			log.Printf("Problem connecting")
			continue
		case <-done:
			//conected
			b.Reset()
			return
		}
	}

}

func main() {
	retry := retryGenerator()
	b := &goback.SimpleBackoff{
		Min:         100 * time.Millisecond,
		Max:         60 * time.Second,
		Factor:      2,
		MaxAttempts: 4,
	}
	connect(retry, b)
	// Duplicates the time each time from a minimum of 100ms to a maximum of 1 min.
	fmt.Println("Connected")
}
