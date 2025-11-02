package main

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
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
	cs := make([]<-chan any, 5)

	for i := range len(cs) {
		c := make(chan any)

		go func() {
			defer close(c)

			for j := range 5 {
				c <- fmt.Sprintf("%d from c%d", j, i)
			}
		}()

		cs[i] = c
	}

	in := fanIn(nil, cs...)

	for x := range in {
		slog.Info(x.(string))
	}

	return nil
}

func fanIn(done chan struct{}, cs ...<-chan any) <-chan any {
	stream := make(chan any)

	var wg sync.WaitGroup

	in := func(c <-chan any) {
		defer wg.Done()

		for x := range c {
			select {
			case <-done:
				return
			case stream <- x:
			}
		}
	}

	for _, c := range cs {
		wg.Add(1)

		go in(c)
	}

	go func() {
		wg.Wait()
		close(stream)
	}()

	return stream
}
