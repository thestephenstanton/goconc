package main

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
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

var wg sync.WaitGroup

func run() error {
	var t TaskScheduler

	t.Start()

	t.Schedule(do("A"), 1000*time.Millisecond)
	t.Schedule(do("B"), 500*time.Millisecond)
	t.Schedule(do("C"), 1000*time.Millisecond)

	wg.Wait()

	t.Shutdown()

	t.Schedule(do("D"), 1*time.Second)

	time.Sleep(2 * time.Second)

	return nil
}

func do(key string) func() {
	wg.Add(1)

	return func() {
		slog.Info(fmt.Sprintf("task %s", key))

		wg.Done()
	}
}

type scheduledTask struct {
	task           func()
	completionTime time.Time
}

func newScheduledTask(task func(), delay time.Duration) scheduledTask {
	return scheduledTask{
		task:           task,
		completionTime: time.Now().Add(delay),
	}
}

type TaskScheduler struct {
	currentTimer         *time.Timer
	currentScheduledTask *scheduledTask
	newScheduledTask     chan scheduledTask

	scheduledTasks []scheduledTask

	done chan struct{}
	mu   *sync.Mutex
	once sync.Once
}

func (t *TaskScheduler) Schedule(task func(), delay time.Duration) {
	select {
	case <-t.done:
		return
	default:
	}

	newScheduledTask := newScheduledTask(task, delay)

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.currentScheduledTask != nil && t.currentScheduledTask.completionTime.Before(newScheduledTask.completionTime) {
		t.scheduledTasks = append(t.scheduledTasks, newScheduledTask)
		return
	}

	// since we are choosing not to run tasks in a go routine, this will block if a
	// task is long running while someone schedules another task, we could buffer
	// to help alleviate this but it kinda just kicks the can down the road
	t.newScheduledTask <- newScheduledTask
}

func (t *TaskScheduler) Start() {
	t.once.Do(func() {
		noStopTimer := time.NewTimer(0)
		noStopTimer.Stop()

		t.currentTimer = noStopTimer
		t.currentScheduledTask = nil
		t.newScheduledTask = make(chan scheduledTask)
		t.done = make(chan struct{})
		t.mu = &sync.Mutex{}

		go t.loop()
	})
}

func (t *TaskScheduler) loop() {
	for {
		select {
		case <-t.done:
			t.currentTimer.Stop()
			return
		case <-t.currentTimer.C:
			t.currentScheduledTask.task() // could run this in a go routine to not block

			t.mu.Lock()

			if len(t.scheduledTasks) == 0 {
				// no tasks so this is how we basically make it so nothing ever happens
				t.currentTimer = time.NewTimer(0)
				t.currentTimer.Stop()

				t.currentScheduledTask = nil

				t.mu.Unlock()

				continue
			}

			// get task that will complete the soonest
			k := 0
			for i := range t.scheduledTasks {
				if t.scheduledTasks[i].completionTime.Before(t.scheduledTasks[k].completionTime) {
					k = i
				}
			}

			newScheduledTask := t.scheduledTasks[k]

			t.currentScheduledTask = &newScheduledTask
			t.currentTimer = time.NewTimer(time.Until(newScheduledTask.completionTime))

			t.scheduledTasks = slices.Delete(t.scheduledTasks, k, k+1)

			t.mu.Unlock()

		case newScheduledTask := <-t.newScheduledTask:
			// this means we got a task scheduled that was going to happen sooner than the current task being run

			// if there is a currently scheduled task, add the current task back to the list
			if t.currentScheduledTask != nil {
				t.scheduledTasks = append(t.scheduledTasks, *t.currentScheduledTask)
			}

			t.currentScheduledTask = &newScheduledTask

			// change the timer to be this new one
			t.currentTimer.Stop()
			t.currentTimer = time.NewTimer(time.Until(newScheduledTask.completionTime))
		}
	}
}

func (t *TaskScheduler) Shutdown() {
	close(t.done)
}
