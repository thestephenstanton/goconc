package main

import (
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
	l := NewReadWriteLock()

	slog.Info("acquiring read 1")
	l.AcquireRead()
	slog.Info("acquired read 1")

	slog.Info("acquiring read 2")
	l.AcquireRead()
	slog.Info("acquired read 2")

	go func() {
		time.Sleep(2 * time.Second)
		l.ReleaseRead()
		l.ReleaseRead()
	}()

	slog.Info("acquiring lock")
	l.AcquireWrite() // should pause here for 2 seconds before acquiring
	slog.Info("acquired lock")

	go func() {
		time.Sleep(2 * time.Second)
		l.ReleaseWrite()
	}()

	slog.Info("acquiring read 3")
	l.AcquireRead() // should pause here for 2 seconds before acquiring
	slog.Info("acquired read 3")

	return nil
}

type ReadWriteLock struct {
	mu *sync.Mutex

	tickets             uint
	currentTicketNumber uint

	numReaders uint
	numWriters uint

	canRead  *sync.Cond
	canWrite *sync.Cond
}

func NewReadWriteLock() *ReadWriteLock {
	mu := &sync.Mutex{}

	return &ReadWriteLock{
		mu:       mu,
		canRead:  sync.NewCond(mu),
		canWrite: sync.NewCond(mu),
	}
}

func (l *ReadWriteLock) AcquireRead() {
	l.mu.Lock()
	defer l.mu.Unlock()

	for l.numWriters > 0 {
		l.canRead.Wait()
	}

	l.numReaders++
}

func (l *ReadWriteLock) ReleaseRead() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.numReaders--

	if l.numReaders == 0 && l.numWriters > 0 {
		// no more readers but there are writers waiting
		l.canWrite.Broadcast()
	}
}

func (l *ReadWriteLock) AcquireWrite() {
	l.mu.Lock()
	defer l.mu.Unlock()

	myTicket := l.tickets
	l.tickets++
	l.numWriters++

	for l.numReaders != 0 || myTicket != l.currentTicketNumber {
		l.canWrite.Wait()
	}
}

func (l *ReadWriteLock) ReleaseWrite() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.numWriters--
	l.currentTicketNumber++

	if l.numWriters > 0 {
		l.canWrite.Broadcast()
	} else {
		l.canRead.Broadcast()
	}
}
