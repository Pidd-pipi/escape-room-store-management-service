package util

import (
	"fmt"
	"time"

	"github.com/lp/escape-room-ops/internal/constants"
)

// FormatDateTime renders a time value using the shared display layout.
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

// FormatDate renders a date-only display string.
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatMoney renders a money value with two decimals.
func FormatMoney(v float64) string {
	return fmt.Sprintf("¥%.2f", v)
}

// SessionStatusText maps a session status to its Chinese label.
func SessionStatusText(s string) string {
	return constants.GameSessionStatusText(s)
}

// RoleText maps a user role to its Chinese label.
func RoleText(role string) string {
	return constants.UserRoleText(role)
}

// ThemeCategoryText maps a theme category to its Chinese label.
func ThemeCategoryText(c string) string {
	return constants.ThemeCategoryText(c)
}

// PropStatusText maps a prop status to its Chinese label.
func PropStatusText(s string) string {
	return constants.PropStatusText(s)
}
