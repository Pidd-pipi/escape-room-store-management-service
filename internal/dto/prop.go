package dto

// CreatePropRequest is the payload for adding a prop.
type CreatePropRequest struct {
	ThemeRoomID    uint   `json:"theme_room_id" binding:"required"`
	Name           string `json:"name" binding:"required,min=1,max=64"`
	Category       string `json:"category" binding:"required"`
	Stock          int    `json:"stock" binding:"gte=0"`
	AlertThreshold int    `json:"alert_threshold" binding:"gte=0"`
}

// CreateMaintenanceRequest is the payload for a maintenance record.
type CreateMaintenanceRequest struct {
	Type string `json:"type" binding:"required"`
	Note string `json:"note"`
}

// ConsumeRequest is the payload for consuming props.
type ConsumeRequest struct {
	Quantity int `json:"quantity" binding:"required,gte=1"`
}
