package repository

import (
	"context"

	"github.com/lp/escape-room-ops/internal/model"
	"github.com/lp/escape-room-ops/internal/util"
	"gorm.io/gorm"
)

// GameSessionRepository persists game session and registration rows.
type GameSessionRepository struct {
	db *gorm.DB
}

// NewGameSessionRepository builds a GameSessionRepository.
func NewGameSessionRepository(db *gorm.DB) *GameSessionRepository {
	return &GameSessionRepository{db: db}
}

// Transaction runs fn inside a database transaction for cross-repository writes.
func (r *GameSessionRepository) Transaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return Transaction(ctx, r.db, fn)
}

// Create inserts a new session.
func (r *GameSessionRepository) Create(ctx context.Context, s *model.GameSession) error {
	return db(ctx, r.db).Create(s).Error
}

// FindByID returns a session by id.
func (r *GameSessionRepository) FindByID(ctx context.Context, id uint) (*model.GameSession, error) {
	var s model.GameSession
	err := db(ctx, r.db).First(&s, id).Error
	if err != nil {
		return nil, normalizeError(err)
	}
	return &s, nil
}

// List filters sessions by theme/status with pagination.
func (r *GameSessionRepository) List(ctx context.Context, themeRoomID uint, status string, page, pageSize int) ([]model.GameSession, int64, error) {
	q := db(ctx, r.db).Model(&model.GameSession{})
	if themeRoomID > 0 {
		q = q.Where("theme_room_id = ?", themeRoomID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.GameSession
	err := q.Order("start_time ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// UpdateStatus sets the session status.
func (r *GameSessionRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	res := db(ctx, r.db).Model(&model.GameSession{}).Where("id = ?", id).Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return util.ErrNotFound
	}
	return nil
}

// UpdateStatusIf performs an atomic session status transition.
func (r *GameSessionRepository) UpdateStatusIf(ctx context.Context, id uint, from, to string) error {
	res := db(ctx, r.db).Model(&model.GameSession{}).Where("id = ? AND status = ?", id, from).Update("status", to)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return util.ErrConflict
	}
	return nil
}

// IncrementBooked atomically adds one booked count when still open.
func (r *GameSessionRepository) IncrementBooked(ctx context.Context, id uint) error {
	res := db(ctx, r.db).Model(&model.GameSession{}).
		Where("id = ? AND status = ? AND booked_count < max_players", id, "open").
		UpdateColumn("booked_count", gorm.Expr("booked_count + 1"))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return util.ErrConflict
	}
	return nil
}

// DecrementBooked reduces booked count when a player quits.
func (r *GameSessionRepository) DecrementBooked(ctx context.Context, id uint) error {
	return db(ctx, r.db).Model(&model.GameSession{}).Where("id = ? AND booked_count > 0", id).
		UpdateColumn("booked_count", gorm.Expr("booked_count - 1")).Error
}

// QueryByIDs returns sessions matching the given ids.
func (r *GameSessionRepository) QueryByIDs(ctx context.Context, ids []uint, out *[]model.GameSession) error {
	return db(ctx, r.db).Where("id IN ?", ids).Find(out).Error
}

// CreateRegistration inserts a session signup.
func (r *GameSessionRepository) CreateRegistration(ctx context.Context, reg *model.SessionRegistration) error {
	return db(ctx, r.db).Create(reg).Error
}

// FindRegistration returns a signup of a user for a session.
func (r *GameSessionRepository) FindRegistration(ctx context.Context, sessionID, userID uint) (*model.SessionRegistration, error) {
	var reg model.SessionRegistration
	err := db(ctx, r.db).Where("session_id = ? AND user_id = ?", sessionID, userID).First(&reg).Error
	if err != nil {
		return nil, normalizeError(err)
	}
	return &reg, nil
}

// DeleteRegistration removes a signup.
func (r *GameSessionRepository) DeleteRegistration(ctx context.Context, sessionID, userID uint) error {
	return db(ctx, r.db).Where("session_id = ? AND user_id = ?", sessionID, userID).Delete(&model.SessionRegistration{}).Error
}

// CountRegistrations returns the signup count of a session.
func (r *GameSessionRepository) CountRegistrations(ctx context.Context, sessionID uint) (int64, error) {
	var n int64
	err := db(ctx, r.db).Model(&model.SessionRegistration{}).Where("session_id = ?", sessionID).Count(&n).Error
	return n, err
}
