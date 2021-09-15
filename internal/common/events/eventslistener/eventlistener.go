package eventslistener

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Listener struct {
	Uuid     string
	Response chan []byte
}

type EventListener struct {
	listeners map[string]*Listener
	logger    *zap.Logger
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

func (e *EventListener) processEvent(m *kafka.Message) {
	mp := map[string]string{}
	if err := json.Unmarshal(m.Value, &mp); err != nil {
		e.logger.Error(err.Error())
		return
	}
	eventId := mp["EventId"]
	for _, l := range e.listeners {
		if eventId == l.Uuid {
			l.Response <- m.Value
		}
	}
}

func NewEventListener(
	ctx context.Context,
	c *config.Config,
	logger *zap.Logger,
) *EventListener {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{c.KafkaUri},
		Topic:    c.ResultsKafkaTopic,
		GroupID:  c.KafkaResultsConsumerGroupId,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	el := EventListener{make(map[string]*Listener), logger}

	go func() {
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				el.logger.Error(err.Error())
				continue
			}
			go el.processEvent(&m)
			select {
			case <-ctx.Done():
				break
			default:
				continue
			}
		}
	}()
	return &el
}
