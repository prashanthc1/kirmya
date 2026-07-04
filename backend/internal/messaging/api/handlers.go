package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"workspace-app/internal/common"
	"workspace-app/internal/messaging/application"
	"workspace-app/internal/messaging/domain"
)

type Handler struct{ svc *application.Service }

func NewHandler(svc *application.Service) *Handler { return &Handler{svc: svc} }

// streamEventDTO is a multiplexed conversation event. Kind is "message",
// "typing", or "read"; fields are populated per kind.
type streamEventDTO struct {
	Kind           string `json:"kind"`
	ConversationID string `json:"conversation_id"`
	SenderID       string `json:"sender_id,omitempty"`
	ReaderID       string `json:"reader_id,omitempty"`
	ID             string `json:"id,omitempty"`
	Body           string `json:"body,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	At             string `json:"at,omitempty"`
}

func toStreamEvent(ev application.StreamEvent) streamEventDTO {
	d := streamEventDTO{Kind: ev.Kind, ConversationID: ev.ConversationID}
	switch ev.Kind {
	case application.EventMessage:
		d.SenderID = ev.Message.SenderID
		d.ID = ev.Message.ID
		d.Body = ev.Message.Body
		d.CreatedAt = ev.Message.CreatedAt.UTC().Format(time.RFC3339)
	case application.EventTyping:
		d.SenderID = ev.ActorID
	case application.EventRead:
		d.ReaderID = ev.ActorID
		d.At = ev.At.UTC().Format(time.RFC3339)
	}
	return d
}

// Stream handles GET /conversations/stream — an SSE stream of incoming messages
// across the user's conversations. Clients connect with fetch (Bearer header).
func (h *Handler) Stream(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	rc := http.NewResponseController(w)
	_ = rc.SetWriteDeadline(time.Time{}) // long-lived: clear the server WriteTimeout

	ch, cancel := h.svc.Subscribe(uid)
	defer cancel()

	fmt.Fprint(w, ": connected\n\n")
	if rc.Flush() != nil {
		return
	}

	ping := time.NewTicker(25 * time.Second)
	defer ping.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case ev, ok := <-ch:
			if !ok {
				return
			}
			data, _ := json.Marshal(toStreamEvent(ev))
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ev.Kind, data)
			if rc.Flush() != nil {
				return
			}
		case <-ping.C:
			fmt.Fprint(w, ": ping\n\n")
			if rc.Flush() != nil {
				return
			}
		}
	}
}

func decode(r *http.Request, dst any) bool {
	return json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20)).Decode(dst) == nil
}

type conversationDTO struct {
	ID           string   `json:"id"`
	IsGroup      bool     `json:"is_group"`
	Title        string   `json:"title,omitempty"`
	Participants []string `json:"participants"`
	UpdatedAt    string   `json:"updated_at"`
}

type messageDTO struct {
	ID        string `json:"id"`
	SenderID  string `json:"sender_id"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

func toConv(c *domain.Conversation) conversationDTO {
	return conversationDTO{
		ID: c.ID, IsGroup: c.IsGroup, Title: c.Title, Participants: c.ParticipantIDs,
		UpdatedAt: c.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func (h *Handler) writeErr(w http.ResponseWriter, err error) {
	var ve application.ValidationError
	switch {
	case errors.As(err, &ve):
		common.WriteValidationError(w, ve.Msg)
	case errors.Is(err, domain.ErrNotParticipant):
		common.WriteForbiddenError(w, "you are not part of this conversation")
	case errors.Is(err, domain.ErrNotFound):
		common.WriteNotFoundError(w, "conversation not found")
	default:
		common.WriteInternalError(w, "something went wrong")
	}
}

type startRequest struct {
	ParticipantIDs []string `json:"participant_ids"`
	Title          string   `json:"title"`
}

func (h *Handler) Start(w http.ResponseWriter, r *http.Request) {
	var req startRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	c, err := h.svc.Start(r.Context(), common.UserIDFromContext(r.Context()), req.ParticipantIDs, req.Title)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusCreated, toConv(c))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	convs, err := h.svc.ListConversations(r.Context(), common.UserIDFromContext(r.Context()))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	out := make([]conversationDTO, 0, len(convs))
	for i := range convs {
		out = append(out, toConv(&convs[i]))
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"conversations": out})
}

func (h *Handler) Messages(w http.ResponseWriter, r *http.Request) {
	msgs, err := h.svc.ListMessages(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	out := make([]messageDTO, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, messageDTO{ID: m.ID, SenderID: m.SenderID, Body: m.Body, CreatedAt: m.CreatedAt.UTC().Format(time.RFC3339)})
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"messages": out})
}

type sendRequest struct {
	Body string `json:"body"`
}

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
	var req sendRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	m, err := h.svc.Send(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"), req.Body)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusCreated, messageDTO{
		ID: m.ID, SenderID: m.SenderID, Body: m.Body, CreatedAt: m.CreatedAt.UTC().Format(time.RFC3339),
	})
}

func (h *Handler) MarkRead(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.MarkRead(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id")); err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"read": true})
}

// Typing handles POST /conversations/{id}/typing — broadcasts an ephemeral
// typing indicator to the other participants.
func (h *Handler) Typing(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Typing(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id")); err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"ok": true})
}
