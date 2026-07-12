package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"workspace-app/internal/common"
	"workspace-app/internal/community/application"
	"workspace-app/internal/community/domain"
)

type Handler struct{ svc *application.Service }

func NewHandler(svc *application.Service) *Handler { return &Handler{svc: svc} }

func decode(r *http.Request, dst any) bool {
	return json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20)).Decode(dst) == nil
}

func (h *Handler) writeErr(w http.ResponseWriter, err error) {
	var ve application.ValidationError
	switch {
	case errors.As(err, &ve):
		common.WriteValidationError(w, ve.Msg)
	case errors.Is(err, domain.ErrInvalidPoll), errors.Is(err, domain.ErrOptionNotInPoll):
		common.WriteValidationError(w, err.Error())
	case errors.Is(err, domain.ErrPollExists):
		common.WriteError(w, common.NewConflictError("this post already has a poll"))
	case errors.Is(err, domain.ErrNotModerator):
		common.WriteForbiddenError(w, "you are not a moderator of this community")
	case errors.Is(err, domain.ErrNotFound), errors.Is(err, domain.ErrPollNotFound):
		common.WriteNotFoundError(w, "not found")
	default:
		common.WriteInternalError(w, "something went wrong")
	}
}

type communityDTO struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"`
	MemberCount int    `json:"member_count"`
}

func toCommunity(c *domain.Community) communityDTO {
	return communityDTO{ID: c.ID, Slug: c.Slug, Name: c.Name, Description: c.Description, Category: c.Category, MemberCount: c.MemberCount}
}

type postDTO struct {
	ID            string `json:"id"`
	AuthorID      string `json:"author_id"`
	Title         string `json:"title"`
	Body          string `json:"body,omitempty"`
	CommentCount  int    `json:"comment_count"`
	ReactionCount int    `json:"reaction_count"`
	CreatedAt     string `json:"created_at"`
}

func toPost(p *domain.Post) postDTO {
	return postDTO{
		ID: p.ID, AuthorID: p.AuthorID, Title: p.Title, Body: p.Body,
		CommentCount: p.CommentCount, ReactionCount: p.ReactionCount,
		CreatedAt: p.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	communities, err := h.svc.ListCommunities(r.Context())
	if err != nil {
		h.writeErr(w, err)
		return
	}
	out := make([]communityDTO, 0, len(communities))
	for i := range communities {
		out = append(out, toCommunity(&communities[i]))
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"communities": out})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.GetCommunity(r.Context(), r.PathValue("slug"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toCommunity(c))
}

func (h *Handler) ToggleJoin(w http.ResponseWriter, r *http.Request) {
	joined, err := h.svc.ToggleJoin(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("slug"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"joined": joined})
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	ctx := r.Context()
	var (
		c     *domain.Community
		posts []domain.Post
		err   error
	)
	if tag := r.URL.Query().Get("tag"); tag != "" {
		c, posts, err = h.svc.PostsByTag(ctx, slug, tag)
	} else {
		c, posts, err = h.svc.ListPosts(ctx, slug)
	}
	if err != nil {
		h.writeErr(w, err)
		return
	}
	out := make([]postDTO, 0, len(posts))
	for i := range posts {
		out = append(out, toPost(&posts[i]))
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"community": toCommunity(c), "posts": out})
}

type createPostRequest struct {
	Title string   `json:"title"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req createPostRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	p, err := h.svc.CreatePost(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("slug"), req.Title, req.Body, req.Tags)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusCreated, toPost(p))
}

func (h *Handler) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.svc.Tags(r.Context(), r.PathValue("slug"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	type tagDTO struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	out := make([]tagDTO, 0, len(tags))
	for _, t := range tags {
		out = append(out, tagDTO{Name: t.Name, Count: t.Count})
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"tags": out})
}

// --- polls ---

type pollDTO struct {
	ID       string          `json:"id"`
	PostID   string          `json:"post_id"`
	Question string          `json:"question"`
	Options  []pollOptionDTO `json:"options"`
}

type pollOptionDTO struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	VoteCount int    `json:"vote_count"`
}

func toPoll(p *domain.Poll) pollDTO {
	opts := make([]pollOptionDTO, 0, len(p.Options))
	for _, o := range p.Options {
		opts = append(opts, pollOptionDTO{ID: o.ID, Label: o.Label, VoteCount: o.VoteCount})
	}
	return pollDTO{ID: p.ID, PostID: p.PostID, Question: p.Question, Options: opts}
}

type createPollRequest struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

func (h *Handler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	var req createPollRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	p, err := h.svc.CreatePoll(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"), req.Question, req.Options)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusCreated, toPoll(p))
}

func (h *Handler) GetPoll(w http.ResponseWriter, r *http.Request) {
	p, err := h.svc.GetPoll(r.Context(), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toPoll(p))
}

type voteRequest struct {
	OptionID string `json:"option_id"`
}

func (h *Handler) Vote(w http.ResponseWriter, r *http.Request) {
	var req voteRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	p, err := h.svc.Vote(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"), req.OptionID)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toPoll(p))
}

// --- moderation / reporting ---

type reportRequest struct {
	Reason string `json:"reason"`
}

func (h *Handler) ReportPost(w http.ResponseWriter, r *http.Request) {
	var req reportRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	rep, err := h.svc.ReportPost(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"), req.Reason)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusCreated, map[string]any{"id": rep.ID, "status": rep.Status})
}

type reportDTO struct {
	ID         string `json:"id"`
	ReporterID string `json:"reporter_id"`
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
	Reason     string `json:"reason"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

func (h *Handler) ListReports(w http.ResponseWriter, r *http.Request) {
	reports, err := h.svc.CommunityReports(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("slug"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	out := make([]reportDTO, 0, len(reports))
	for _, rep := range reports {
		out = append(out, reportDTO{
			ID: rep.ID, ReporterID: rep.ReporterID, TargetType: rep.TargetType, TargetID: rep.TargetID,
			Reason: rep.Reason, Status: rep.Status, CreatedAt: rep.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"reports": out})
}

func (h *Handler) HidePost(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.HidePost(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("slug"), r.PathValue("id")); err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"hidden": true})
}

type commentRequest struct {
	Body string `json:"body"`
}

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	var req commentRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request payload")
		return
	}
	c, err := h.svc.AddComment(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"), req.Body)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusCreated, map[string]any{
		"id": c.ID, "body": c.Body, "created_at": c.CreatedAt.UTC().Format(time.RFC3339),
	})
}

func (h *Handler) ListComments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.svc.Comments(r.Context(), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	type dto struct {
		ID        string `json:"id"`
		AuthorID  string `json:"author_id"`
		Body      string `json:"body"`
		CreatedAt string `json:"created_at"`
	}
	out := make([]dto, 0, len(comments))
	for _, c := range comments {
		out = append(out, dto{ID: c.ID, AuthorID: c.AuthorID, Body: c.Body, CreatedAt: c.CreatedAt.UTC().Format(time.RFC3339)})
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"comments": out})
}

func (h *Handler) ToggleReaction(w http.ResponseWriter, r *http.Request) {
	reacted, err := h.svc.ToggleReaction(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]bool{"reacted": reacted})
}

type createCommunityRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createCommunityRequest
	if !decode(r, &req) {
		common.WriteValidationError(w, "invalid request body")
		return
	}

	userID := common.UserIDFromContext(r.Context())
	if userID == "" {
		common.WriteUnauthorizedError(w, "unauthorized")
		return
	}

	c, err := h.svc.CreateCommunity(r.Context(), userID, req.Name, req.Slug, req.Description, req.Category)
	if err != nil {
		h.writeErr(w, err)
		return
	}

	common.WriteSuccess(w, http.StatusCreated, toCommunity(c))
}
