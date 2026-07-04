package application

import (
	"context"
	"errors"
	"testing"
)

func newSvc() *Service {
	return NewService(newFakeRepo(), newFakeStorage(), fakeParser{text: "Experience Education Skills jane@x.com led improved 30%"}, nil)
}

func TestUploadCreatesResumeWithScore(t *testing.T) {
	svc := newSvc()
	ctx := context.Background()

	res, err := svc.Upload(ctx, "user-1", UploadInput{Filename: "cv.txt", ContentType: "text/plain", Data: []byte("content")})
	if err != nil {
		t.Fatalf("upload: %v", err)
	}
	if res.ID == "" || res.Latest == nil || res.Latest.VersionNo != 1 {
		t.Fatalf("expected version 1 attached, got %+v", res)
	}
	if res.Score == nil {
		t.Fatal("expected a score attached to the uploaded resume")
	}
	if len(res.Versions) != 1 {
		t.Fatalf("expected 1 version, got %d", len(res.Versions))
	}
}

func TestUploadEmptyRejected(t *testing.T) {
	svc := newSvc()
	var ve ValidationError
	if _, err := svc.Upload(context.Background(), "user-1", UploadInput{Filename: "cv.txt", Data: nil}); !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError for empty upload, got %v", err)
	}
}

func TestAddVersionIncrementsAndOwnership(t *testing.T) {
	svc := newSvc()
	ctx := context.Background()
	res, _ := svc.Upload(ctx, "owner", UploadInput{Filename: "cv.txt", ContentType: "text/plain", Data: []byte("v1")})

	// Non-owner cannot add a version.
	if _, err := svc.AddVersion(ctx, "intruder", res.ID, UploadInput{Filename: "cv.txt", Data: []byte("v2")}); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}

	updated, err := svc.AddVersion(ctx, "owner", res.ID, UploadInput{Filename: "cv.txt", ContentType: "text/plain", Data: []byte("v2")})
	if err != nil {
		t.Fatalf("add version: %v", err)
	}
	if updated.Latest.VersionNo != 2 || len(updated.Versions) != 2 {
		t.Fatalf("expected version 2, got latest=%v count=%d", updated.Latest.VersionNo, len(updated.Versions))
	}
}

func TestGetAndScoreOwnership(t *testing.T) {
	svc := newSvc()
	ctx := context.Background()
	res, _ := svc.Upload(ctx, "owner", UploadInput{Filename: "cv.txt", ContentType: "text/plain", Data: []byte("v1")})

	if _, err := svc.Get(ctx, "intruder", res.ID); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden on get, got %v", err)
	}
	score, err := svc.Score(ctx, "owner", res.ID)
	if err != nil {
		t.Fatalf("score: %v", err)
	}
	if score == nil {
		t.Fatal("expected a score")
	}
}
