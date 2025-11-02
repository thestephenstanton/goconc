package main

import (
	"log/slog"
	"os"
	"sync"
	"time"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	err := run()
	if err != nil {
		slog.Error("starting application", "error", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	c := make(chan any)

	go func() {
		defer close(c)

		for i := range 5 {
			c <- i
		}
	}()

	c1, c2 := tee(nil, c)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for i := range c1 {
			slog.Info("from c1", "i", i)
		}
	}()

	go func() {
		defer wg.Done()

		for i := range c2 {
			time.Sleep(1 * time.Second)
			slog.Info("from c2", "i", i)
		}
	}()

	wg.Wait()

	return nil
}

func tee(done chan struct{}, c <-chan any) (<-chan any, <-chan any) {
	c1, c2 := make(chan any), make(chan any)

	go func() {
		defer close(c1)
		defer close(c2)

	loop:
		for v := range c {
			sc1, sc2 := c1, c2

			for range 2 {
				select {
				case <-done:
					break loop
				case sc1 <- v:
					sc1 = nil
				case sc2 <- v:
					sc2 = nil
				}
			}
		}
	}()

	return c1, c2
}
