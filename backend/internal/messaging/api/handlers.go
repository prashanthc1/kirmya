package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"

	"workspace-app/internal/common"
	"workspace-app/internal/messaging/application"
	"workspace-app/internal/messaging/domain"
)

type Handler struct{ svc *application.Service }

func NewHandler(svc *application.Service) *Handler { return &Handler{svc: svc} }

func decode(r *http.Request, dst any) bool {
	return json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20)).Decode(dst) == nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Rely on JWT auth and standard CORS configurations
	},
}

type wsMessage struct {
	Type           string `json:"type"`            // "ping" | "typing"
	ConversationID string `json:"conversation_id"` // for typing indicators
}

// authenticateWS parses and verifies the JWT token from the query string
func authenticateWS(tokenStr string) (string, string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", "", errors.New("JWT_SECRET is not configured")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid claims")
	}

	sub, _ := claims["sub"].(string)
	email, _ := claims["email"].(string)
	if sub == "" {
		return "", "", errors.New("subject claim missing")
	}

	return sub, email, nil
}

// WebSocket handles GET /api/v1/ws. Upgrades connection and routes real-time events.
func (h *Handler) WebSocket(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		common.WriteUnauthorizedError(w, "missing token query parameter")
		return
	}

	userID, _, err := authenticateWS(tokenStr)
	if err != nil {
		common.WriteUnauthorizedError(w, "invalid or expired token")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ws] upgrade failed for user %s: %v", userID, err)
		return
	}
	defer conn.Close()

	log.Printf("[ws] user %s connected", userID)

	// Set presence to online
	_ = h.svc.UpdateUserPresence(context.Background(), userID, true)
	defer func() {
		_ = h.svc.UpdateUserPresence(context.Background(), userID, false)
		log.Printf("[ws] user %s disconnected", userID)
	}()

	// Subscribe to hub messages
	ch, cancel := h.svc.Subscribe(userID)
	defer cancel()

	// 1. Writer loop: Send events from hub to client
	go func() {
		for ev := range ch {
			dto := toStreamEvent(ev)
			if err := conn.WriteJSON(dto); err != nil {
				return
			}
		}
	}()

	// Keep alive ticker
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if err := conn.WriteJSON(map[string]string{"type": "ping"}); err != nil {
				return
			}
		}
	}()

	// 2. Reader loop: Read messages from client
	for {
		var msg wsMessage
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}

		switch msg.Type {
		case "ping":
			_ = conn.WriteJSON(map[string]string{"type": "pong"})
			_ = h.svc.UpdateUserPresence(context.Background(), userID, true)
		case "typing":
			if msg.ConversationID != "" {
				_ = h.svc.Typing(context.Background(), userID, msg.ConversationID)
			}
		}
	}
}

// streamEventDTO is a multiplexed conversation event.
type streamEventDTO struct {
	Kind           string `json:"kind"`
	ConversationID string `json:"conversation_id"`
	SenderID       string `json:"sender_id,omitempty"`
	ReaderID       string `json:"reader_id,omitempty"`
	ID             string `json:"id,omitempty"`
	Body           string `json:"body,omitempty"`
	ContentType    string `json:"content_type,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	At             string `json:"at,omitempty"`
}

func toStreamEvent(ev application.StreamEvent) streamEventDTO {
	d := streamEventDTO{Kind: ev.Kind, ConversationID: ev.ConversationID}
	switch ev.Kind {
	case application.EventMessage:
		d.SenderID = ev.Message.SenderID
		d.ID = ev.Message.ID
		d.Body = ev.Message.Content
		d.ContentType = ev.Message.ContentType
		d.CreatedAt = ev.Message.CreatedAt.UTC().Format(time.RFC3339)
	case application.EventTyping:
		d.SenderID = ev.ActorID
	case application.EventRead:
		d.ReaderID = ev.ActorID
		d.At = ev.At.UTC().Format(time.RFC3339)
	}
	return d
}

// Stream handles GET /conversations/stream (fallback SSE)
func (h *Handler) Stream(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	rc := http.NewResponseController(w)
	_ = rc.SetWriteDeadline(time.Time{})

	ch, cancel := h.svc.Subscribe(uid)
	defer cancel()

	_, _ = fmt.Fprint(w, ": connected\n\n")
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
			_, _ = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", ev.Kind, data)
			if rc.Flush() != nil {
				return
			}
		case <-ping.C:
			_, _ = fmt.Fprint(w, ": ping\n\n")
			if rc.Flush() != nil {
				return
			}
		}
	}
}

type conversationDTO struct {
	ID                 string   `json:"id"`
	Type               string   `json:"type"`
	Title              string   `json:"title,omitempty"`
	Participants       []string `json:"participants"`
	UpdatedAt          string   `json:"updated_at"`
	UnreadCount        int      `json:"unread_count"`
	LastMessagePreview string   `json:"last_message_preview,omitempty"`
	IsPinned           bool     `json:"is_pinned"`
	IsArchived         bool     `json:"is_archived"`
}

type messageDTO struct {
	ID          string `json:"id"`
	SenderID    string `json:"sender_id"`
	Body        string `json:"body"`
	ContentType string `json:"content_type"`
	CreatedAt   string `json:"created_at"`
	EditedAt    string `json:"edited_at,omitempty"`
	DeletedAt   string `json:"deleted_at,omitempty"`
}

func toConv(c *domain.Conversation) conversationDTO {
	return conversationDTO{
		ID:                 c.ID,
		Type:               c.Type,
		Title:              c.Title,
		Participants:       c.ParticipantIDs,
		UpdatedAt:          c.UpdatedAt.UTC().Format(time.RFC3339),
		UnreadCount:        c.UnreadCount,
		LastMessagePreview: c.LastMessagePreview,
		IsPinned:           c.IsPinned,
		IsArchived:         c.IsArchived,
	}
}

func (h *Handler) writeErr(w http.ResponseWriter, err error) {
	var ve application.ValidationError
	switch {
	case errors.As(err, &ve):
		if strings.Contains(strings.ToLower(ve.Msg), "forbidden") {
			common.WriteForbiddenError(w, ve.Msg)
		} else {
			common.WriteValidationError(w, ve.Msg)
		}
	case errors.Is(err, domain.ErrNotParticipant):
		common.WriteForbiddenError(w, "you are not part of this conversation")
	case errors.Is(err, domain.ErrNotFound):
		common.WriteNotFoundError(w, "conversation not found")
	default:
		common.WriteInternalError(w, err.Error())
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
	uid := common.UserIDFromContext(r.Context())
	convID := r.PathValue("id")
	q := r.URL.Query().Get("q")

	var msgs []domain.Message
	var err error

	if q != "" {
		msgs, err = h.svc.SearchMessages(r.Context(), uid, convID, q)
	} else {
		msgs, err = h.svc.ListMessages(r.Context(), uid, convID)
	}

	if err != nil {
		h.writeErr(w, err)
		return
	}

	out := make([]messageDTO, 0, len(msgs))
	for _, m := range msgs {
		dto := messageDTO{
			ID:          m.ID,
			SenderID:    m.SenderID,
			Body:        m.Content,
			ContentType: m.ContentType,
			CreatedAt:   m.CreatedAt.UTC().Format(time.RFC3339),
		}
		if m.EditedAt != nil {
			dto.EditedAt = m.EditedAt.UTC().Format(time.RFC3339)
		}
		if m.DeletedAt != nil {
			dto.DeletedAt = m.DeletedAt.UTC().Format(time.RFC3339)
			dto.Body = "This message was deleted."
		}
		out = append(out, dto)
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"messages": out})
}

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	convID := r.PathValue("id")

	var content string
	var contentType string = "text"
	var fileData []byte
	var fileName string

	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		// Multipart file upload
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			common.WriteValidationError(w, "failed to parse multipart form")
			return
		}
		content = r.FormValue("body")
		contentType = r.FormValue("content_type")
		if contentType == "" {
			contentType = "text"
		}

		file, header, err := r.FormFile("attachment")
		if err == nil {
			defer file.Close()
			fileName = header.Filename
			fileData, err = io.ReadAll(file)
			if err != nil {
				common.WriteValidationError(w, "failed to read uploaded file")
				return
			}
		}
	} else {
		// Standard JSON body
		var req struct {
			Body        string `json:"body"`
			ContentType string `json:"content_type"`
		}
		if !decode(r, &req) {
			common.WriteValidationError(w, "invalid request payload")
			return
		}
		content = req.Body
		contentType = req.ContentType
		if contentType == "" {
			contentType = "text"
		}
	}

	m, err := h.svc.Send(r.Context(), uid, convID, content, contentType, fileData, fileName)
	if err != nil {
		h.writeErr(w, err)
		return
	}

	dto := messageDTO{
		ID:          m.ID,
		SenderID:    m.SenderID,
		Body:        m.Content,
		ContentType: m.ContentType,
		CreatedAt:   m.CreatedAt.UTC().Format(time.RFC3339),
	}
	common.WriteSuccess(w, http.StatusCreated, dto)
}

func (h *Handler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	messageID := r.PathValue("messageID")

	if err := h.svc.DeleteMessage(r.Context(), uid, messageID); err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"deleted": true})
}

func (h *Handler) MarkRead(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.MarkRead(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id")); err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"read": true})
}

func (h *Handler) Typing(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Typing(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id")); err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) Archive(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	convID := r.PathValue("id")
	var req struct {
		Archive bool `json:"archive"`
	}
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}

	if err := h.svc.ArchiveConversation(r.Context(), uid, convID, req.Archive); err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"archived": req.Archive})
}

func (h *Handler) Pin(w http.ResponseWriter, r *http.Request) {
	uid := common.UserIDFromContext(r.Context())
	convID := r.PathValue("id")
	var req struct {
		Pin bool `json:"pin"`
	}
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}

	if err := h.svc.PinConversation(r.Context(), uid, convID, req.Pin); err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"pinned": req.Pin})
}
