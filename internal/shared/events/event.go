package events

import (
	"context"
	"time"
)

type Event interface {
	EventType() string
	OccurredAt() time.Time
}

type EventHandler func(ctx context.Context, event Event) error

type EventBus interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(eventType string, handler EventHandler)
	Start(ctx context.Context) error
	Stop() error
}
