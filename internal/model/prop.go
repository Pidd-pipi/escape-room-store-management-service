package model

import "time"

// Prop is a clue prop belonging to a theme room.
type Prop struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ThemeRoomID    uint      `gorm:"index;not null" json:"theme_room_id"`
	Name           string    `gorm:"size:64;not null" json:"name"`
	Category       string    `gorm:"size:24;not null" json:"category"`
	Stock          int       `gorm:"not null;default:1" json:"stock"`
	AlertThreshold int       `gorm:"not null;default:2" json:"alert_threshold"`
	Status         string    `gorm:"size:16;not null;default:normal" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// PropMaintenance is a maintenance/repair record for a prop.
type PropMaintenance struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PropID    uint      `gorm:"index;not null" json:"prop_id"`
	Type      string    `gorm:"size:24;not null" json:"type"`
	Note      string    `gorm:"type:text" json:"note"`
	CreatedAt time.Time `json:"created_at"`
}
