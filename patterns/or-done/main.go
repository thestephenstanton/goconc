package main

import (
	"fmt"
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

	fmt.Printf("waited %s", time.Since(start).String())
}

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
		defer close(done)

		switch len(c) {
		case 2:
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
