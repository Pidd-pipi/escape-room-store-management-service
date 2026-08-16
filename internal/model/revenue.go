package model

import "time"

// Revenue is a daily revenue record generated when a session finishes.
type Revenue struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SessionID   uint      `gorm:"index;not null" json:"session_id"`
	ThemeRoomID uint      `gorm:"index;not null" json:"theme_room_id"`
	Amount      float64   `gorm:"not null" json:"amount"`
	BookedCount int       `gorm:"not null" json:"booked_count"`
	Date        time.Time `gorm:"index;not null" json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}
