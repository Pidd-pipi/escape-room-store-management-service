package repository

import (
	"context"

	"github.com/lp/escape-room-ops/internal/model"
	"gorm.io/gorm"
)

// ThemeRoomRepository persists theme room rows.
type ThemeRoomRepository struct {
	db *gorm.DB
}

// NewThemeRoomRepository builds a ThemeRoomRepository.
func NewThemeRoomRepository(db *gorm.DB) *ThemeRoomRepository {
	return &ThemeRoomRepository{db: db}
}

// Create inserts a new theme room.
func (r *ThemeRoomRepository) Create(ctx context.Context, t *model.ThemeRoom) error {
	return db(ctx, r.db).Create(t).Error
}

// FindByID returns a theme room by id.
func (r *ThemeRoomRepository) FindByID(ctx context.Context, id uint) (*model.ThemeRoom, error) {
	var t model.ThemeRoom
	err := db(ctx, r.db).First(&t, id).Error
	if err != nil {
		return nil, normalizeError(err)
	}
	return &t, nil
}

// List filters theme rooms by category/stars with pagination.
func (r *ThemeRoomRepository) List(ctx context.Context, category string, stars, page, pageSize int) ([]model.ThemeRoom, int64, error) {
	q := db(ctx, r.db).Model(&model.ThemeRoom{})
	if category != "" {
		q = q.Where("category = ?", category)
	}
	if stars > 0 {
		q = q.Where("difficulty_stars = ?", stars)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.ThemeRoom
	err := q.Order("created_at DESC").Offset(page * pageSize).Limit(pageSize).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// Count returns the total theme count.
func (r *ThemeRoomRepository) Count(ctx context.Context) (int64, error) {
	var n int64
	err := db(ctx, r.db).Model(&model.ThemeRoom{}).Count(&n).Error
	return n, err
}
