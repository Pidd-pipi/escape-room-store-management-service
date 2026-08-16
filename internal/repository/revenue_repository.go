package repository

import (
	"context"
	"time"

	"github.com/lp/escape-room-ops/internal/model"
	"gorm.io/gorm"
)

// RevenueRepository persists revenue rows and provides analytics queries.
type RevenueRepository struct {
	db *gorm.DB
}

// NewRevenueRepository builds a RevenueRepository.
func NewRevenueRepository(db *gorm.DB) *RevenueRepository {
	return &RevenueRepository{db: db}
}

// Create inserts a revenue record.
func (r *RevenueRepository) Create(ctx context.Context, rev *model.Revenue) error {
	return db(ctx, r.db).Create(rev).Error
}

// ListSince returns revenue records on or after a date.
func (r *RevenueRepository) ListSince(ctx context.Context, since time.Time) ([]model.Revenue, error) {
	var items []model.Revenue
	err := db(ctx, r.db).Where("date >= ?", since).Order("date ASC").Find(&items).Error
	return items, err
}

// ListAll returns all revenue records.
func (r *RevenueRepository) ListAll(ctx context.Context) ([]model.Revenue, error) {
	var items []model.Revenue
	err := db(ctx, r.db).Order("date DESC").Find(&items).Error
	return items, err
}

// Count returns the revenue record count.
func (r *RevenueRepository) Count(ctx context.Context) (int64, error) {
	var n int64
	err := db(ctx, r.db).Model(&model.Revenue{}).Count(&n).Error
	return n, err
}
