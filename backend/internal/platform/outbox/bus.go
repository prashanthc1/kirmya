package outbox

import (
	"context"

	"workspace-app/internal/platform/eventbus"
)

// Bus embeds the native eventbus.Bus to delegate NATS subscriptions and raw SSE
// broadcasts, but overrides the Publish method to write events atomically to the
// outbox table.
type Bus struct {
	*eventbus.Bus
	pub *OutboxPublisher
}

// NewBus creates an Outbox-backed Bus.
func NewBus(bus *eventbus.Bus, pub *OutboxPublisher) *Bus {
	return &Bus{
		Bus: bus,
		pub: pub,
	}
}

// Publish writes the event data to the database event_outbox table.
func (b *Bus) Publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) error {
	return b.pub.Publish(ctx, eventType, aggregateID, payload)
}
