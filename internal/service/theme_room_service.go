package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lp/escape-room-ops/internal/constants"
	"github.com/lp/escape-room-ops/internal/dto"
	"github.com/lp/escape-room-ops/internal/model"
	"github.com/lp/escape-room-ops/internal/util"
)

// ThemeRoomRepository is the data access contract for theme room rows.
type ThemeRoomRepository interface {
	Create(ctx context.Context, t *model.ThemeRoom) error
	FindByID(ctx context.Context, id uint) (*model.ThemeRoom, error)
	List(ctx context.Context, category string, stars, page, pageSize int) ([]model.ThemeRoom, int64, error)
	Count(ctx context.Context) (int64, error)
}

// ThemeRoomService manages escape room themes.
type ThemeRoomService struct {
	themes ThemeRoomRepository
	logger *slog.Logger
}

// NewThemeRoomService wires the theme room service dependencies.
func NewThemeRoomService(themes ThemeRoomRepository, logger *slog.Logger) *ThemeRoomService {
	return &ThemeRoomService{themes: themes, logger: logger}
}

// Create adds a new theme room.
func (s *ThemeRoomService) Create(ctx context.Context, req *dto.CreateThemeRoomRequest) (*model.ThemeRoom, error) {
	t := &model.ThemeRoom{
		Name: req.Name, Category: req.Category, DifficultyStars: req.DifficultyStars,
		MinPlayers: req.MinPlayers, MaxPlayers: req.MaxPlayers, DurationMinutes: req.DurationMinutes,
		Story: req.Story, PosterURL: req.PosterURL, SceneImages: req.SceneImages, Status: "active",
	}
	if err := s.themes.Create(ctx, t); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogThemeCreateFailed, req.Name, err))
		return nil, util.WrapAppError(fmt.Errorf("theme_room[name=%s] create: %w", req.Name, err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	s.logger.Info(fmt.Sprintf(constants.LogThemeCreateSuccess, t.ID, t.Name, t.Category))
	return t, nil
}

// Get returns one theme room.
func (s *ThemeRoomService) Get(ctx context.Context, id uint) (*model.ThemeRoom, error) {
	t, err := s.themes.FindByID(ctx, id)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("theme_room[id=%d] get: %w", id, err), 404, constants.CodeNotFound, constants.MsgNotFound)
	}
	return t, nil
}

// List filters theme rooms.
func (s *ThemeRoomService) List(ctx context.Context, q *dto.ListThemeQuery) (*dto.PageResult, error) {
	q.Normalize()
	items, total, err := s.themes.List(ctx, q.Category, q.Stars, q.Page, q.PageSize)
	if err != nil {
		return nil, util.WrapAppError(fmt.Errorf("theme_room list: %w", err), 500, constants.CodeInternalError, constants.MsgInternalError)
	}
	return &dto.PageResult{Items: items, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}
