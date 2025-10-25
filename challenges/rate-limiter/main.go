package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"log/slog"
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

type rateLimiter interface {
	Allow(ctx context.Context) bool
}

func run() error {
	if len(os.Args) <= 1 {
		return errors.New("requires arg")
	}

	var r rateLimiter
	switch os.Args[1] {
	case "token-bucket":
		r = New(5, 1).Start()
	default:
		return fmt.Errorf("unknown rate limiter '%s'", os.Args[1])
	}

	limit := time.After(10 * time.Second)

loop:
	for {
		select {
		case <-limit:
			break loop
		case <-time.After(500 * time.Millisecond):
			ok := r.Allow(context.Background())
			slog.Info("hit", "ok", ok)
		}
	}

	return nil
}
