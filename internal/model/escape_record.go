package model

import "time"

// EscapeRecord is a team's escape outcome for the leaderboard.
type EscapeRecord struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	SessionID       uint      `gorm:"index;not null" json:"session_id"`
	ThemeRoomID     uint      `gorm:"index;not null" json:"theme_room_id"`
	TeamName        string    `gorm:"size:64;not null" json:"team_name"`
	DurationMinutes int       `gorm:"not null" json:"duration_minutes"`
	HintCount       int       `gorm:"not null;default:0" json:"hint_count"`
	Escaped         bool      `gorm:"not null" json:"escaped"`
	CreatedAt       time.Time `json:"created_at"`
}
