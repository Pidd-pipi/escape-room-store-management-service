package repository

import (
	"context"

	"github.com/lp/escape-room-ops/internal/model"
	"gorm.io/gorm"
)

// EscapeRecordRepository persists escape outcome rows.
type EscapeRecordRepository struct {
	db *gorm.DB
}

// NewEscapeRecordRepository builds an EscapeRecordRepository.
func NewEscapeRecordRepository(db *gorm.DB) *EscapeRecordRepository {
	return &EscapeRecordRepository{db: db}
}

// Create inserts a new escape record.
func (r *EscapeRecordRepository) Create(ctx context.Context, e *model.EscapeRecord) error {
	return db(ctx, r.db).Create(e).Error
}

// ListAll returns all records for leaderboard aggregation.
func (r *EscapeRecordRepository) ListAll(ctx context.Context) ([]model.EscapeRecord, error) {
	var items []model.EscapeRecord
	err := db(ctx, r.db).Order("duration_minutes ASC").Find(&items).Error
	return items, err
}

// ListByUser returns records of the themes the user played (via registrations handled in service).
func (r *EscapeRecordRepository) ListByTheme(ctx context.Context, themeRoomID uint) ([]model.EscapeRecord, error) {
	var items []model.EscapeRecord
	err := db(ctx, r.db).Where("theme_room_id = ?", themeRoomID).Order("duration_minutes ASC").Find(&items).Error
	return items, err
}
