package application

import (
	"context"
	"errors"
	"testing"

	"workspace-app/internal/resume/domain"
)

// errParser always reports the given error from ExtractText.
type errParser struct{ err error }

func (p errParser) ExtractText(_, _ string, _ []byte) (string, error) { return "", p.err }

func TestUploadUnsupportedTypeRejected(t *testing.T) {
	svc := NewService(newFakeRepo(), newFakeStorage(), errParser{err: domain.ErrUnsupported}, nil)
	var ve ValidationError
	_, err := svc.Upload(context.Background(), "u1", UploadInput{Filename: "cv.png", ContentType: "image/png", Data: []byte("x")})
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError for unsupported type, got %v", err)
	}
}

func TestListAttachesLatestScore(t *testing.T) {
	svc := newSvc()
	ctx := context.Background()
	if _, err := svc.Upload(ctx, "u1", UploadInput{Filename: "a.txt", ContentType: "text/plain", Data: []byte("a")}); err != nil {
		t.Fatalf("upload a: %v", err)
	}
	if _, err := svc.Upload(ctx, "u1", UploadInput{Filename: "b.txt", ContentType: "text/plain", Data: []byte("b")}); err != nil {
		t.Fatalf("upload b: %v", err)
	}
	// A resume belonging to someone else must not appear.
	if _, err := svc.Upload(ctx, "u2", UploadInput{Filename: "c.txt", ContentType: "text/plain", Data: []byte("c")}); err != nil {
		t.Fatalf("upload c: %v", err)
	}

	list, err := svc.List(ctx, "u1")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 resumes for u1, got %d", len(list))
	}
	for _, r := range list {
		if r.Score == nil {
			t.Fatalf("expected a latest score attached to resume %s", r.ID)
		}
	}
}

func TestVersionsOwnership(t *testing.T) {
	svc := newSvc()
	ctx := context.Background()
	res, _ := svc.Upload(ctx, "owner", UploadInput{Filename: "cv.txt", ContentType: "text/plain", Data: []byte("v1")})

	if _, err := svc.Versions(ctx, "intruder", res.ID); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
	vs, err := svc.Versions(ctx, "owner", res.ID)
	if err != nil {
		t.Fatalf("versions: %v", err)
	}
	if len(vs) != 1 {
		t.Fatalf("expected 1 version, got %d", len(vs))
	}
}

func TestReviewRescoresLatestVersion(t *testing.T) {
	svc := newSvc()
	ctx := context.Background()
	res, _ := svc.Upload(ctx, "owner", UploadInput{Filename: "cv.txt", ContentType: "text/plain", Data: []byte("v1")})

	if _, err := svc.Review(ctx, "intruder", res.ID); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
	score, err := svc.Review(ctx, "owner", res.ID)
	if err != nil {
		t.Fatalf("review: %v", err)
	}
	if score == nil || score.VersionID == "" {
		t.Fatalf("expected a score bound to the latest version, got %+v", score)
	}
}

func TestDeleteOwnership(t *testing.T) {
	svc := newSvc()
	ctx := context.Background()
	res, _ := svc.Upload(ctx, "owner", UploadInput{Filename: "cv.txt", ContentType: "text/plain", Data: []byte("v1")})

	if err := svc.Delete(ctx, "intruder", res.ID); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
	if err := svc.Delete(ctx, "owner", res.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := svc.Get(ctx, "owner", res.ID); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}
