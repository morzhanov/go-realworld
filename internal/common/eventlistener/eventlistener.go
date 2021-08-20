package eventlistener

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// TODO: move to common package
type EventMessage struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Listener struct {
	Uuid     string
	Response chan []byte
}

// TODO: rename if name not fits
type EventListener struct {
	listeners map[string]*Listener
}

func check(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func (e *EventListener) AddListener(l *Listener) error {
	if e.listeners[l.Uuid] != nil {
		return fmt.Errorf("Listener already exists, uuid: %v", l.Uuid)
	}
	e.listeners[l.Uuid] = l
	return nil
}

func (e *EventListener) RemoveListener(l *Listener) error {
	if e.listeners[l.Uuid] == nil {
		return fmt.Errorf("Listener not found, uuid: %v", l.Uuid)
	}
	delete(e.listeners, l.Uuid)
	return nil
}

func (e *EventListener) processEvent(b *[]byte) {
	data := EventMessage{}
	err := json.Unmarshal(*b, &data)
	check(err)
	for _, l := range e.listeners {
		if data.Key == l.Uuid {
			l.Response <- []byte(data.Value)
		}
	}
}

func NewEventListener(topic string, partition int) *EventListener {
	// TODO: provide kafka uri
	uri := "192.168.0.180:32181"
	conn, _ := kafka.DialLeader(context.Background(), "tcp", uri, topic, partition)
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	el := EventListener{make(map[string]*Listener)}

	go func() {
		// TODO: maybe this initialization should be in a loop
		b := make([]byte, 10e3) // 10KB max per message
		for {
			_, err := batch.Read(b)
			if err != nil {
				break
			}
			el.processEvent(&b)
		}

		if err := batch.Close(); err != nil {
			log.Fatal("failed to close batch:", err)
		}
		if err := conn.Close(); err != nil {
			log.Fatal("failed to close connection:", err)
		}
	}()
	return &el
}
