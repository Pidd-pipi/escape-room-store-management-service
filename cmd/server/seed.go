package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/lp/escape-room-ops/internal/constants"
	"github.com/lp/escape-room-ops/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// seed populates default users, theme rooms, sessions and props on first boot.
func seed(ctx context.Context, db *gorm.DB, logger *slog.Logger) error {
	var total int64
	if err := db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return fmt.Errorf("seed count users: %w", err)
	}
	if total > 0 {
		logger.Info("seed skipped: users already exist", slog.Int64("count", total))
		return nil
	}
	hash := func(p string) (string, error) {
		b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	type seedUser struct {
		phone, password, nickname, role string
	}
	defs := []seedUser{
		{phone: "13300000001", password: "123456", nickname: "密室玩家甲", role: constants.UserRolePlayer},
		{phone: "13300000002", password: "123456", nickname: "密室玩家乙", role: constants.UserRolePlayer},
		{phone: "13400000001", password: "admin123", nickname: "店长", role: constants.UserRoleAdmin},
	}
	users := make([]model.User, 0, len(defs))
	for _, d := range defs {
		h, err := hash(d.password)
		if err != nil {
			return fmt.Errorf("seed hash user %s: %w", d.phone, err)
		}
		users = append(users, model.User{
			Phone: d.phone, PasswordHash: h, Nickname: d.nickname, Role: d.role,
		})
	}
	if err := db.WithContext(ctx).Create(&users).Error; err != nil {
		return fmt.Errorf("seed users: %w", err)
	}
	themes := []model.ThemeRoom{
		{Name: "午夜医院", Category: constants.ThemeCategoryHorror, DifficultyStars: 4, MinPlayers: 2, MaxPlayers: 6, DurationMinutes: 60, Story: "深夜的废弃医院传来哭声……", Status: "active"},
		{Name: "雾都疑云", Category: constants.ThemeCategorySuspense, DifficultyStars: 3, MinPlayers: 2, MaxPlayers: 5, DurationMinutes: 60, Story: "雾都连环失踪案等你破解", Status: "active"},
		{Name: "星际方舟", Category: constants.ThemeCategoryScifi, DifficultyStars: 5, MinPlayers: 3, MaxPlayers: 6, DurationMinutes: 90, Story: "飞船即将坠毁，找到逃生舱", Status: "active"},
	}
	if err := db.WithContext(ctx).Create(&themes).Error; err != nil {
		return fmt.Errorf("seed themes: %w", err)
	}
	sessions := []model.GameSession{
		{ThemeRoomID: themes[0].ID, StartTime: time.Now().Add(24 * time.Hour), MaxPlayers: 6, BookedCount: 2, Status: constants.GameSessionStatusOpen},
		{ThemeRoomID: themes[1].ID, StartTime: time.Now().Add(48 * time.Hour), MaxPlayers: 5, BookedCount: 5, Status: constants.GameSessionStatusLocked},
		{ThemeRoomID: themes[2].ID, StartTime: time.Now().Add(72 * time.Hour), MaxPlayers: 6, BookedCount: 0, Status: constants.GameSessionStatusOpen},
	}
	if err := db.WithContext(ctx).Create(&sessions).Error; err != nil {
		return fmt.Errorf("seed sessions: %w", err)
	}
	props := []model.Prop{
		{ThemeRoomID: themes[0].ID, Name: "病房钥匙", Category: constants.PropCategoryKey, Stock: 5, AlertThreshold: 2, Status: constants.PropStatusNormal},
		{ThemeRoomID: themes[0].ID, Name: "密码箱", Category: constants.PropCategoryPasswordBox, Stock: 2, AlertThreshold: 1, Status: constants.PropStatusNormal},
		{ThemeRoomID: themes[2].ID, Name: "齿轮机关", Category: constants.PropCategoryMechanism, Stock: 3, AlertThreshold: 2, Status: constants.PropStatusNormal},
	}
	if err := db.WithContext(ctx).Create(&props).Error; err != nil {
		return fmt.Errorf("seed props: %w", err)
	}
	logger.Info(fmt.Sprintf(constants.LogSeedingCompleted, len(users), len(themes)))
	return nil
}
