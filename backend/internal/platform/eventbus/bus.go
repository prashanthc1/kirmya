// Package eventbus is the platform's publish/subscribe seam. By default it runs
// in-process; when NATS_URL is set it bridges through NATS JetStream so events
// are delivered across replicas (durable, at-least-once) without changing the
// Subscribe/Publish API that modules already use. If NATS is configured but
// unreachable it logs and falls back to in-process delivery.
package eventbus

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// Event is the envelope delivered to subscribers.
type Event struct {
	Type        string         `json:"type"`
	AggregateID string         `json:"aggregate_id"`
	OccurredAt  time.Time      `json:"occurred_at"`
	Payload     map[string]any `json:"payload"`
}

// Handler consumes an event. Handlers run synchronously in registration order;
// a panic in one handler does not stop the others.
type Handler func(ctx context.Context, e Event)

const (
	streamName  = "EVENTS"
	subjectRoot = "events." // subject per type: events.<EventType>
	queueGroup  = "kirmya"
)

// Bus is a concurrency-safe pub/sub bus. With NATS configured, Publish writes to
// JetStream and a single queue subscription per event type dispatches to all
// local handlers (so each event is processed once cluster-wide while preserving
// in-process fan-out to every handler).
type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
	subbed   map[string]bool // event types with an active NATS subscription

	nc *nats.Conn
	js nats.JetStreamContext
}

// New builds the bus. It connects to NATS JetStream when NATS_URL is set and
// reachable; otherwise it runs in-process.
func New() *Bus {
	b := &Bus{handlers: map[string][]Handler{}, subbed: map[string]bool{}}

	url := os.Getenv("NATS_URL")
	if url == "" {
		log.Printf("[eventbus] NATS_URL not set; using in-process event bus")
		return b
	}

	nc, err := nats.Connect(url,
		nats.Name("kirmya"),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		log.Printf("[eventbus] NATS unreachable (%v); falling back to in-process", err)
		return b
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Printf("[eventbus] JetStream init failed (%v); falling back to in-process", err)
		nc.Close()
		return b
	}
	// Ensure the stream exists (idempotent — ignore "already in use").
	if _, err := js.AddStream(&nats.StreamConfig{
		Name:      streamName,
		Subjects:  []string{subjectRoot + ">"},
		Storage:   nats.FileStorage,
		Retention: nats.LimitsPolicy,
	}); err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.Printf("[eventbus] stream setup failed (%v); falling back to in-process", err)
		nc.Close()
		return b
	}

	b.nc, b.js = nc, js
	log.Printf("[eventbus] connected to NATS JetStream at %s", url)
	return b
}

// Subscribe registers a handler for an event type. When NATS is enabled, the
// first handler for a type also opens a durable queue subscription that fans the
// received message out to every local handler for that type.
func (b *Bus) Subscribe(eventType string, h Handler) {
	b.mu.Lock()
	b.handlers[eventType] = append(b.handlers[eventType], h)
	needSub := b.js != nil && !b.subbed[eventType]
	if needSub {
		b.subbed[eventType] = true
	}
	b.mu.Unlock()

	if needSub {
		b.openNatsSub(eventType)
	}
}

func (b *Bus) openNatsSub(eventType string) {
	subject := subjectRoot + eventType
	_, err := b.js.QueueSubscribe(subject, queueGroup, func(m *nats.Msg) {
		var e Event
		if err := json.Unmarshal(m.Data, &e); err == nil {
			b.dispatch(context.Background(), e)
		}
		_ = m.Ack()
	}, nats.Durable("cb_"+eventType), nats.ManualAck(), nats.AckExplicit())
	if err != nil {
		log.Printf("[eventbus] subscribe %s failed: %v", subject, err)
	}
}

// Publish delivers an event. With NATS enabled it writes to JetStream (delivery
// to handlers happens on receive, cluster-wide); otherwise it dispatches to the
// local handlers directly.
func (b *Bus) Publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) error {
	e := Event{Type: eventType, AggregateID: aggregateID, OccurredAt: time.Now().UTC(), Payload: payload}

	if b.js != nil {
		data, err := json.Marshal(e)
		if err != nil {
			return err
		}
		if _, err := b.js.Publish(subjectRoot+eventType, data); err != nil {
			log.Printf("[eventbus] publish %s failed: %v", eventType, err)
			return err
		}
		return nil
	}

	b.dispatch(ctx, e)
	return nil
}

// dispatch invokes all local handlers registered for the event's type.
func (b *Bus) dispatch(ctx context.Context, e Event) {
	b.mu.RLock()
	handlers := append([]Handler(nil), b.handlers[e.Type]...)
	b.mu.RUnlock()

	if len(handlers) == 0 {
		log.Printf("[event] %s aggregate=%s (no subscribers)", e.Type, e.AggregateID)
		return
	}
	for _, h := range handlers {
		safeDispatch(ctx, h, e)
	}
}

func safeDispatch(ctx context.Context, h Handler, e Event) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[event] handler panic for %s: %v", e.Type, r)
		}
	}()
	h(ctx, e)
}

// Close drains the NATS connection (no-op in-process).
func (b *Bus) Close() {
	if b.nc != nil {
		_ = b.nc.Drain()
	}
}

// ----- Fanout (core NATS) -------------------------------------------------
//
// These power cross-instance SSE delivery. Unlike Publish/Subscribe (JetStream
// queue groups, processed once cluster-wide), fanout uses plain core NATS with
// no queue group, so EVERY connected instance receives EVERY message — exactly
// what's needed to reach a user's SSE connection wherever it happens to live.

// HasNATS reports whether a NATS connection is available for fanout. Safe to
// call on a nil *Bus (returns false).
func (b *Bus) HasNATS() bool { return b != nil && b.nc != nil }

// BroadcastFanout publishes raw bytes to a fanout subject (no-op without NATS).
func (b *Bus) BroadcastFanout(subject string, data []byte) {
	if b.HasNATS() {
		if err := b.nc.Publish(subject, data); err != nil {
			log.Printf("[eventbus] fanout publish %s: %v", subject, err)
		}
	}
}

// SubscribeFanout subscribes to a fanout subject with no queue group, so this
// instance receives every message published to it (no-op without NATS).
func (b *Bus) SubscribeFanout(subject string, handler func(data []byte)) {
	if !b.HasNATS() {
		return
	}
	if _, err := b.nc.Subscribe(subject, func(m *nats.Msg) { handler(m.Data) }); err != nil {
		log.Printf("[eventbus] fanout subscribe %s: %v", subject, err)
	}
}
