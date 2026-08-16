package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/lp/escape-room-ops/internal/constants"
	"github.com/lp/escape-room-ops/internal/dto"
	"github.com/lp/escape-room-ops/internal/util"
)

func TestThemeCategoryValidationBoundary(t *testing.T) {
	if constants.IsThemeCategory("fantasy") {
		t.Fatal("IsThemeCategory('fantasy') should be false")
	}
	if util.ThemeCategoryText("fantasy") != "未知" {
		t.Fatalf("ThemeCategoryText('fantasy') = %s, want 未知", util.ThemeCategoryText("fantasy"))
	}
	svc := NewThemeRoomService(newFakeThemeRepo(), slog.Default())
	req := &dto.CreateThemeRoomRequest{Name: "非法分类主题", Category: "fantasy", DifficultyStars: 3, MinPlayers: 2, MaxPlayers: 6, DurationMinutes: 60}
	if _, err := svc.Create(context.Background(), req); err == nil {
		t.Fatal("invalid theme category should be rejected")
	}
}
