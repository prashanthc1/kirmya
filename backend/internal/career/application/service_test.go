package application

import (
	"context"
	"testing"
)

func TestServicePaths(t *testing.T) {
	svc := NewService()
	p := svc.Paths(context.Background(), "Operations Manager")
	if p.From != "Operations Manager" {
		t.Fatalf("expected From echoed, got %q", p.From)
	}
	if p.Target != "VP of Operations" || len(p.Steps) == 0 {
		t.Fatalf("expected curated ladder, got %+v", p)
	}
}
