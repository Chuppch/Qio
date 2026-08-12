package user

import "regexp"

// 账号规则。
//
// 迁移自 v1 的 AccountValidator 与 UserConstants。
//
// v1 存在一处自相矛盾：UserConstants.USERNAME_MIN_LENGTH 声明为 6，但
// AccountValidator 的正则实际允许 4 位（首字母 + 3 位）。这里以正则为准，
// 统一为 4 位下限，并把长度常量与正则收在同一处，避免再次分叉。
const (
	// UsernameMinLength 用户名最短长度
	UsernameMinLength = 4
	// UsernameMaxLength 用户名最长长度
	UsernameMaxLength = 20

	// PasswordMinLength 密码最短长度
	PasswordMinLength = 6
	// PasswordMaxLength 密码最长长度
	PasswordMaxLength = 20
)

// usernamePattern 要求 4-20 位，由字母开头，其余为字母、数字或下划线。
var usernamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{3,19}$`)

// emailPattern 与 v1 保持一致。
var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$`)

// ValidUsername 校验用户名是否符合规则。
func ValidUsername(s string) bool {
	return usernamePattern.MatchString(s)
}

// ValidEmail 校验邮箱格式是否合法。
func ValidEmail(s string) bool {
	return emailPattern.MatchString(s)
}

// ValidPassword 校验密码长度是否在允许区间内。
func ValidPassword(s string) bool {
	n := len([]rune(s))
	return n >= PasswordMinLength && n <= PasswordMaxLength
}

// ValidAccount 校验登录账号，用户名或邮箱二者之一合法即可。
func ValidAccount(s string) bool {
	return ValidUsername(s) || ValidEmail(s)
}
