package main

import (
	"context"
	"log/slog"
	"os"
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
	b := NewBoundedBuffer[int](2)

	b.Put(1)
	b.Put(2)
	println(b.Take()) // 1

	b.Put(3)

	ok := b.TryPut(4, 1*time.Second)
	println(ok) // false

	go func() {
		time.AfterFunc(500*time.Millisecond, func() {
			b.Take()
		})
	}()

	ok = b.TryPut(4, 1*time.Second)
	println(ok) // true

	return nil
}

type BoundedBuffer[T any] struct {
	items chan T
}

func NewBoundedBuffer[T any](size int) BoundedBuffer[T] {
	return BoundedBuffer[T]{
		items: make(chan T, size),
	}
}

func (b *BoundedBuffer[T]) Put(item T) {
	b.items <- item
}

func (b *BoundedBuffer[T]) Take() T {
	return <-b.items
}

func (b *BoundedBuffer[T]) TryPut(item T, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return false
	case b.items <- item:
		return true
	}
}

func (b *BoundedBuffer[T]) TryTake(timeout time.Duration) *T {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil
	case item := <-b.items:
		return &item
	}
}
