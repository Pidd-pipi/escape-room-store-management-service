package repository

import (
	"context"

	"github.com/lp/escape-room-ops/internal/model"
	"github.com/lp/escape-room-ops/internal/util"
	"gorm.io/gorm"
)

// PropRepository persists prop and maintenance rows.
type PropRepository struct {
	db *gorm.DB
}

// NewPropRepository builds a PropRepository.
func NewPropRepository(db *gorm.DB) *PropRepository {
	return &PropRepository{db: db}
}

// Transaction runs fn inside a database transaction for multi-write prop flows.
func (r *PropRepository) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return Transaction(ctx, r.db, fn)
}

// Create inserts a new prop.
func (r *PropRepository) Create(ctx context.Context, p *model.Prop) error {
	return db(ctx, r.db).Create(p).Error
}

// FindByID returns a prop by id.
func (r *PropRepository) FindByID(ctx context.Context, id uint) (*model.Prop, error) {
	var p model.Prop
	err := db(ctx, r.db).First(&p, id).Error
	if err != nil {
		return nil, normalizeError(err)
	}
	return &p, nil
}

// List returns props filtered by theme with pagination.
func (r *PropRepository) List(ctx context.Context, themeRoomID uint, page, pageSize int) ([]model.Prop, int64, error) {
	q := db(ctx, r.db).Model(&model.Prop{})
	if themeRoomID > 0 {
		q = q.Where("theme_room_id = ?", themeRoomID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.Prop
	err := q.Order("created_at DESC").Offset(page * pageSize).Limit(pageSize).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// UpdateStock sets stock and auto status based on the threshold.
func (r *PropRepository) UpdateStock(ctx context.Context, id uint, stock int) error {
	return db(ctx, r.db).Model(&model.Prop{}).Where("id = ?", id).
		Updates(map[string]interface{}{"stock": stock}).Error
}

// RefreshStatus recomputes status from stock vs threshold.
func (r *PropRepository) RefreshStatus(ctx context.Context, id uint) error {
	return db(ctx, r.db).Model(&model.Prop{}).Where("id = ?", id).
		UpdateColumn("status", gorm.Expr("CASE WHEN stock = 0 THEN 'repairing' WHEN stock <= alert_threshold THEN 'low_stock' ELSE 'normal' END")).Error
}

// Consume decreases stock by quantity; returns ErrConflict when stock is insufficient.
func (r *PropRepository) Consume(ctx context.Context, id uint, quantity int) error {
	res := db(ctx, r.db).Model(&model.Prop{}).Where("id = ? AND stock >= ?", id, quantity).
		UpdateColumn("stock", gorm.Expr("stock - ?", quantity))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return util.ErrConflict
	}
	return nil
}

// CreateMaintenance inserts a maintenance record and sets the prop to repairing.
func (r *PropRepository) CreateMaintenance(ctx context.Context, m *model.PropMaintenance) error {
	return db(ctx, r.db).Create(m).Error
}

// MarkRepairing sets a prop status to repairing.
func (r *PropRepository) MarkRepairing(ctx context.Context, id uint) error {
	return db(ctx, r.db).Model(&model.Prop{}).Where("id = ?", id).Update("status", "repairing").Error
}

// ListMaintenance returns maintenance records of a prop.
func (r *PropRepository) ListMaintenance(ctx context.Context, propID uint) ([]model.PropMaintenance, error) {
	var items []model.PropMaintenance
	err := db(ctx, r.db).Where("prop_id = ?", propID).Order("created_at DESC").Find(&items).Error
	return items, err
}
