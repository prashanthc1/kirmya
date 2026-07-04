package application

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"workspace-app/internal/messaging/domain"
)

const fanoutSubject = "sse.messages"

// Broadcaster is the cross-instance fanout transport (the platform event bus
// satisfies it). When NATS is available, conversation events are broadcast so
// whichever backend instance holds a participant's SSE connection delivers them.
type Broadcaster interface {
	HasNATS() bool
	BroadcastFanout(subject string, data []byte)
	SubscribeFanout(subject string, handler func([]byte))
}

type fanoutEnvelope struct {
	UserID string      `json:"user_id"`
	Event  StreamEvent `json:"event"`
}

// Stream event kinds delivered over the conversation SSE stream.
const (
	EventMessage = "message"
	EventTyping  = "typing"
	EventRead    = "read"
)

// StreamEvent is the envelope pushed to a participant's live conversation
// stream. ActorID is the message sender, the typing user, or the reader,
// depending on Kind. Message is set only for Kind == EventMessage.
type StreamEvent struct {
	Kind           string
	ConversationID string
	ActorID        string
	Message        *domain.Message
	At             time.Time
}

// Hub fans conversation events out to a participant's connected SSE subscribers,
// keyed by user id. With a NATS-backed Broadcaster it works across instances;
// otherwise it is a single-instance in-memory hub. Delivery is best-effort and
// non-blocking (drops on a full buffer; the client re-syncs on reload).
type Hub struct {
	mu    sync.RWMutex
	subs  map[string]map[chan StreamEvent]struct{}
	bcast Broadcaster
}

func NewHub(bcast Broadcaster) *Hub {
	h := &Hub{subs: map[string]map[chan StreamEvent]struct{}{}, bcast: bcast}
	if bcast != nil && bcast.HasNATS() {
		bcast.SubscribeFanout(fanoutSubject, func(data []byte) {
			var env fanoutEnvelope
			if json.Unmarshal(data, &env) == nil {
				h.localPublish(env.UserID, env.Event)
			}
		})
		log.Printf("[messaging] cross-instance SSE fanout enabled")
	}
	return h
}

// Subscribe registers a subscriber for a user and returns its channel plus a
// cancel func that must be called when the subscriber goes away.
func (h *Hub) Subscribe(userID string) (<-chan StreamEvent, func()) {
	ch := make(chan StreamEvent, 32)
	h.mu.Lock()
	if h.subs[userID] == nil {
		h.subs[userID] = map[chan StreamEvent]struct{}{}
	}
	h.subs[userID][ch] = struct{}{}
	h.mu.Unlock()

	var once sync.Once
	cancel := func() {
		once.Do(func() {
			h.mu.Lock()
			if set, ok := h.subs[userID]; ok {
				delete(set, ch)
				if len(set) == 0 {
					delete(h.subs, userID)
				}
			}
			h.mu.Unlock()
			close(ch)
		})
	}
	return ch, cancel
}

// Publish delivers an event to a user. With NATS it broadcasts to every instance
// (each delivers to its own local subscribers); otherwise it delivers locally.
func (h *Hub) Publish(userID string, ev StreamEvent) {
	if h.bcast != nil && h.bcast.HasNATS() {
		if data, err := json.Marshal(fanoutEnvelope{UserID: userID, Event: ev}); err == nil {
			h.bcast.BroadcastFanout(fanoutSubject, data)
		}
		return
	}
	h.localPublish(userID, ev)
}

// localPublish delivers to this instance's subscribers (non-blocking).
func (h *Hub) localPublish(userID string, ev StreamEvent) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.subs[userID] {
		select {
		case ch <- ev:
		default: // buffer full — drop; the client re-syncs on reload
		}
	}
}
