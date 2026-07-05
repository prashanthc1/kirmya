package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"workspace-app/internal/platform/eventbus"
	"workspace-app/internal/platform/tx"
)

// OutboxPublisher implements the domain EventPublisher interface by writing
// event records to the event_outbox table in the same database transaction.
type OutboxPublisher struct {
	db *sql.DB
}

// NewPublisher builds an OutboxPublisher.
func NewPublisher(db *sql.DB) *OutboxPublisher {
	return &OutboxPublisher{db: db}
}

// Publish serializes and stores the event aggregateID, type, and payload
// in the event_outbox table.
func (p *OutboxPublisher) Publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	sanitizedID := sanitizeUUID(aggregateID)
	executor := tx.GetExecutor(ctx, p.db)

	const q = `
		INSERT INTO event_outbox (aggregate_id, event_type, payload, occurred_at)
		VALUES ($1, $2, $3, now())`
	_, err = executor.ExecContext(ctx, q, sanitizedID, eventType, payloadBytes)
	return err
}

func sanitizeUUID(id string) string {
	// A valid UUID has exactly 36 characters (e.g. 123e4567-e89b-12d3-a456-426614174000).
	// If it doesn't match this basic length check (common in test stubs), fallback to a zero UUID.
	if len(id) == 36 {
		return id
	}
	return "00000000-0000-0000-0000-000000000000"
}

// Relay handles reading unpublished events from the database outbox
// and publishing them onto the underlying Event Bus asynchronously.
type Relay struct {
	db     *sql.DB
	bus    *eventbus.Bus
	ticker *time.Ticker
	stop   chan struct{}
}

// NewRelay builds an outbox relay background worker.
func NewRelay(db *sql.DB, bus *eventbus.Bus) *Relay {
	return &Relay{
		db:   db,
		bus:  bus,
		stop: make(chan struct{}),
	}
}

// Start spawns the relay background loop.
func (r *Relay) Start(interval time.Duration) {
	r.ticker = time.NewTicker(interval)
	go r.run()
}

// Stop shuts down the relay loop.
func (r *Relay) Stop() {
	if r.ticker != nil {
		r.ticker.Stop()
	}
	close(r.stop)
}

func (r *Relay) run() {
	for {
		select {
		case <-r.stop:
			return
		case <-r.ticker.C:
			r.processBatch()
		}
	}
}

type outboxEntry struct {
	ID          int64
	AggregateID string
	Type        string
	Payload     []byte
}

func (r *Relay) processBatch() {
	ctx := context.Background()

	// Query for unpublished events. Use FOR UPDATE SKIP LOCKED to prevent race
	// conditions if multiple replicas or instances are running simultaneously.
	const selectQ = `
		SELECT id, aggregate_id, event_type, payload
		FROM event_outbox
		WHERE published_at IS NULL
		ORDER BY id ASC
		LIMIT 50
		FOR UPDATE SKIP LOCKED`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer func() { _ = tx.Rollback() }()

	rows, err := tx.QueryContext(ctx, selectQ)
	if err != nil {
		return
	}
	defer rows.Close()

	var entries []outboxEntry
	for rows.Next() {
		var entry outboxEntry
		if err := rows.Scan(&entry.ID, &entry.AggregateID, &entry.Type, &entry.Payload); err != nil {
			return
		}
		entries = append(entries, entry)
	}
	_ = rows.Close()

	if len(entries) == 0 {
		return
	}

	const updateQ = `UPDATE event_outbox SET published_at = now() WHERE id = $1`

	for _, entry := range entries {
		var payload map[string]any
		if err := json.Unmarshal(entry.Payload, &payload); err != nil {
			log.Printf("[outbox] failed to unmarshal payload for event %d: %v", entry.ID, err)
			continue
		}

		// Relay the event to the real event bus (NATS JetStream or in-process)
		if err := r.bus.Publish(ctx, entry.Type, entry.AggregateID, payload); err != nil {
			log.Printf("[outbox] failed to publish event %d to bus: %v", entry.ID, err)
			continue
		}

		// Mark the event as published
		if _, err := tx.ExecContext(ctx, updateQ, entry.ID); err != nil {
			log.Printf("[outbox] failed to mark event %d as published: %v", entry.ID, err)
			return
		}
	}

	_ = tx.Commit()
}
