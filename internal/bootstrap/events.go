package bootstrap

import (
	"context"
	"log/slog"

	"github.com/S4F4Y4T/goWebService/internal/product"
	"github.com/S4F4Y4T/goWebService/internal/shared/domain"
	"github.com/S4F4Y4T/goWebService/internal/shared/event"
	"github.com/S4F4Y4T/goWebService/internal/user"
)

// RegisterEventHandlers wires all domain event subscribers to the dispatcher.
// Add new event handlers here as the application grows — main.go stays clean.
func RegisterEventHandlers(dispatcher *event.Dispatcher) {
	registerUserEvents(dispatcher)
	registerProductEvents(dispatcher)
}

func registerUserEvents(dispatcher *event.Dispatcher) {
	dispatcher.Subscribe(user.UserCreatedTopic, func(ctx context.Context, ev domain.DomainEvent) error {
		if e, ok := ev.(user.UserCreated); ok {
			slog.Info("[EVENT] User Created",
				"userID", e.UserID,
				"email", e.Email,
				"occurredAt", e.OccurredAt())
		}
		return nil
	})
}

func registerProductEvents(dispatcher *event.Dispatcher) {
	dispatcher.Subscribe(product.ProductCreatedTopic, func(ctx context.Context, ev domain.DomainEvent) error {
		if e, ok := ev.(product.ProductCreated); ok {
			slog.Info("[EVENT] Product Created",
				"productID", e.ProductID,
				"name", e.Name,
				"occurredAt", e.OccurredAt())
		}
		return nil
	})
}
