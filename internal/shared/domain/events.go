package domain

import (
	"time"
)

// DomainEvent is an interface for all domain events
type DomainEvent interface {
	OccurredAt() time.Time
	Topic() string
}

// BaseEvent provides common functionality for domain events
type BaseEvent struct {
	Timestamp time.Time
}

func (e BaseEvent) OccurredAt() time.Time {
	return e.Timestamp
}

// AggregateRoot is a base struct to be embedded in entities that need to record events
type AggregateRoot struct {
	events []DomainEvent
}

func (a *AggregateRoot) RecordEvent(event DomainEvent) {
	a.events = append(a.events, event)
}

func (a *AggregateRoot) GetEvents() []DomainEvent {
	return a.events
}

func (a *AggregateRoot) ClearEvents() {
	a.events = nil
}
