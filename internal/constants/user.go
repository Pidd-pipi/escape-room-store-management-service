package constants

// UserRole defines user role enum values shared with the frontend.
const (
	UserRolePlayer = "player"
	UserRoleAdmin  = "admin"
)

// UserRoles lists all valid user roles.
var UserRoles = []string{UserRolePlayer, UserRoleAdmin}

// IsUserRole reports whether the given role is valid.
func IsUserRole(r string) bool {
	for _, v := range UserRoles {
		if v == r {
			return true
		}
	}
	return false
}

// UserRoleText returns the Chinese label of a user role.
func UserRoleText(r string) string {
	switch r {
	case UserRolePlayer:
		return "玩家"
	case UserRoleAdmin:
		return "管理员"
	default:
		return "未知"
	}
}
