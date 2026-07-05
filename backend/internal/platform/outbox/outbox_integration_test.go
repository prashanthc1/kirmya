//go:build integration

package outbox

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"workspace-app/internal/platform/eventbus"
	"workspace-app/internal/platform/tx"
	"workspace-app/internal/testsupport"
)

func TestOutbox_TransactionalPublishAndRelay(t *testing.T) {
	db := testsupport.OpenTestDB(t)
	ctx := context.Background()

	bus := eventbus.New()
	pub := NewPublisher(db)
	txMgr := tx.NewTxManager(db)

	const eventType = "TestEventOccurred"
	const aggregateID = "e2c842b0-9f5b-4c07-b08e-324c4786720d" // Valid UUID
	payload := map[string]any{"key": "value"}

	// Scenario 1: Write inside a transaction that gets rolled back
	err := txMgr.RunInTx(ctx, func(txCtx context.Context) error {
		err := pub.Publish(txCtx, eventType, aggregateID, payload)
		if err != nil {
			t.Fatalf("unexpected publish error: %v", err)
		}
		// Return an error to force rollback
		return tx.WithTx(txCtx, nil).Err() // arbitrary error to trigger rollback, or custom error
	})
	if err == nil {
		t.Fatalf("expected error on rollback transaction, got nil")
	}

	// Verify database is empty (no record in event_outbox)
	var count int
	err = db.QueryRowContext(ctx, "SELECT count(*) FROM event_outbox WHERE aggregate_id = $1", aggregateID).Scan(&count)
	if err != nil {
		t.Fatalf("failed to query outbox: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 outbox events after rollback, got %d", count)
	}

	// Scenario 2: Write inside a transaction that commits successfully
	err = txMgr.RunInTx(ctx, func(txCtx context.Context) error {
		return pub.Publish(txCtx, eventType, aggregateID, payload)
	})
	if err != nil {
		t.Fatalf("unexpected transaction failure: %v", err)
	}

	// Verify record exists and is unpublished
	var dbPayload []byte
	var publishedAt *time.Time
	err = db.QueryRowContext(ctx, "SELECT payload, published_at FROM event_outbox WHERE aggregate_id = $1", aggregateID).Scan(&dbPayload, &publishedAt)
	if err != nil {
		t.Fatalf("failed to read outbox row: %v", err)
	}
	if publishedAt != nil {
		t.Fatalf("expected published_at to be nil, got %v", publishedAt)
	}

	var parsedPayload map[string]any
	if err := json.Unmarshal(dbPayload, &parsedPayload); err != nil {
		t.Fatalf("failed to unmarshal db payload: %v", err)
	}
	if parsedPayload["key"] != "value" {
		t.Fatalf("expected payload key 'value', got '%v'", parsedPayload["key"])
	}

	// Scenario 3: Verify UUID sanitization fallback for non-UUID strings
	const invalidUUID = "not-a-uuid"
	err = txMgr.RunInTx(ctx, func(txCtx context.Context) error {
		return pub.Publish(txCtx, eventType, invalidUUID, payload)
	})
	if err != nil {
		t.Fatalf("failed to publish event with invalid UUID: %v", err)
	}

	// Query fallback UUID record
	err = db.QueryRowContext(ctx, "SELECT count(*) FROM event_outbox WHERE aggregate_id = '00000000-0000-0000-0000-000000000000'").Scan(&count)
	if err != nil {
		t.Fatalf("failed to check fallback UUID record: %v", err)
	}
	if count == 0 {
		t.Fatalf("expected fallback record under zero-UUID namespace, but none was found")
	}

	// Scenario 4: Start outbox Relay and verify delivery to EventBus
	delivered := make(chan eventbus.Event, 2)
	bus.Subscribe(eventType, func(ctx context.Context, e eventbus.Event) {
		delivered <- e
	})

	relay := NewRelay(db, bus)
	relay.Start(10 * time.Millisecond)
	defer relay.Stop()

	// Wait for delivery
	select {
	case e := <-delivered:
		if e.AggregateID != aggregateID {
			t.Errorf("expected aggregate ID %s, got %s", aggregateID, e.AggregateID)
		}
		if e.Payload["key"] != "value" {
			t.Errorf("expected payload key 'value', got %v", e.Payload["key"])
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for outbox event delivery to eventbus")
	}

	// Wait a moment for database mark-as-published update to complete
	time.Sleep(50 * time.Millisecond)

	// Verify outbox entry is now marked as published
	err = db.QueryRowContext(ctx, "SELECT published_at FROM event_outbox WHERE aggregate_id = $1", aggregateID).Scan(&publishedAt)
	if err != nil {
		t.Fatalf("failed to read outbox row again: %v", err)
	}
	if publishedAt == nil {
		t.Fatal("expected outbox row to be marked published (published_at not nil)")
	}
}
