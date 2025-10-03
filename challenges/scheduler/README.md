# Concurrent Task Scheduler

## Problem Statement

Design and implement a thread-safe task scheduler that can execute tasks at scheduled times. Multiple threads will be submitting tasks concurrently, and a separate worker thread (or threads) should execute tasks when their scheduled time arrives.

## Requirements

1. Implement a `TaskScheduler` class with the following methods:
   - `scheduleTask(task, delayInMillis)`: Schedule a task to run after the specified delay
   - `start()`: Start the scheduler (begins processing scheduled tasks)
   - `shutdown()`: Gracefully shutdown the scheduler

2. Tasks should execute as close to their scheduled time as possible:
   - Tasks scheduled for the same time can execute in any order
   - Tasks should not execute before their scheduled time
   - If multiple tasks are ready, they should be executed in the order they were scheduled

3. The scheduler must be thread-safe:
   - Multiple threads can call `scheduleTask()` concurrently
   - Tasks should execute exactly once
   - No race conditions or deadlocks

4. Your solution should be efficient:
   - Worker threads should not busy-wait
   - Use appropriate synchronization primitives (condition variables, mutexes, etc.)


## Example Usage

TaskScheduler scheduler = new TaskScheduler();
scheduler.start();
// Thread 1
scheduler.scheduleTask(() -> System.out.println("Task A"), 1000); // runs after 1 second
// Thread 2 (concurrent)
scheduler.scheduleTask(() -> System.out.println("Task B"), 500);  // runs after 0.5 seconds
// Thread 3 (concurrent)
scheduler.scheduleTask(() -> System.out.println("Task C"), 1000); // runs after 1 second
// Expected output (approximately):
// (after 500ms)  Task B
// (after 1000ms) Task A
// (after 1000ms) Task C
scheduler.shutdown();

## Evaluation Criteria

- **Correctness**: Tasks execute at the correct time and exactly once
- **Thread-Safety**: Proper synchronization for concurrent access
- **Efficiency**: No busy-waiting, appropriate use of condition variables
- **Shutdown Handling**: Graceful shutdown without resource leaks
- **Code Quality**: Clean, readable, well-commented code
