package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type SlidingWindowCounter struct {
	limit      float64
	windowSize time.Duration
	lastWindow time.Time
	prevCount  float64
	curCount   float64
	mu         *sync.Mutex
}

func NewSlidingWindowCounter(limit int, windowSize time.Duration) *SlidingWindowCounter {
	return &SlidingWindowCounter{
		limit:      float64(limit),
		windowSize: windowSize,
		lastWindow: time.Now().Truncate(windowSize),
		prevCount:  0,
		curCount:   0,
		mu:         &sync.Mutex{},
	}
}

func (s *SlidingWindowCounter) Allow(ctx context.Context) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	curWindow := now.Truncate(s.windowSize)

	if curWindow != s.lastWindow {
		if curWindow.Sub(s.lastWindow) >= 2*s.windowSize {
			// over 2 window sizes, reset everything
			s.prevCount = 0
			s.curCount = 0
		} else {
			// we have just moved 1 window
			s.prevCount = s.curCount
			s.curCount = 0
		}

		s.lastWindow = curWindow
	}

	elapsedTimeInWindow := now.Sub(curWindow)
	weight := 1 - (float64(elapsedTimeInWindow) / float64(s.windowSize))
	adjustedCount := (weight * s.prevCount) + s.curCount

	fmt.Println(adjustedCount, s.prevCount, s.curCount, weight)

	if adjustedCount >= s.limit {
		return false
	}

	s.curCount++

	return true
}
