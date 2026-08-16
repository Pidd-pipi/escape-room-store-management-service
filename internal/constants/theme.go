// Package constants centralizes business constants for the escape-room-ops backend.
package constants

// ThemeCategory defines escape room theme category enum values shared with the frontend.
const (
	ThemeCategoryHorror   = "horror"
	ThemeCategorySuspense = "suspense"
	ThemeCategoryScifi    = "scifi"
)

// ThemeCategories lists all valid theme categories.
var ThemeCategories = []string{ThemeCategoryHorror, ThemeCategorySuspense, ThemeCategoryScifi}

// IsThemeCategory reports whether the given category is valid.
func IsThemeCategory(c string) bool {
	for _, v := range ThemeCategories {
		if v == c {
			return true
		}
	}
	return false
}

// ThemeCategoryText returns the Chinese label of a theme category.
func ThemeCategoryText(c string) string {
	switch c {
	case ThemeCategoryHorror:
		return "恐怖"
	case ThemeCategorySuspense:
		return "悬疑"
	case ThemeCategoryScifi:
		return "科幻"
	default:
		return "未知"
	}
}
