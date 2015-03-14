package main

import (
	"fmt"
	"time"

	"github.com/carlescere/goback"
)

func main() {
	// Jitter backoff randomises the retry duration to minimise condending clients
	b := &goback.JitterBackoff{
		Min:    100 * time.Millisecond,
		Max:    60 * time.Second,
		Factor: 2,
	}
	fmt.Println(b.NextAttempt())
	fmt.Println(b.NextAttempt())
	fmt.Println(b.NextAttempt())
	fmt.Println(b.NextAttempt())
	fmt.Println(b.NextAttempt())

}
