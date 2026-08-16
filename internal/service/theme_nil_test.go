package service

import (
	"context"
	"log/slog"
	"testing"
)

func TestGetMissingThemeReturnsError(t *testing.T) {
	svc := NewThemeRoomService(newFakeThemeRepo(), slog.Default())
	theme, err := svc.Get(context.Background(), 999)
	if err == nil {
		t.Fatalf("expected error for missing theme, got theme=%+v", theme)
	}
	if theme != nil {
		t.Fatalf("theme should be nil on error, got %+v", theme)
	}
}
