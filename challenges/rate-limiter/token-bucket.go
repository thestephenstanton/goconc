package main

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type tokenBucket struct {
	done chan any

	tokens      int
	maxTokens   int
	refreshRate int
	mu          sync.Mutex
}

func New(maxTokens int, refreshRatePerSecond int) *tokenBucket {
	return &tokenBucket{
		done:        make(chan any),
		tokens:      maxTokens,
		maxTokens:   maxTokens,
		refreshRate: refreshRatePerSecond,
		mu:          sync.Mutex{},
	}
}

func (t *tokenBucket) Start() *tokenBucket {
	go func() {
		for {
			select {
			case <-t.done:
				return
			case <-time.After(1 * time.Second):
				t.mu.Lock()

				if t.tokens < t.maxTokens {
					t.tokens += 1
					slog.Info("reload")
				}

				t.mu.Unlock()
			}
		}
	}()

	return t
}

func (t *tokenBucket) Stop() {
	close(t.done)
}

func (t *tokenBucket) Allow(ctx context.Context) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.tokens == 0 {
		return false
	}

	t.tokens -= 1

	return true
}
