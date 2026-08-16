// Package model defines the GORM entity structures for escape-room-ops.
package model

import "time"

// User represents a player or admin.
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Phone        string    `gorm:"size:20;uniqueIndex;not null" json:"phone"`
	PasswordHash string    `gorm:"size:100;not null" json:"-"`
	Nickname     string    `gorm:"size:32;not null" json:"nickname"`
	Role         string    `gorm:"size:16;not null;default:player" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
