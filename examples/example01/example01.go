package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/carlescere/goback"
)

// Creates a function that will fail 6 times before connecting
func faultyTCPGenerator() func(string, string, string) (net.Listener, error) {
	var failedAttempts = 6
	return func(protocol, address, port string) (net.Listener, error) {
		if failedAttempts > 0 {
			failedAttempts--
			return nil, fmt.Errorf("Haha!") // :)
		}
		l, err := net.Listen(protocol, fmt.Sprintf("%s:%s", address, port))
		return l, err
	}
}

func main() {
	tcpConnect := faultyTCPGenerator()
	// Duplicates the time each time from a minimum of 100ms to a maximum of 1 min.
	b := &goback.SimpleBackoff{
		Min:    100 * time.Millisecond,
		Max:    60 * time.Second,
		Factor: 2,
	}
	for {
		l, err := tcpConnect("tcp", "localhost", "5000")
		if err != nil { // fail to connect
			log.Printf("Error connecting: %v", err)
			goback.Wait(b) // Exponential backoff
			continue
		}
		defer l.Close()
		// connected
		log.Printf("Connected!")
		b.Reset() // Reset number of attempts. Brings backoff time to the minimum
		break     // Here be dragons
	}
}
