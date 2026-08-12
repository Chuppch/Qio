package user

import "errors"

// 用户域的错误值。
//
// 逐条对应 v1 的异常类与 i18n 消息码，注释中给出出处以便核对。transport 层负责
// 把这些错误映射为 HTTP 状态码与业务码，domain 与 infra 只负责返回。
//
// 仓储实现须把存储层错误翻译为这里的错误，例如把 gorm.ErrRecordNotFound 转成
// ErrNotFound，避免存储细节泄漏到上层。
var (
	// ---- 账号 ----

	// ErrNotFound 用户不存在。对应 UserNotExistsException（user.not.exists）
	ErrNotFound = errors.New("user not found")

	// ErrLoginNotFound 登录账号不存在。对应 UserLoginNotExistsException（user.login.not.exists）
	ErrLoginNotFound = errors.New("login account not found")

	// ErrWrongPassword 密码不匹配。对应 UserPasswordNotMatchException（user.password.not.match）
	ErrWrongPassword = errors.New("password mismatch")

	// ErrPasswordRetryLimit 密码错误次数超限。对应 UserPasswordRetryLimitExceedException
	ErrPasswordRetryLimit = errors.New("password retry limit exceeded")

	// ErrPasswordConfirmMismatch 两次输入的密码不一致。对应 UserConfirmPasswordNotEqualsException
	ErrPasswordConfirmMismatch = errors.New("password confirmation mismatch")

	// ErrAccountDisabled 账号已停用或已删除
	ErrAccountDisabled = errors.New("account disabled")

	// ---- 注册与资料 ----

	// ErrUsernameTaken 用户名已被占用。对应 user.username.exists
	ErrUsernameTaken = errors.New("username already taken")

	// ErrUsernameInvalid 用户名格式不合法。对应 UserNameNotMatchException 与 user.username.length
	ErrUsernameInvalid = errors.New("invalid username")

	// ErrEmailTaken 邮箱已注册。对应 email.exists
	ErrEmailTaken = errors.New("email already registered")

	// ErrEmailInvalid 邮箱格式不合法。对应 email.format.error
	ErrEmailInvalid = errors.New("invalid email format")

	// ErrEmailNotFound 邮箱未注册。对应 email.not.exists
	ErrEmailNotFound = errors.New("email not registered")

	// ErrPasswordInvalid 密码长度不合法。对应 user.password.length
	ErrPasswordInvalid = errors.New("invalid password length")

	// ---- 验证码 ----

	// ErrCodeExpired 验证码已过期。对应 CodeTimeoutException（user.code.expire）
	ErrCodeExpired = errors.New("verify code expired")

	// ErrCodeMismatch 验证码不正确。对应 CodeErrorException（user.code.error）
	ErrCodeMismatch = errors.New("verify code mismatch")

	// ErrCodeSendLimit 验证码发送过于频繁。对应 user.sent.code.limit
	ErrCodeSendLimit = errors.New("verify code send rate limited")

	// ErrCodeGetLimit 获取图形验证码过于频繁。对应 user.get.code.limit
	ErrCodeGetLimit = errors.New("captcha request rate limited")

	// ErrCodeSendFailed 验证码发送失败。对应 user.sent.code.failed 与
	// user.sent.code.failed.by.email
	ErrCodeSendFailed = errors.New("verify code send failed")

	// ---- 地址簿 ----

	// ErrAddressNotFound 地址不存在。对应 user.address.not.exists
	ErrAddressNotFound = errors.New("address not found")

	// ErrDefaultAddressUndeletable 默认地址不允许删除。对应 user.address.default.not.delete
	ErrDefaultAddressUndeletable = errors.New("default address cannot be deleted")

	// ---- 余额与成长 ----

	// ErrNotEnoughMoney 余额不足。对应 user.money.not.enough
	ErrNotEnoughMoney = errors.New("insufficient balance")

	// ErrAlreadySigned 今日已签到。对应 user.sign.today
	ErrAlreadySigned = errors.New("already signed today")

	// ErrSignAwardFailed 签到奖励发放失败。对应 user.sign.award.error
	ErrSignAwardFailed = errors.New("sign award grant failed")

	// ErrTaskNotClaimable 任务未完成或奖励已领取
	ErrTaskNotClaimable = errors.New("task reward not claimable")
)
