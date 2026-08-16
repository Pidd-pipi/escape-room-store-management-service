package constants

// PropCategory defines prop category enum values shared with the frontend.
const (
	PropCategoryKey        = "key"
	PropCategoryPasswordBox = "password_box"
	PropCategoryMechanism  = "mechanism"
)

// PropStatus defines prop inventory state enum values shared with the frontend.
const (
	PropStatusNormal     = "normal"
	PropStatusLowStock   = "low_stock"
	PropStatusRepairing  = "repairing"
)

// IsPropStatus reports whether the given status is valid.
func IsPropStatus(s string) bool {
	switch s {
	case PropStatusNormal, PropStatusLowStock, PropStatusRepairing:
		return true
	default:
		return false
	}
}

// PropStatusText returns the Chinese label of a prop status.
func PropStatusText(s string) string {
	switch s {
	case PropStatusNormal:
		return "正常"
	case PropStatusLowStock:
		return "库存不足"
	case PropStatusRepairing:
		return "维修中"
	default:
		return "未知"
	}
}
