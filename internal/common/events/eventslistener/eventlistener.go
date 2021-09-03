package eventslistener

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Listener struct {
	Uuid     string
	Response chan []byte
}

type EventListener struct {
	listeners map[string]*Listener
}

func (e *EventListener) AddListener(l *Listener) error {
	if e.listeners[l.Uuid] != nil {
		return fmt.Errorf("listener already exists, uuid: %v", l.Uuid)
	}
	e.listeners[l.Uuid] = l
	return nil
}

func (e *EventListener) RemoveListener(l *Listener) error {
	if e.listeners[l.Uuid] == nil {
		return fmt.Errorf("listener not found, uuid: %v", l.Uuid)
	}
	delete(e.listeners, l.Uuid)
	return nil
}

func (e *EventListener) processEvent(b *[]byte) error {
	data := events.EventMessage{}
	if err := json.Unmarshal(*b, &data); err != nil {
		return err
	}

	for _, l := range e.listeners {
		if data.Key == l.Uuid {
			l.Response <- []byte(data.Value)
		}
	}
	return nil
}

func NewEventListener(topic string, partition int, c *config.Config, log *zap.Logger) *EventListener {
	conn, _ := kafka.DialLeader(context.Background(), "tcp", c.KafkaUri, topic, partition)
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
			log.Error("failed to close batch in event listener")
		}
		if err := conn.Close(); err != nil {
			log.Error("failed to close connection in event listener")
		}
	}()
	return &el
}
