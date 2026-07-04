package eventbus

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"
)

func TestInProcessPubSub(t *testing.T) {
	os.Unsetenv("NATS_URL")
	bus := New()

	var mu sync.Mutex
	got := map[string]string{}
	bus.Subscribe("Thing.Happened", func(_ context.Context, e Event) {
		mu.Lock()
		got["a"] = e.AggregateID
		mu.Unlock()
	})
	bus.Subscribe("Thing.Happened", func(_ context.Context, e Event) {
		mu.Lock()
		got["b"], _ = e.Payload["k"].(string)
		mu.Unlock()
	})

	if err := bus.Publish(context.Background(), "Thing.Happened", "agg-1", map[string]any{"k": "v"}); err != nil {
		t.Fatalf("publish: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if got["a"] != "agg-1" || got["b"] != "v" {
		t.Fatalf("both handlers should fire, got %+v", got)
	}
}

func TestHandlerPanicIsolated(t *testing.T) {
	os.Unsetenv("NATS_URL")
	bus := New()

	ran := false
	bus.Subscribe("E", func(context.Context, Event) { panic("boom") })
	bus.Subscribe("E", func(context.Context, Event) { ran = true })

	if err := bus.Publish(context.Background(), "E", "1", nil); err != nil {
		t.Fatalf("publish: %v", err)
	}
	if !ran {
		t.Fatal("a panicking handler must not stop the others")
	}
}

// TestNATSIntegration exercises the JetStream path. It runs only when NATS_URL
// is set (e.g. against a `nats -js` container).
func TestNATSIntegration(t *testing.T) {
	if os.Getenv("NATS_URL") == "" {
		t.Skip("NATS_URL not set; skipping JetStream integration test")
	}
	bus := New()
	defer bus.Close()
	if bus.js == nil {
		t.Fatal("expected JetStream to be connected when NATS_URL is set")
	}

	done := make(chan Event, 1)
	bus.Subscribe("JobPosted", func(_ context.Context, e Event) { done <- e })

	if err := bus.Publish(context.Background(), "JobPosted", "job-42", map[string]any{"title": "Ops Manager"}); err != nil {
		t.Fatalf("publish: %v", err)
	}

	select {
	case e := <-done:
		if e.AggregateID != "job-42" || e.Payload["title"] != "Ops Manager" {
			t.Fatalf("unexpected event delivered: %+v", e)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for the event via NATS")
	}
}
