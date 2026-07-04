package api

import (
	"errors"
	"io"
	"net/http"

	"workspace-app/internal/common"
	"workspace-app/internal/resume/application"
	"workspace-app/internal/resume/domain"
)

const maxUpload = 10 << 20 // 10 MiB

type Handler struct {
	svc *application.Service
}

func NewHandler(svc *application.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) writeErr(w http.ResponseWriter, err error) {
	var ve application.ValidationError
	switch {
	case errors.As(err, &ve):
		common.WriteValidationError(w, ve.Msg)
	case errors.Is(err, application.ErrForbidden):
		common.WriteForbiddenError(w, "you do not own this resume")
	case errors.Is(err, domain.ErrNotFound):
		common.WriteNotFoundError(w, "resume not found")
	default:
		common.WriteInternalError(w, "something went wrong")
	}
}

// readUpload extracts the multipart file + title.
func readUpload(w http.ResponseWriter, r *http.Request) (application.UploadInput, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUpload+(1<<20))
	if err := r.ParseMultipartForm(maxUpload); err != nil {
		common.WriteValidationError(w, "file exceeds the 10MB limit or form is malformed")
		return application.UploadInput{}, false
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		common.WriteValidationError(w, "a 'file' field is required")
		return application.UploadInput{}, false
	}
	defer func() { _ = file.Close() }()
	data, err := io.ReadAll(file)
	if err != nil {
		common.WriteValidationError(w, "could not read uploaded file")
		return application.UploadInput{}, false
	}
	return application.UploadInput{
		Title:       r.FormValue("title"),
		Filename:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		Data:        data,
	}, true
}

// Upload handles POST /resumes (multipart).
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	in, ok := readUpload(w, r)
	if !ok {
		return
	}
	res, err := h.svc.Upload(r.Context(), common.UserIDFromContext(r.Context()), in)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusCreated, toResume(res))
}

// AddVersion handles POST /resumes/{id}/versions (multipart).
func (h *Handler) AddVersion(w http.ResponseWriter, r *http.Request) {
	in, ok := readUpload(w, r)
	if !ok {
		return
	}
	res, err := h.svc.AddVersion(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"), in)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusCreated, toResume(res))
}

// List handles GET /resumes.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.List(r.Context(), common.UserIDFromContext(r.Context()))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"resumes": toResumes(list)})
}

// Get handles GET /resumes/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	res, err := h.svc.Get(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toResume(res))
}

// Versions handles GET /resumes/{id}/versions.
func (h *Handler) Versions(w http.ResponseWriter, r *http.Request) {
	vs, err := h.svc.Versions(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	out := make([]versionResponse, 0, len(vs))
	for i := range vs {
		out = append(out, toVersion(&vs[i]))
	}
	common.WriteSuccess(w, http.StatusOK, map[string]any{"versions": out})
}

// Score handles GET /resumes/{id}/score.
func (h *Handler) Score(w http.ResponseWriter, r *http.Request) {
	s, err := h.svc.Score(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toScore(s))
}

// Review handles POST /resumes/{id}/review.
func (h *Handler) Review(w http.ResponseWriter, r *http.Request) {
	s, err := h.svc.Review(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, toScore(s))
}

// Delete handles DELETE /resumes/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.Context(), common.UserIDFromContext(r.Context()), r.PathValue("id")); err != nil {
		h.writeErr(w, err)
		return
	}
	common.WriteSuccess(w, http.StatusOK, map[string]string{"status": "deleted"})
}
