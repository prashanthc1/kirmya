package application

import (
	"testing"
	"time"

	"workspace-app/internal/messaging/domain"
)

func TestMessageHubDeliversToParticipant(t *testing.T) {
	h := NewHub(nil)
	ch, cancel := h.Subscribe("u1")
	defer cancel()

	h.Publish("u1", StreamEvent{Kind: EventMessage, ConversationID: "c1", Message: &domain.Message{ID: "m1", Body: "hi"}})
	h.Publish("u2", StreamEvent{Kind: EventMessage, ConversationID: "c9", Message: &domain.Message{ID: "m2"}})

	select {
	case ev := <-ch:
		if ev.Kind != EventMessage || ev.ConversationID != "c1" || ev.Message.ID != "m1" {
			t.Fatalf("unexpected event %+v", ev)
		}
	case <-time.After(time.Second):
		t.Fatal("expected an event")
	}
	select {
	case ev := <-ch:
		t.Fatalf("did not expect another event, got %+v", ev)
	default:
	}
}

func TestMessageHubTypingAndRead(t *testing.T) {
	h := NewHub(nil)
	ch, cancel := h.Subscribe("u1")
	defer cancel()

	h.Publish("u1", StreamEvent{Kind: EventTyping, ConversationID: "c1", ActorID: "u2"})
	h.Publish("u1", StreamEvent{Kind: EventRead, ConversationID: "c1", ActorID: "u2", At: time.Now()})

	if ev := <-ch; ev.Kind != EventTyping || ev.ActorID != "u2" {
		t.Fatalf("expected typing from u2, got %+v", ev)
	}
	if ev := <-ch; ev.Kind != EventRead || ev.ActorID != "u2" {
		t.Fatalf("expected read from u2, got %+v", ev)
	}
}

func TestMessageHubCancelCloses(t *testing.T) {
	h := NewHub(nil)
	ch, cancel := h.Subscribe("u1")
	cancel()
	if _, ok := <-ch; ok {
		t.Fatal("expected channel closed after cancel")
	}
	h.Publish("u1", StreamEvent{Kind: EventMessage, Message: &domain.Message{ID: "late"}}) // must not panic
}
