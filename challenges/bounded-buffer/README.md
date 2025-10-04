# Concurrent Bounded Buffer

## Problem Statement

Implement a thread-safe bounded buffer (also known as a circular buffer or ring buffer) that supports multiple concurrent producers and consumers. This is a classic producer-consumer problem where producers add items to the buffer and consumers remove items from the buffer.

## Requirements

1. Implement a `BoundedBuffer<T>` class with the following methods:
   - `put(item)`: Add an item to the buffer. If the buffer is full, block until space becomes available.
   - `take()`: Remove and return an item from the buffer. If the buffer is empty, block until an item becomes available.
   - Optional: `tryPut(item, timeoutMs)`: Try to add an item with a timeout, return true if successful, false if timeout expires
   - Optional: `tryTake(timeoutMs)`: Try to remove an item with a timeout, return the item if successful, null if timeout expires

2. The buffer has a fixed capacity specified at construction time

3. The implementation must be thread-safe:
   - Multiple producer threads can call `put()` concurrently
   - Multiple consumer threads can call `take()` concurrently
   - Producers and consumers can operate concurrently with each other

4. Use proper synchronization:
   - Producers should block when the buffer is full
   - Consumers should block when the buffer is empty
   - Use condition variables (or equivalent) to avoid busy-waiting
   - No race conditions or deadlocks

## Example Usage

BoundedBuffer<Integer> buffer = new BoundedBuffer<>(5); // capacity of 5
// Producer thread 1
buffer.put(1);
buffer.put(2);
buffer.put(3);

// Producer thread 2 (concurrent)
buffer.put(4);
buffer.put(5);
buffer.put(6); // blocks because buffer is full

// Consumer thread 1 (concurrent)
int val1 = buffer.take(); // returns 1, now buffer.put(6) can proceed
int val2 = buffer.take(); // returns 2

// Consumer thread 2 (concurrent)
int val3 = buffer.take(); // returns 3

## Evaluation Criteria

- **Correctness**: Items are added and removed in FIFO order without loss or duplication
- **Thread-Safety**: Proper use of synchronization primitives (mutexes, condition variables, semaphores, etc.)
- **Blocking Behavior**: Proper blocking when full/empty, no busy-waiting
- **Efficiency**: Minimal lock contention, efficient notification of waiting threads
- **Code Quality**: Clean, readable, well-commented code
