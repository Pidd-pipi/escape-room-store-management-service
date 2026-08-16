package dto

import "time"

// CreateGameSessionRequest is the payload for scheduling a session.
type CreateGameSessionRequest struct {
	ThemeRoomID uint      `json:"theme_room_id" binding:"required"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	MaxPlayers  int       `json:"max_players" binding:"required,gte=1"`
}

// ListSessionQuery adds theme/status filters to pagination.
type ListSessionQuery struct {
	PageQuery
	ThemeRoomID uint   `form:"theme_room_id"`
	Status      string `form:"status"`
}
