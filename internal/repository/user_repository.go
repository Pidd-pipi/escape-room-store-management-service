package repository

import (
	"context"

	"github.com/lp/escape-room-ops/internal/model"
	"gorm.io/gorm"
)

// UserRepository persists user rows.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository builds a UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user.
func (r *UserRepository) Create(ctx context.Context, u *model.User) error {
	return db(ctx, r.db).Create(u).Error
}

// FindByPhone returns the user with the given phone.
func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	var u model.User
	err := db(ctx, r.db).Where("phone = ?", phone).First(&u).Error
	if err != nil {
		return nil, normalizeError(err)
	}
	return &u, nil
}

// FindByID returns the user with the given id.
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var u model.User
	err := db(ctx, r.db).First(&u, id).Error
	if err != nil {
		return nil, normalizeError(err)
	}
	return &u, nil
}

// UpdateProfile updates the nickname.
func (r *UserRepository) UpdateProfile(ctx context.Context, id uint, nickname string) error {
	return db(ctx, r.db).Model(&model.User{}).Where("id = ?", id).Update("nickname", nickname).Error
}

// Count returns the total number of users.
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	var n int64
	err := db(ctx, r.db).Model(&model.User{}).Count(&n).Error
	return n, err
}
