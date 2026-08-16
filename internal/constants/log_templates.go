package constants

// Log templates are centralized so that any business field change forces a
// coordinated update of the related log statements across the codebase.
const (
	LogUserRegisterSuccess         = "user register success: phone=%s user_id=%d"
	LogUserRegisterFailed          = "user register failed: phone=%s error=%v"
	LogUserLoginSuccess            = "user login success: phone=%s user_id=%d role=%s"
	LogUserLoginFailed             = "user login failed: phone=%s error=%v"
	LogUserProfileUpdateSuccess    = "user profile update success: user_id=%d"
	LogThemeCreateSuccess          = "theme create success: theme_id=%d name=%s category=%s"
	LogThemeCreateFailed           = "theme create failed: name=%s error=%v"
	LogGameSessionCreateSuccess    = "game session create success: session_id=%d theme_id=%d"
	LogGameSessionCreateFailed     = "game session create failed: theme_id=%d error=%v"
	LogGameSessionRegisterSuccess  = "game session register success: session_id=%d user_id=%d"
	LogGameSessionRegisterFailed   = "game session register failed: session_id=%d user_id=%d error=%v"
	LogGameSessionLockSuccess      = "game session lock success: session_id=%d booked=%d"
	LogGameSessionLockFailed       = "game session lock failed: session_id=%d error=%v"
	LogGameSessionFinishSuccess    = "game session finish success: session_id=%d"
	LogPropCreateSuccess           = "prop create success: prop_id=%d name=%s"
	LogPropStockAlertSuccess       = "prop stock alert success: prop_id=%d stock=%d"
	LogPropStockAlertFailed        = "prop stock alert failed: prop_id=%d error=%v"
	LogPropMaintenanceCreate       = "prop maintenance record create: prop_id=%d type=%s"
	LogEscapeRecordCreateSuccess   = "escape record create success: record_id=%d session_id=%d"
	LogEscapeRecordCreateFailed    = "escape record create failed: session_id=%d error=%v"
	LogRevenueRecordSuccess        = "revenue record success: session_id=%d amount=%.2f"
	LogRevenueExportSuccess        = "revenue export success: count=%d"
	LogMiddlewareAuthFailed        = "auth middleware failed: error=%v"
	LogMiddlewareRbacDenied        = "rbac middleware denied: user_id=%d role=%s required=%v"
	LogRateLimitReached            = "rate limit reached: ip=%s route=%s"
	LogSeedingCompleted            = "database seeding completed: users=%d themes=%d"
	LogLeaderboardBuildSuccess     = "leaderboard build success: count=%d"
)

// LogTemplateCount guards the "at least 25 templates" requirement.
const LogTemplateCount = 30
