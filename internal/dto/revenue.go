package dto

// RevenueQuery is the analytics filter (day/week/month export).
type RevenueQuery struct {
	Range string `form:"range" binding:"omitempty,oneof=day week month"`
}
