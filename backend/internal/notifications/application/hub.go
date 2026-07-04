package application

import (
	"encoding/json"
	"log"
	"sync"

	"workspace-app/internal/notifications/domain"
)

const fanoutSubject = "sse.notifications"

// Broadcaster is the cross-instance fanout transport (the platform event bus
// satisfies it). When NATS is available, notifications are broadcast so whichever
// backend instance holds a user's SSE connection delivers them — not only the
// instance that created the notification.
type Broadcaster interface {
	HasNATS() bool
	BroadcastFanout(subject string, data []byte)
	SubscribeFanout(subject string, handler func([]byte))
}

type fanoutEnvelope struct {
	UserID       string              `json:"user_id"`
	Notification domain.Notification `json:"notification"`
}

// Hub fans newly-created notifications out to a user's connected SSE subscribers.
// With a NATS-backed Broadcaster it works across instances; otherwise it is a
// single-instance in-memory hub. Delivery is best-effort and non-blocking.
type Hub struct {
	mu    sync.RWMutex
	subs  map[string]map[chan domain.Notification]struct{}
	bcast Broadcaster
}

func NewHub(bcast Broadcaster) *Hub {
	h := &Hub{subs: map[string]map[chan domain.Notification]struct{}{}, bcast: bcast}
	if bcast != nil && bcast.HasNATS() {
		bcast.SubscribeFanout(fanoutSubject, func(data []byte) {
			var env fanoutEnvelope
			if json.Unmarshal(data, &env) == nil {
				h.localPublish(env.UserID, env.Notification)
			}
		})
		log.Printf("[notifications] cross-instance SSE fanout enabled")
	}
	return h
}

// Subscribe registers a subscriber for a user and returns its channel plus a
// cancel func that must be called when the subscriber goes away.
func (h *Hub) Subscribe(userID string) (<-chan domain.Notification, func()) {
	ch := make(chan domain.Notification, 16)
	h.mu.Lock()
	if h.subs[userID] == nil {
		h.subs[userID] = map[chan domain.Notification]struct{}{}
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

// Publish delivers a notification. With NATS it broadcasts to every instance
// (each delivers to its own local subscribers, so it's delivered exactly once
// per connection); otherwise it delivers to local subscribers directly.
func (h *Hub) Publish(n domain.Notification) {
	if h.bcast != nil && h.bcast.HasNATS() {
		if data, err := json.Marshal(fanoutEnvelope{UserID: n.UserID, Notification: n}); err == nil {
			h.bcast.BroadcastFanout(fanoutSubject, data)
		}
		return
	}
	h.localPublish(n.UserID, n)
}

// localPublish delivers to this instance's subscribers (non-blocking).
func (h *Hub) localPublish(userID string, n domain.Notification) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.subs[userID] {
		select {
		case ch <- n:
		default: // buffer full — drop; the client re-syncs on reconnect
		}
	}
}
