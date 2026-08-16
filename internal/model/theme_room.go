package model

import "time"

// ThemeRoom is an escape room theme.
type ThemeRoom struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:64;not null" json:"name"`
	Category       string    `gorm:"size:16;index;not null" json:"category"`
	DifficultyStars int      `gorm:"not null;default:3" json:"difficulty_stars"`
	MinPlayers     int       `gorm:"not null;default:2" json:"min_players"`
	MaxPlayers     int       `gorm:"not null;default:6" json:"max_players"`
	DurationMinutes int      `gorm:"not null;default:60" json:"duration_minutes"`
	Story          string    `gorm:"type:text" json:"story"`
	PosterURL      string    `gorm:"size:255" json:"poster_url"`
	SceneImages    string    `gorm:"type:text" json:"scene_images"`
	Status         string    `gorm:"size:16;not null;default:active" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}
