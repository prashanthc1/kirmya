package network

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"workspace-app/internal/network/api"
	"workspace-app/internal/network/application"
	"workspace-app/internal/network/domain"
	"workspace-app/internal/network/infrastructure/postgres"
	"workspace-app/internal/platform/eventbus"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB, authMiddleware func(http.Handler) http.Handler, bus *eventbus.Bus, limit func(string) func(http.Handler) http.Handler) {
	repo := postgres.NewRepository(db)
	svc := application.NewService(repo)
	h := api.NewHandler(svc)
	api.RegisterRoutes(mux, h, authMiddleware, limit)

	if bus != nil {
		subscribeToAutoCreateEvents(db, bus, svc)
	}
}

func subscribeToAutoCreateEvents(db *sql.DB, bus *eventbus.Bus, svc *application.Service) {
	// 1. Mentorship Confirmation
	bus.Subscribe("SessionConfirmed", func(ctx context.Context, e eventbus.Event) {
		sessionID := e.AggregateID
		if sessionID == "" {
			return
		}
		// Query database directly to get mentee_id and mentor's user_id
		var mentorUserID, menteeID string
		err := db.QueryRowContext(ctx, `
			SELECT mp.user_id, ms.mentee_id
			FROM mentorship_sessions ms
			JOIN mentor_profiles mp ON ms.mentor_id = mp.id
			WHERE ms.id = $1
		`, sessionID).Scan(&mentorUserID, &menteeID)

		if err != nil {
			log.Printf("[network-auto-connect] failed to query session %s: %v", sessionID, err)
			return
		}

		_, err = svc.AutoCreateConnection(ctx, mentorUserID, menteeID, domain.OriginMentorshipMatch)
		if err != nil {
			log.Printf("[network-auto-connect] failed to auto-create connection between mentor=%s and mentee=%s: %v", mentorUserID, menteeID, err)
		} else {
			log.Printf("[network-auto-connect] auto-created accepted connection for confirmed mentorship match between %s and %s", mentorUserID, menteeID)
		}
	})

	// 2. Referral Acceptance
	bus.Subscribe("ReferralAccepted", func(ctx context.Context, e eventbus.Event) {
		seekerID, _ := e.Payload["seeker_id"].(string)
		referrerID, _ := e.Payload["referrer_id"].(string)
		if seekerID == "" || referrerID == "" {
			return
		}

		_, err := svc.AutoCreateConnection(ctx, referrerID, seekerID, domain.OriginReferralRequest)
		if err != nil {
			log.Printf("[network-auto-connect] failed to auto-create connection between referrer=%s and seeker=%s: %v", referrerID, seekerID, err)
		} else {
			log.Printf("[network-auto-connect] auto-created accepted connection for accepted referral between referrer=%s and seeker=%s", referrerID, seekerID)
		}
	})
}
