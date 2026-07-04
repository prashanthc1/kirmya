// Package notifications is the composition root for the Notifications context.
// It also subscribes to the platform event bus so cross-module events
// (referrals, messages, mentorship) generate user notifications.
package notifications

import (
	"context"
	"database/sql"
	"net/http"

	"workspace-app/internal/notifications/api"
	"workspace-app/internal/notifications/application"
	"workspace-app/internal/notifications/infrastructure/postgres"
	"workspace-app/internal/platform/eventbus"
)

// RegisterRoutes wires the module and subscribes to domain events.
func RegisterRoutes(mux *http.ServeMux, db *sql.DB, authMiddleware func(http.Handler) http.Handler, bus *eventbus.Bus, prefs application.PrefChecker) {
	svc := application.NewService(postgres.NewRepository(db), application.NewHub(bus))
	svc.SetPrefChecker(prefs)
	api.RegisterRoutes(mux, api.NewHandler(svc), authMiddleware)
	if bus != nil {
		subscribe(bus, svc)
	}
}

func subscribe(bus *eventbus.Bus, svc *application.Service) {
	bus.Subscribe("ReferralRequested", func(ctx context.Context, e eventbus.Event) {
		_ = svc.Notify(ctx, str(e.Payload, "referrer_id"), "referral_request",
			"New referral request", "Someone asked you for a referral.", "/referrals")
	})
	bus.Subscribe("ReferralAccepted", func(ctx context.Context, e eventbus.Event) {
		_ = svc.Notify(ctx, str(e.Payload, "seeker_id"), "referral_accepted",
			"Referral accepted 🎉", "A referrer accepted your request.", "/referrals")
	})
	bus.Subscribe("ReferralDeclined", func(ctx context.Context, e eventbus.Event) {
		_ = svc.Notify(ctx, str(e.Payload, "seeker_id"), "referral_declined",
			"Referral update", "A referrer responded to your request.", "/referrals")
	})
	bus.Subscribe("MentorshipBooked", func(ctx context.Context, e eventbus.Event) {
		_ = svc.Notify(ctx, str(e.Payload, "mentor_user_id"), "mentorship_booked",
			"New session request", "Someone requested a mentorship session.", "/mentorship")
	})
	bus.Subscribe("MessageSent", func(ctx context.Context, e eventbus.Event) {
		sender := str(e.Payload, "sender_id")
		recips, _ := e.Payload["recipient_ids"].([]string)
		for _, uid := range recips {
			if uid != "" && uid != sender {
				_ = svc.Notify(ctx, uid, "message", "New message", "You have a new message.", "/messages")
			}
		}
	})
}

func str(p map[string]any, key string) string {
	if v, ok := p[key].(string); ok {
		return v
	}
	return ""
}
