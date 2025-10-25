# Concurrent Read-Write Lock with Writer Priority

## Problem Statement

Implement a thread-safe read-write lock that allows multiple concurrent readers OR a single writer, with writer priority. This means that when a writer is waiting, no new readers should be allowed to acquire the lock, ensuring writers don't starve.

## Requirements

1. Implement a `ReadWriteLock` class with the following methods:
   - `acquireRead()`: Acquire the lock for reading. Blocks if a writer holds the lock or writers are waiting.
   - `releaseRead()`: Release the read lock.
   - `acquireWrite()`: Acquire the lock for writing. Blocks if any readers or another writer holds the lock.
   - `releaseWrite()`: Release the write lock.

2. Concurrency rules:
   - Multiple readers can hold the lock simultaneously (if no writer is active or waiting)
   - Only one writer can hold the lock at a time
   - Readers and writers cannot hold the lock simultaneously
   - **Writer Priority**: When a writer is waiting, no new readers should acquire the lock (prevents writer starvation)

3. Your implementation must:
   - Be thread-safe for any number of concurrent readers and writers
   - Use low-level synchronization primitives (mutexes, condition variables, atomics)
   - Avoid deadlocks and race conditions
   - Not use built-in read-write lock implementations from standard libraries

4. Handle edge cases:
   - Multiple writers waiting (should be served in FIFO order)
   - Readers currently holding the lock when a writer arrives
   - Writer finishing while both readers and writers are waiting

## Example Usage

```
ReadWriteLock rwLock = new ReadWriteLock();
int sharedData = 0;

// Reader thread 1
rwLock.acquireRead();
int value1 = sharedData; // safe to read
rwLock.releaseRead();

// Reader thread 2 (concurrent with reader 1)
rwLock.acquireRead();
int value2 = sharedData; // safe to read
rwLock.releaseRead();

// Writer thread (while readers are active)
rwLock.acquireWrite(); // blocks until all readers finish
sharedData = 42;        // exclusive access
rwLock.releaseWrite();

// Reader thread 3 (arrives while writer is waiting)
rwLock.acquireRead();   // blocks even though readers are active
                        // (writer priority)
```

## Example Scenario Demonstrating Writer Priority

```
Time | Thread | Action                  | State
-----|--------|-------------------------|---------------------------
t0   | R1     | acquireRead()           | R1 reading
t1   | R2     | acquireRead()           | R1, R2 reading
t2   | W1     | acquireWrite() [blocks] | R1, R2 reading; W1 waiting
t3   | R3     | acquireRead() [blocks]  | R1, R2 reading; W1 waiting; R3 blocked (writer priority!)
t4   | R1     | releaseRead()           | R2 reading; W1 waiting; R3 blocked
t5   | R2     | releaseRead()           | W1 acquires lock (R3 still blocked)
t6   | W1     | releaseWrite()          | R3 can now proceed
```

## Evaluation Criteria

- **Correctness**: Properly enforces read-write semantics and writer priority
- **Thread-Safety**: Correct use of synchronization primitives without deadlocks or race conditions
- **Writer Priority**: New readers are blocked when writers are waiting
- **Fairness**: Writers waiting are served in FIFO order
- **Efficiency**: Minimal blocking and lock contention
- **Code Quality**: Clean, readable, well-commented code with clear state management
```
