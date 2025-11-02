package main

import (
	"fmt"
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
	c := sync.NewCond(&sync.Mutex{})

	for i := range 10 {
		waiter(i, c)
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("---------")
	time.Sleep(1 * time.Second)

	for range 10 {
		c.Signal()
	}

	time.Sleep(1 * time.Second)

	fmt.Println(a)

	return nil
}

var a = make([]int, 10)

func waiter(x int, c *sync.Cond) {
	go func() {
		fmt.Printf("%d is waiting\n", x)

		c.L.Lock()
		defer c.L.Unlock()
		c.Wait()

		a[x] = x

		fmt.Printf("%d has awakened\n", x)
	}()
}
