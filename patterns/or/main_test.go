package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOr_NumGoRoutines(t *testing.T) {
	testCases := []struct {
		dones    int
		expected int32
	}{
		{dones: 1, expected: 0},
		{dones: 2, expected: 1},
		{dones: 3, expected: 1},
		{dones: 4, expected: 2},
		{dones: 5, expected: 2},
		{dones: 6, expected: 3},
		{dones: 7, expected: 3},
		{dones: 8, expected: 4},
		{dones: 9, expected: 4},
		{dones: 10, expected: 5},
		{dones: 20, expected: 10},
		{dones: 50, expected: 25},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("dones (%d)", tc.dones), func(t *testing.T) {
			i.Store(0)

			cs := make([]<-chan any, tc.dones)

			for i := 0; i < tc.dones; i++ {
				cs[i] = later(time.Duration(i+1) * time.Millisecond)
			}

			<-or(cs...)

			assert.Equal(t, tc.expected, i.Load())
		})
	}
}
