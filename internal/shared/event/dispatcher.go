package event

import (
	"context"
	"log/slog"
	"sync"

	"github.com/S4F4Y4T/goWebService/internal/shared/domain"
)

type HandlerFunc func(ctx context.Context, event domain.DomainEvent) error

// Dispatcher is an in-memory synchronous event dispatcher
type Dispatcher struct {
	mu       sync.RWMutex
	handlers map[string][]HandlerFunc
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string][]HandlerFunc),
	}
}

func (d *Dispatcher) Subscribe(topic string, handler HandlerFunc) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[topic] = append(d.handlers[topic], handler)
}

func (d *Dispatcher) Dispatch(ctx context.Context, events []domain.DomainEvent) {
	for _, event := range events {
		d.mu.RLock()
		handlers, ok := d.handlers[event.Topic()]
		d.mu.RUnlock()

		if !ok {
			continue
		}

		for _, handler := range handlers {
			if err := handler(ctx, event); err != nil {
				slog.Error("Failed to handle domain event", 
					"topic", event.Topic(), 
					"error", err)
			}
		}
	}
}
