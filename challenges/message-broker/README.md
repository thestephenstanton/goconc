# Concurrent Message Broker with Topic Subscriptions

## Problem Statement

Build a thread-safe in-memory message broker similar to Kafka, NATS, or Redis Pub/Sub. Multiple publishers can send messages to topics, and multiple subscribers can consume messages from topics they're interested in. The broker must handle subscriptions, message delivery, and backpressure efficiently.

## Requirements

1. Implement a `MessageBroker` with these methods:
   - `Publish(topic string, message Message) error`: Publish a message to a topic
   - `Subscribe(topic string, subscriberID string, bufferSize int) (<-chan Message, error)`: Subscribe to a topic and get a channel for messages
   - `Unsubscribe(topic string, subscriberID string) error`: Unsubscribe from a topic
   - `Close()`: Shutdown the broker gracefully

2. Message delivery semantics:
   - Messages are delivered to all active subscribers of a topic
   - If a subscriber's buffer is full, the broker should either drop the message or block (you choose, but document your decision)
   - Messages should be delivered in the order they were published
   - Subscribers can consume at their own pace

3. The broker should handle:
   - Multiple publishers publishing to the same topic concurrently
   - Multiple subscribers on the same topic receiving messages independently
   - Subscribers joining/leaving while messages are being published
   - Topics being created dynamically (first publish creates the topic)

## Example Usage
```go
package main

import (
    "fmt"
    "time"
)

type Message struct {
    ID        string
    Payload   string
    Timestamp time.Time
}

func main() {
    broker := NewMessageBroker()

    // Subscriber 1 subscribes to "orders" topic
    ordersCh1, _ := broker.Subscribe("orders", "subscriber-1", 10)
    go func() {
        for msg := range ordersCh1 {
            fmt.Printf("Sub1 received: %s\n", msg.Payload)
            time.Sleep(100 * time.Millisecond) // Slow consumer
        }
    }()

    // Subscriber 2 subscribes to same topic
    ordersCh2, _ := broker.Subscribe("orders", "subscriber-2", 10)
    go func() {
        for msg := range ordersCh2 {
            fmt.Printf("Sub2 received: %s\n", msg.Payload)
            time.Sleep(10 * time.Millisecond) // Fast consumer
        }
    }()

    // Publisher 1 publishes to "orders"
    go func() {
        for i := 0; i < 20; i++ {
            broker.Publish("orders", Message{
                ID:        fmt.Sprintf("order-%d", i),
                Payload:   fmt.Sprintf("Order data %d", i),
                Timestamp: time.Now(),
            })
            time.Sleep(50 * time.Millisecond)
        }
    }()

    // Publisher 2 publishes to different topic
    paymentsCh, _ := broker.Subscribe("payments", "subscriber-3", 5)
    go func() {
        for msg := range paymentsCh {
            fmt.Printf("Sub3 received: %s\n", msg.Payload)
        }
    }()

    go func() {
        for i := 0; i < 10; i++ {
            broker.Publish("payments", Message{
                ID:      fmt.Sprintf("payment-%d", i),
                Payload: fmt.Sprintf("Payment data %d", i),
            })
            time.Sleep(100 * time.Millisecond)
        }
    }()

    // Let it run for a bit
    time.Sleep(3 * time.Second)

    // Unsubscribe one subscriber
    broker.Unsubscribe("orders", "subscriber-1")
    
    time.Sleep(1 * time.Second)
    broker.Close()
}
```

## Backpressure Scenario
```go
func backpressureExample() {
    broker := NewMessageBroker()

    // Slow subscriber with small buffer
    slowCh, _ := broker.Subscribe("data", "slow-consumer", 2)
    go func() {
        for msg := range slowCh {
            fmt.Printf("Slow consumer processing: %s\n", msg.Payload)
            time.Sleep(1 * time.Second) // Very slow
        }
    }()

    // Fast publisher
    go func() {
        for i := 0; i < 10; i++ {
            err := broker.Publish("data", Message{
                ID:      fmt.Sprintf("msg-%d", i),
                Payload: fmt.Sprintf("Message %d", i),
            })
            if err != nil {
                fmt.Printf("Failed to publish msg-%d: %v\n", i, err)
            }
            time.Sleep(100 * time.Millisecond)
        }
    }()

    time.Sleep(5 * time.Second)
    
    // What happens when slow consumer's buffer fills up?
    // Option 1: Drop messages and return error
    // Option 2: Block publisher until space available
    // Option 3: Drop oldest message from buffer
    // You decide and document your choice!
}
```

## Timeline Example
```
Time | Topic: "orders"                    | Action
-----|------------------------------------|-----------------------------------------
t0   | Subscribers: []                    | broker.Subscribe("orders", "sub-1", 5)
t1   | Subscribers: [sub-1]               | broker.Subscribe("orders", "sub-2", 5)
t2   | Subscribers: [sub-1, sub-2]        | broker.Publish("orders", msg1)
t3   | sub-1 buffer: [msg1]               | Both subscribers receive msg1
     | sub-2 buffer: [msg1]               |
t4   | Subscribers: [sub-1, sub-2]        | broker.Publish("orders", msg2)
t5   | sub-1 buffer: [msg1, msg2]         | Both receive msg2
     | sub-2 buffer: [msg1, msg2]         |
t6   | sub-1 buffer: [msg2]               | sub-1 consumes msg1
     | sub-2 buffer: []                   | sub-2 already consumed both (faster)
t7   | Subscribers: [sub-2]               | broker.Unsubscribe("orders", "sub-1")
t8   | Subscribers: [sub-2]               | broker.Publish("orders", msg3)
t9   | sub-2 buffer: [msg3]               | Only sub-2 receives msg3 (sub-1 gone)
```

## Evaluation Criteria

- **Correctness**: All subscribers receive all messages published to their topics after they subscribe
- **Thread-Safety**: Concurrent publishes, subscribes, and unsubscribes work correctly
- **Backpressure Handling**: Clear strategy for handling slow consumers
- **Resource Management**: Channels closed properly on unsubscribe, no goroutine leaks
- **Efficiency**: No unnecessary blocking, good use of Go channels and goroutines
- **Code Quality**: Clean design, well-organized code, documented decisions
