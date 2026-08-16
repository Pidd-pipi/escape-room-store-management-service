package constants

// Messages centralizes user-facing prompts, backend responses and log wording.
const (
	MsgOK                   = "ok"
	MsgValidationFailed     = "参数校验失败"
	MsgUnauthorized         = "未登录或登录已过期"
	MsgForbidden            = "没有操作权限"
	MsgNotFound             = "资源不存在"
	MsgConflict             = "资源状态冲突"
	MsgRateLimited          = "请求过于频繁，请稍后再试"
	MsgInternalError        = "服务器内部错误"
	MsgPhoneAlreadyUsed     = "该手机号已被注册"
	MsgPhoneOrPassword      = "手机号或密码错误"
	MsgSessionNotOpen       = "该场次当前不可报名"
	MsgSessionFull          = "该场次人数已满"
	MsgAlreadyRegistered    = "您已报名该场次"
	MsgNotRegistered        = "您未报名该场次"
	MsgSessionStatusInvalid = "当前场次状态不可操作"
	MsgPropLowStock         = "道具库存不足"
)
