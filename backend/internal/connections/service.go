package connections

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"workspace-app/internal/platform/tx"
)

// EventPublisher is structurally identical to outbox.Bus or eventbus.Bus
type EventPublisher interface {
	Publish(ctx context.Context, eventType, aggregateID string, payload map[string]any) error
}

type Service struct {
	repo  *Repository
	txMgr *tx.TxManager
	bus   EventPublisher
	rdb   *redis.Client

	// In-memory rate limiter fallback for testing and environments without Redis
	mu          sync.Mutex
	localCounts map[string]int
	localExpiry map[string]time.Time
}

func NewService(db *sql.DB, repo *Repository, bus EventPublisher) *Service {
	var rdb *redis.Client
	url := os.Getenv("REDIS_URL")
	if url == "" {
		if addr := os.Getenv("REDIS_ADDR"); addr != "" {
			url = "redis://" + addr
		}
	}
	if url != "" {
		opts, err := redis.ParseURL(url)
		if err == nil {
			rdb = redis.NewClient(opts)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := rdb.Ping(ctx).Err(); err != nil {
				rdb = nil
			}
		}
	}

	return &Service{
		repo:        repo,
		txMgr:       tx.NewTxManager(db),
		bus:         bus,
		rdb:         rdb,
		localCounts: make(map[string]int),
		localExpiry: make(map[string]time.Time),
	}
}

// SetRedisClient allows overriding the redis client for testing
func (s *Service) SetRedisClient(client *redis.Client) {
	s.rdb = client
}

// CanMessage checks if two users have an accepted connection
func CanMessage(ctx context.Context, db *sql.DB, userA, userB string) (bool, error) {
	if userA == userB {
		return false, nil
	}
	repo := NewRepository(db)
	conn, err := repo.GetConnection(ctx, userA, userB)
	if errors.Is(err, ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return conn.Status == StatusAccepted, nil
}

// SendRequest sends a connection request from one user to another
func (s *Service) SendRequest(ctx context.Context, fromUser, toUser string, note *string, source *ConnectionSource) error {
	if fromUser == toUser {
		return ErrSelfRequest
	}

	// 1. Rate limiting check
	if err := s.checkRateLimit(ctx, fromUser); err != nil {
		return err
	}

	return s.txMgr.RunInTx(ctx, func(ctx context.Context) error {
		// 2. Block check (both directions)
		blocked, err := s.repo.IsBlocked(ctx, fromUser, toUser)
		if err != nil {
			return err
		}
		if blocked {
			return ErrBlocked
		}

		// 3. Existing connection row checks
		conn, err := s.repo.GetConnection(ctx, fromUser, toUser)
		if err == nil {
			if conn.Status == StatusPending {
				return ErrAlreadyPending
			}
			if conn.Status == StatusAccepted {
				return ErrAlreadyConnected
			}
			if conn.Status == StatusBlocked {
				return ErrBlocked
			}
			if conn.Status == StatusDeclined {
				// Cooldown check: must be older than 30 days
				cooldownLimit := time.Now().AddDate(0, 0, -30)
				if conn.RespondedAt != nil && conn.RespondedAt.After(cooldownLimit) {
					return ErrCooldown
				}
				// If cooldown has expired, we can re-request.
				// We update the existing row rather than inserting to maintain unique constraint
				err = s.repo.UpdateConnectionStatus(ctx, conn.ID, StatusPending, nil)
				if err != nil {
					return err
				}
				// Force reset requested_by, created_at, and updated_at
				exec := tx.GetExecutor(ctx, s.repo.db)
				_, err = exec.ExecContext(ctx, `
					UPDATE connections 
					SET requested_by = $1, created_at = now(), updated_at = now() 
					WHERE id = $2
				`, fromUser, conn.ID)
				if err != nil {
					return err
				}
				// Update or insert request meta
				if note != nil || source != nil {
					_, err = exec.ExecContext(ctx, `
						INSERT INTO connection_requests_meta (connection_id, note, source)
						VALUES ($1, $2, $3)
						ON CONFLICT (connection_id) DO UPDATE
						SET note = EXCLUDED.note, source = EXCLUDED.source
					`, conn.ID, note, source)
					if err != nil {
						return err
					}
				}
				// Reconcile counts
				if err := s.repo.ReconcileCounts(ctx, fromUser); err != nil {
					return err
				}
				if err := s.repo.ReconcileCounts(ctx, toUser); err != nil {
					return err
				}
				// Publish notification event
				if s.bus != nil {
					_ = s.bus.Publish(ctx, "ConnectionRequested", conn.ID, map[string]any{
						"connection_id": conn.ID,
						"requester_id":  fromUser,
						"receiver_id":   toUser,
					})
				}
				return nil
			}
		} else if !errors.Is(err, ErrNotFound) {
			return err
		}

		// 4. Create new connection row
		newConn, err := s.repo.CreateConnection(ctx, fromUser, toUser, StatusPending, note, source)
		if err != nil {
			return err
		}

		// 5. Reconcile counts
		if err := s.repo.ReconcileCounts(ctx, fromUser); err != nil {
			return err
		}
		if err := s.repo.ReconcileCounts(ctx, toUser); err != nil {
			return err
		}

		// 6. Publish notification event
		if s.bus != nil {
			_ = s.bus.Publish(ctx, "ConnectionRequested", newConn.ID, map[string]any{
				"connection_id": newConn.ID,
				"requester_id":  fromUser,
				"receiver_id":   toUser,
			})
		}

		return nil
	})
}

// AcceptRequest accepts a pending connection request
func (s *Service) AcceptRequest(ctx context.Context, connectionID, respondingUser string) error {
	return s.txMgr.RunInTx(ctx, func(ctx context.Context) error {
		conn, err := s.repo.GetConnectionByID(ctx, connectionID)
		if err != nil {
			return err
		}

		// Verify responder is part of the connection and NOT the original requester
		if (conn.UserAID != respondingUser && conn.UserBID != respondingUser) || conn.RequestedBy == respondingUser {
			return ErrForbidden
		}

		if conn.Status != StatusPending {
			return ErrForbidden
		}

		now := time.Now().UTC()
		err = s.repo.UpdateConnectionStatus(ctx, connectionID, StatusAccepted, &now)
		if err != nil {
			return err
		}

		// Reconcile counts
		if err := s.repo.ReconcileCounts(ctx, conn.UserAID); err != nil {
			return err
		}
		if err := s.repo.ReconcileCounts(ctx, conn.UserBID); err != nil {
			return err
		}

		// Publish event
		if s.bus != nil {
			_ = s.bus.Publish(ctx, "ConnectionAccepted", connectionID, map[string]any{
				"connection_id": connectionID,
				"requester_id":  conn.RequestedBy,
				"receiver_id":   respondingUser,
			})
		}

		return nil
	})
}

// DeclineRequest declines a pending connection request
func (s *Service) DeclineRequest(ctx context.Context, connectionID, respondingUser string) error {
	return s.txMgr.RunInTx(ctx, func(ctx context.Context) error {
		conn, err := s.repo.GetConnectionByID(ctx, connectionID)
		if err != nil {
			return err
		}

		// Verify responder is part of the connection and NOT the original requester
		if (conn.UserAID != respondingUser && conn.UserBID != respondingUser) || conn.RequestedBy == respondingUser {
			return ErrForbidden
		}

		if conn.Status != StatusPending {
			return ErrForbidden
		}

		now := time.Now().UTC()
		err = s.repo.UpdateConnectionStatus(ctx, connectionID, StatusDeclined, &now)
		if err != nil {
			return err
		}

		// Reconcile counts
		if err := s.repo.ReconcileCounts(ctx, conn.UserAID); err != nil {
			return err
		}
		if err := s.repo.ReconcileCounts(ctx, conn.UserBID); err != nil {
			return err
		}

		return nil
	})
}

// RemoveConnection removes an existing accepted connection
func (s *Service) RemoveConnection(ctx context.Context, connectionID, requestingUser string) error {
	return s.txMgr.RunInTx(ctx, func(ctx context.Context) error {
		conn, err := s.repo.GetConnectionByID(ctx, connectionID)
		if err != nil {
			return err
		}

		// Verify requesting user is part of the connection
		if conn.UserAID != requestingUser && conn.UserBID != requestingUser {
			return ErrForbidden
		}

		if conn.Status != StatusAccepted {
			return ErrForbidden
		}

		err = s.repo.DeleteConnection(ctx, connectionID)
		if err != nil {
			return err
		}

		// Reconcile counts
		if err := s.repo.ReconcileCounts(ctx, conn.UserAID); err != nil {
			return err
		}
		if err := s.repo.ReconcileCounts(ctx, conn.UserBID); err != nil {
			return err
		}

		return nil
	})
}

// BlockUser blocks a user, cancelling any pending/accepted connection
func (s *Service) BlockUser(ctx context.Context, blocker, blocked, reason string) error {
	if blocker == blocked {
		return ErrSelfRequest
	}

	return s.txMgr.RunInTx(ctx, func(ctx context.Context) error {
		// Idempotently insert block row
		err := s.repo.InsertBlock(ctx, blocker, blocked, reason)
		if err != nil {
			return err
		}

		// Cancel any existing connection (set status=blocked, requested_by=blocker)
		conn, err := s.repo.GetConnection(ctx, blocker, blocked)
		if err == nil {
			err = s.repo.UpdateConnectionStatus(ctx, conn.ID, StatusBlocked, nil)
			if err != nil {
				return err
			}
			// Force update requested_by to blocker
			exec := tx.GetExecutor(ctx, s.repo.db)
			_, err = exec.ExecContext(ctx, "UPDATE connections SET requested_by = $1, updated_at = now() WHERE id = $2", blocker, conn.ID)
			if err != nil {
				return err
			}
		} else if errors.Is(err, ErrNotFound) {
			// Create a blocked connection row if none existed
			_, err = s.repo.CreateConnection(ctx, blocker, blocked, StatusBlocked, nil, nil)
			if err != nil {
				return err
			}
		} else {
			return err
		}

		// Reconcile counts
		if err := s.repo.ReconcileCounts(ctx, blocker); err != nil {
			return err
		}
		if err := s.repo.ReconcileCounts(ctx, blocked); err != nil {
			return err
		}

		return nil
	})
}

// UnblockUser removes a block, returning the pair to a "no relationship" state
func (s *Service) UnblockUser(ctx context.Context, blocker, blocked string) error {
	if blocker == blocked {
		return ErrSelfRequest
	}

	return s.txMgr.RunInTx(ctx, func(ctx context.Context) error {
		err := s.repo.DeleteBlock(ctx, blocker, blocked)
		if err != nil {
			return err
		}

		// Delete the connections row that represented the block
		conn, err := s.repo.GetConnection(ctx, blocker, blocked)
		if err == nil {
			if conn.Status == StatusBlocked {
				err = s.repo.DeleteConnection(ctx, conn.ID)
				if err != nil {
					return err
				}
			}
		} else if !errors.Is(err, ErrNotFound) {
			return err
		}

		// Reconcile counts
		if err := s.repo.ReconcileCounts(ctx, blocker); err != nil {
			return err
		}
		if err := s.repo.ReconcileCounts(ctx, blocked); err != nil {
			return err
		}

		return nil
	})
}

// checkRateLimit enforces the 20 requests per 24h limit
func (s *Service) checkRateLimit(ctx context.Context, userID string) error {
	dateStr := time.Now().UTC().Format("2006-01-02")
	key := fmt.Sprintf("conn_req:%s:%s", userID, dateStr)

	// 1. If Redis is available, run INCR + EXPIRE
	if s.rdb != nil {
		val, err := s.rdb.Incr(ctx, key).Result()
		if err != nil {
			// Fallback to local memory limiter if Redis throws errors
			return s.checkRateLimitLocal(userID)
		}
		if val == 1 {
			s.rdb.Expire(ctx, key, 24*time.Hour)
		}
		if val > 20 {
			return ErrRateLimited
		}
		return nil
	}

	// 2. Fallback to local in-memory thread-safe rate limiter
	return s.checkRateLimitLocal(userID)
}

func (s *Service) checkRateLimitLocal(userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	dateStr := now.UTC().Format("2006-01-02")
	key := fmt.Sprintf("conn_req:%s:%s", userID, dateStr)

	// Clean up stale local entries periodically
	for k, exp := range s.localExpiry {
		if now.After(exp) {
			delete(s.localCounts, k)
			delete(s.localExpiry, k)
		}
	}

	count := s.localCounts[key]
	if count >= 20 {
		return ErrRateLimited
	}

	s.localCounts[key] = count + 1
	if count == 0 {
		s.localExpiry[key] = now.Add(24 * time.Hour)
	}

	return nil
}
