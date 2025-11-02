package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	start := time.Now()

	<-or(
		later(2*time.Second),
		later(5*time.Second),
		later(1*time.Second),
		later(3*time.Second),
		later(4*time.Second),
	)

	fmt.Printf("waited %s\n", time.Since(start).String())
	fmt.Printf("%d go routines spawned\n", i.Load())
}

var i atomic.Int32

func later(t time.Duration) <-chan any {
	c := make(chan any)

	go func() {
		time.Sleep(t)
		close(c)
	}()

	return c
}

func or(c ...<-chan any) <-chan any {
	switch len(c) {
	case 0:
		return nil
	case 1:
		return c[0]
	}

	done := make(chan any)

	go func() {
		i.Add(1) // just here for test
		defer close(done)

		switch len(c) {
		case 2:
			// doing this allows us to have len(c)/2 go routines spawned
			// and if we didn't then it would be len(c)-1 go routines
			select {
			case <-c[0]:
			case <-c[1]:
			}
		default:
			select {
			case <-c[0]:
			case <-c[1]:
			case <-c[2]:
			case <-or(append(c[3:], done)...):
			}
		}
	}()

	return done
}
