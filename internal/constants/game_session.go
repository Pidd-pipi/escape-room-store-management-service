package constants

// GameSessionStatus defines session state machine values shared with the frontend.
const (
	GameSessionStatusOpen      = "open"
	GameSessionStatusLocked    = "locked"
	GameSessionStatusCancelled = "cancelled"
	GameSessionStatusFinished  = "finished"
)

// GameSessionStatuses lists all valid session statuses in flow order.
var GameSessionStatuses = []string{
	GameSessionStatusOpen, GameSessionStatusLocked, GameSessionStatusCancelled, GameSessionStatusFinished,
}

// IsGameSessionStatus reports whether the given status is valid.
func IsGameSessionStatus(s string) bool {
	for _, v := range GameSessionStatuses {
		if v == s {
			return true
		}
	}
	return false
}

// GameSessionStatusText returns the Chinese label of a session status.
func GameSessionStatusText(s string) string {
	switch s {
	case GameSessionStatusOpen:
		return "开放报名"
	case GameSessionStatusLocked:
		return "已锁定"
	case GameSessionStatusCancelled:
		return "已取消"
	case GameSessionStatusFinished:
		return "已结束"
	default:
		return "未知"
	}
}
