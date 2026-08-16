// Package dto defines request and response payloads for the escape-room-ops API.
package dto

// PageQuery is the unified pagination query.
type PageQuery struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
}

// Pagination defaults and bounds.
const (
	defaultPage     = 1
	defaultPageSize = 10
	maxPageSize     = 100
)

// Normalize fills pagination defaults and clamps out-of-range values.
func (p *PageQuery) Normalize() {
	if p.Page <= 0 {
		p.Page = defaultPage
	}
	if p.PageSize <= 0 || p.PageSize > maxPageSize {
		p.PageSize = defaultPageSize
	}
}

// PageResult is the standard paginated payload.
type PageResult struct {
	Items    interface{} `json:"items"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}
