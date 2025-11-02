package main

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// NOTE: so I read things that said that the order that Signal() calls the waiting go routines
// are not preserved but then I saw other spots where it said it was. So I wanted to write a test
// to see what actually happens. The important part here is the sig chan. With this, we can control
// the next Signal() gets called. If you uncomment both of those sig parts, you'll see that the
// test will actually fail. The reason for this is because Cond does in fact maintain order (it actually
// a double linked list under the hood) however with Signal() wakes up the go routines, you have
// no control over what order the are scehduled in. The sig helps synchronize these and proves
// that they activate in the same order that they originally called Wait() in
func TestSignal(t *testing.T) {
	s := make([]int, 0, 100)

	c := sync.NewCond(&sync.Mutex{})

	var wg sync.WaitGroup
	wg.Add(cap(s))

	sig := make(chan struct{})

	do := func(i int) {
		defer wg.Done()

		c.L.Lock()

		c.Wait()
		<-sig // signal the signal loop to go to next go routine

		s = append(s, i)

		c.L.Unlock()
	}

	for i := range cap(s) {
		go do(i)

		// NOTE: this is a little hacky because these go routines can technically be scheduled out of order
		// but this should do for not
		time.Sleep(10 * time.Millisecond)
	}

	// hacky way to give the go routines some time to run and get to the c.Wait()
	time.Sleep(1 * time.Second)

	for range cap(s) {
		c.Signal()

		// this will allow us to not signal our next signal until the go routine as actually woken up and started again
		sig <- struct{}{}
	}

	wg.Wait()

	expected := make([]int, cap(s))
	for i := range cap(s) {
		expected[i] = i
	}

	assert.Equal(t, expected, s)
}
