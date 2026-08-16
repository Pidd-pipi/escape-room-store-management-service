package model

import "time"

// SessionRegistration is a player's signup for a session.
type SessionRegistration struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SessionID uint      `gorm:"uniqueIndex:uniq_session_registration,priority:1;not null" json:"session_id"`
	UserID    uint      `gorm:"uniqueIndex:uniq_session_registration,priority:2;not null" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// GameSession is a scheduled play session for a theme room.
type GameSession struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ThemeRoomID uint      `gorm:"index;not null" json:"theme_room_id"`
	StartTime   time.Time `gorm:"not null" json:"start_time"`
	MaxPlayers  int       `gorm:"not null;default:6" json:"max_players"`
	BookedCount int       `gorm:"not null;default:0" json:"booked_count"`
	Status      string    `gorm:"size:16;index;not null;default:open" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
