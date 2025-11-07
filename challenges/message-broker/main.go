package main

import (
	"log/slog"
	"os"
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
	return nil
}

type Broker struct {
	topics      map[string]*Topic
	subscribers map[string]*Subscriber
}

func (b *Broker) Publish(topic string, payload string) {
	t, found := b.topics[topic]
	if !found {
		t = NewTopic(topic)
		b.topics[topic] = t
	}

	t.Add(payload)
}

// func (b Broker) Subscribe(topic string, subscriberID string, bufferSize int) (<-chan Message, error) {
// }

type Subscriber struct {
	topic *Topic
	cur   *node
	c     chan Message
}

func (s *Subscriber) Start(done <-chan any) {
	go func() {
		for {
			if s.cur != nil {
				s.c <- s.cur.message
				s.cur = s.cur.next
				continue
			}

			s.c <- s.cur.message

			s.cur = s.cur.next
		}
	}()
}

type node struct {
	message Message
	next    *node
}

type Topic struct {
	name string
	tail *node
	head *node
}

func NewTopic(name string) *Topic {
	tail := &node{}
	head := tail

	return &Topic{
		name: name,
		tail: tail,
		head: head,
	}
}

func (t *Topic) Add(payload string) {
	message := Message{
		ID:        t.head.message.ID + 1, // TODO: probably should make this globally unique
		Payload:   payload,
		Timestamp: time.Now(),
	}

	next := &node{
		message: message,
	}

	t.head.next = next
	t.head = next
}

type Message struct {
	ID        int
	Payload   string
	Timestamp time.Time
}
