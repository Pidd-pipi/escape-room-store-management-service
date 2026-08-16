// Package dto defines request and response payloads for the escape-room-ops API.
package dto

// PageQuery is the unified pagination query.
type PageQuery struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

// Normalize fills pagination defaults.
func (p *PageQuery) Normalize() {
	if p.Page <= 0 {
		p.Page = 0
	}
	if p.PageSize <= 0 || p.PageSize > 100 {
		p.PageSize = 0
	}
}

// PageResult is the standard paginated payload.
type PageResult struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}
