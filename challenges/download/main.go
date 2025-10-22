package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
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
	concurrency := 10

	urls := make([]string, 0, 100)
	for i := range cap(urls) {
		urls = append(urls, fmt.Sprintf("url_%d", i))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	start := time.Now()
	m, err := downloadAll(ctx, urls, concurrency)
	if err != nil {
		return fmt.Errorf("downloading all: %w", err)
	}

	fmt.Printf("finished in %s\n", time.Since(start).String())

	fmt.Println("sum", m["sum"])

	return nil
}

func downloadAll(ctx context.Context, urls []string, maxWorkers int) (map[string]int, error) {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(maxWorkers)

	m := make(map[string]int, len(urls))

	var l sync.Mutex

	for _, url := range urls {
		g.Go(func() error {
			result, err := download(ctx, url)
			if err != nil {
				return fmt.Errorf("downloading file: %w", err)
			}

			l.Lock()
			m[url] = result
			m["sum"] += result
			l.Unlock()

			return nil
		})
	}

	err := g.Wait()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func download(ctx context.Context, url string) (int, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		// simulate work
		time.Sleep(100 * time.Millisecond)
		return rand.Intn(10) + 1, nil
	}
}
