package main

import (
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSlidingWindowCounterAllow(t *testing.T) {
	testCases := []struct {
		desc        string
		limiter     *SlidingWindowCounter
		timeElapsed time.Duration
		expected    bool
	}{
		{
			desc: "allows first request",
			limiter: &SlidingWindowCounter{
				limit:      2,
				windowSize: time.Minute,
				prevCount:  0,
				curCount:   0,
				mu:         &sync.Mutex{},
			},
			timeElapsed: 1 * time.Second,
			expected:    true,
		},
		{
			desc: "reject when at limit in current window",
			limiter: &SlidingWindowCounter{
				limit:      2,
				windowSize: time.Minute,
				prevCount:  0,
				curCount:   2,
				mu:         &sync.Mutex{},
			},
			timeElapsed: 1 * time.Second,
			expected:    false,
		},
		{
			desc: "weight of previous window doesn't put us over",
			limiter: &SlidingWindowCounter{
				limit:      10,
				windowSize: time.Minute,
				prevCount:  5, // 25% = 1.25
				curCount:   7,
				mu:         &sync.Mutex{},
			},
			timeElapsed: 1*time.Minute + 45*time.Second, // 75% of current window
			expected:    true,
		},
		{
			desc: "weight of previous window does put us over",
			limiter: &SlidingWindowCounter{
				limit:      10,
				windowSize: time.Minute,
				prevCount:  8, // 25% = 2
				curCount:   8,
				mu:         &sync.Mutex{},
			},
			timeElapsed: 45 * time.Second, // 75% of current window
			expected:    false,
		},
		{
			desc: "missed a window",
			limiter: &SlidingWindowCounter{
				limit:      10,
				windowSize: time.Minute,
				prevCount:  10, // minute 1 count
				curCount:   10, // minute 2 count
				mu:         &sync.Mutex{},
			},
			// at minute 4 (2+2) we will now have fully reset
			timeElapsed: 2 * time.Minute, // the prev count was actually from 5 minutes ago
			expected:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				tc.limiter.lastWindow = time.Now().Truncate(tc.limiter.windowSize)

				time.Sleep(tc.timeElapsed)

				actual := tc.limiter.Allow(t.Context())

				assert.Equal(t, tc.expected, actual)
			})
		})
	}
}
