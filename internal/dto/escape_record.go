package dto

// CreateEscapeRecordRequest is the payload for recording a team outcome.
type CreateEscapeRecordRequest struct {
	SessionID       uint   `json:"session_id" binding:"required"`
	ThemeRoomID     uint   `json:"theme_room_id" binding:"required"`
	TeamName        string `json:"team_name" binding:"required,min=1,max=64"`
	DurationMinutes int    `json:"duration_minutes" binding:"gte=1"`
	HintCount       int    `json:"hint_count" binding:"gte=0"`
	Escaped         bool   `json:"escaped"`
}
