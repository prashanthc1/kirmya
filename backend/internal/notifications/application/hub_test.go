package application

import (
	"testing"
	"time"

	"workspace-app/internal/notifications/domain"
)

func TestHubDeliversToSubscriber(t *testing.T) {
	h := NewHub(nil)
	ch, cancel := h.Subscribe("u1")
	defer cancel()

	h.Publish(domain.Notification{UserID: "u1", Title: "hi"})
	// Another user's notification must not arrive.
	h.Publish(domain.Notification{UserID: "u2", Title: "nope"})

	select {
	case n := <-ch:
		if n.Title != "hi" {
			t.Fatalf("unexpected notification %+v", n)
		}
	case <-time.After(time.Second):
		t.Fatal("expected a notification")
	}

	select {
	case n := <-ch:
		t.Fatalf("did not expect a second notification, got %+v", n)
	default:
	}
}

func TestHubCancelStopsDelivery(t *testing.T) {
	h := NewHub(nil)
	ch, cancel := h.Subscribe("u1")
	cancel()

	// Channel is closed after cancel.
	if _, ok := <-ch; ok {
		t.Fatal("expected channel to be closed after cancel")
	}
	// Publishing after cancel must not panic (no subscribers).
	h.Publish(domain.Notification{UserID: "u1", Title: "late"})
}

func TestHubNonBlockingWhenBufferFull(t *testing.T) {
	h := NewHub(nil)
	_, cancel := h.Subscribe("u1") // never drained
	defer cancel()

	// Far more than the buffer (16) — Publish must not block or panic.
	done := make(chan struct{})
	go func() {
		for i := 0; i < 100; i++ {
			h.Publish(domain.Notification{UserID: "u1", Title: "x"})
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Publish blocked on a full subscriber buffer")
	}
}
