package service

import (
	"context"
	"log/slog"
	"testing"

	"github.com/lp/escape-room-ops/internal/constants"
	"github.com/lp/escape-room-ops/internal/dto"
	"github.com/lp/escape-room-ops/internal/model"
	"github.com/lp/escape-room-ops/internal/util"
)

type fakeThemeRepo struct {
	themes map[uint]*model.ThemeRoom
	nextID uint
}

func newFakeThemeRepo() *fakeThemeRepo {
	return &fakeThemeRepo{themes: map[uint]*model.ThemeRoom{}, nextID: 1}
}

func (f *fakeThemeRepo) Create(_ context.Context, t *model.ThemeRoom) error {
	t.ID = f.nextID
	f.nextID++
	f.themes[t.ID] = t
	return nil
}

func (f *fakeThemeRepo) FindByID(_ context.Context, id uint) (*model.ThemeRoom, error) {
	if t, ok := f.themes[id]; ok {
		cp := *t
		return &cp, nil
	}
	return nil, util.ErrNotFound
}

func (f *fakeThemeRepo) List(_ context.Context, category string, stars, page, pageSize int) ([]model.ThemeRoom, int64, error) {
	var out []model.ThemeRoom
	for _, t := range f.themes {
		if category != "" && t.Category != category {
			continue
		}
		out = append(out, *t)
	}
	return out, int64(len(out)), nil
}

func (f *fakeThemeRepo) Count(context.Context) (int64, error) { return int64(len(f.themes)), nil }

func TestThemeRoomServiceCreate(t *testing.T) {
	svc := NewThemeRoomService(newFakeThemeRepo(), slog.Default())
	tests := []struct {
		name     string
		category string
		wantErr  bool
	}{
		{name: "valid horror", category: constants.ThemeCategoryHorror, wantErr: false},
		{name: "valid scifi", category: constants.ThemeCategoryScifi, wantErr: false},
		{name: "invalid category", category: "fantasy", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &dto.CreateThemeRoomRequest{Name: "测试主题", Category: tt.category, DifficultyStars: 3, MinPlayers: 2, MaxPlayers: 6, DurationMinutes: 60}
			_, err := svc.Create(context.Background(), req)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestThemeRoomServiceGetNotFound(t *testing.T) {
	svc := NewThemeRoomService(newFakeThemeRepo(), slog.Default())
	if _, err := svc.Get(context.Background(), 999); err == nil {
		t.Fatalf("expected not found error")
	}
}
