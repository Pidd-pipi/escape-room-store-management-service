package dto

// CreateThemeRoomRequest is the payload for adding a theme room.
type CreateThemeRoomRequest struct {
	Name            string `json:"name" binding:"required,min=1,max=64"`
	Category        string `json:"category" binding:"required"`
	DifficultyStars int    `json:"difficulty_stars" binding:"required,gte=1,lte=5"`
	MinPlayers      int    `json:"min_players" binding:"required,gte=1"`
	MaxPlayers      int    `json:"max_players" binding:"required,gte=1"`
	DurationMinutes int    `json:"duration_minutes" binding:"required,gte=10"`
	Story           string `json:"story"`
	PosterURL       string `json:"poster_url"`
	SceneImages     string `json:"scene_images"`
}

// ListThemeQuery adds category/difficulty filters to pagination.
type ListThemeQuery struct {
	PageQuery
	Category string `form:"category"`
	Stars    int    `form:"stars"`
}
