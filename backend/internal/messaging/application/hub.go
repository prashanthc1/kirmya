package application

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"workspace-app/internal/messaging/domain"
)

const (
	redisFanoutChannel = "kirmya:messaging:events"
	presencePrefix     = "presence:"
)

type StreamEvent struct {
	Kind           string          `json:"kind"`
	ConversationID string          `json:"conversation_id"`
	ActorID        string          `json:"actor_id"`
	Message        *domain.Message `json:"message,omitempty"`
	At             time.Time       `json:"at"`
}

const (
	EventMessage = "message"
	EventTyping  = "typing"
	EventRead    = "read"
)

type Broadcaster interface {
	HasNATS() bool
	BroadcastFanout(subject string, data []byte)
	SubscribeFanout(subject string, handler func([]byte))
}

type Hub struct {
	mu          sync.RWMutex
	subs        map[string]map[chan StreamEvent]struct{}
	rdb         *redis.Client
	bcast       Broadcaster
	pubsubClose func() error
}

func NewHub(bcast Broadcaster) *Hub {
	h := &Hub{
		subs:  map[string]map[chan StreamEvent]struct{}{},
		bcast: bcast,
	}

	// Initialize Redis client if REDIS_URL or REDIS_ADDR is present
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		if addr := os.Getenv("REDIS_ADDR"); addr != "" {
			redisURL = "redis://" + addr
		}
	}

	if redisURL != "" {
		opts, err := redis.ParseURL(redisURL)
		if err == nil {
			h.rdb = redis.NewClient(opts)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := h.rdb.Ping(ctx).Err(); err == nil {
				log.Printf("[messaging] Redis pub/sub and presence enabled")
				h.startRedisSubscription()
			} else {
				log.Printf("[messaging] Redis unreachable (%v); falling back to local-only", err)
				h.rdb = nil
			}
		} else {
			log.Printf("[messaging] invalid REDIS_URL (%v); falling back to local-only", err)
		}
	}

	// Fallback to NATS if Redis is not configured but NATS is available
	if h.rdb == nil && bcast != nil && bcast.HasNATS() {
		bcast.SubscribeFanout("sse.messages", func(data []byte) {
			var env fanoutEnvelope
			if json.Unmarshal(data, &env) == nil {
				h.localPublish(env.UserID, env.Event)
			}
		})
		log.Printf("[messaging] cross-instance NATS fanout enabled")
	}

	return h
}

func (h *Hub) startRedisSubscription() {
	pubsub := h.rdb.Subscribe(context.Background(), redisFanoutChannel)
	h.pubsubClose = pubsub.Close

	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			var env fanoutEnvelope
			if err := json.Unmarshal([]byte(msg.Payload), &env); err == nil {
				h.localPublish(env.UserID, env.Event)
			}
		}
	}()
}

type fanoutEnvelope struct {
	UserID string      `json:"user_id"`
	Event  StreamEvent `json:"event"`
}

func (h *Hub) Subscribe(userID string) (<-chan StreamEvent, func()) {
	ch := make(chan StreamEvent, 64)
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

func (h *Hub) Publish(userID string, ev StreamEvent) {
	// 1. Publish to Redis if configured
	if h.rdb != nil {
		env := fanoutEnvelope{UserID: userID, Event: ev}
		data, err := json.Marshal(env)
		if err == nil {
			h.rdb.Publish(context.Background(), redisFanoutChannel, string(data))
		}
		return
	}

	// 2. Publish to NATS if configured
	if h.bcast != nil && h.bcast.HasNATS() {
		env := fanoutEnvelope{UserID: userID, Event: ev}
		if data, err := json.Marshal(env); err == nil {
			h.bcast.BroadcastFanout("sse.messages", data)
		}
		return
	}

	// 3. Fallback to local publish
	h.localPublish(userID, ev)
}

func (h *Hub) localPublish(userID string, ev StreamEvent) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.subs[userID] {
		select {
		case ch <- ev:
		default: // Buffer full, drop event
		}
	}
}

// SetPresence records user online/offline status in Redis.
func (h *Hub) SetPresence(ctx context.Context, userID string, online bool) error {
	if h.rdb == nil {
		return nil
	}

	key := presencePrefix + userID
	if online {
		return h.rdb.Set(ctx, key, "online", 60*time.Second).Err()
	}

	return h.rdb.Del(ctx, key).Err()
}

// GetPresence checks user presence status in Redis.
func (h *Hub) GetPresence(ctx context.Context, userID string) (string, error) {
	if h.rdb == nil {
		return "offline", nil
	}

	val, err := h.rdb.Get(ctx, presencePrefix+userID).Result()
	if err == redis.Nil {
		return "offline", nil
	}
	if err != nil {
		return "offline", err
	}
	return val, nil
}

func (h *Hub) Close() {
	if h.pubsubClose != nil {
		_ = h.pubsubClose()
	}
}
