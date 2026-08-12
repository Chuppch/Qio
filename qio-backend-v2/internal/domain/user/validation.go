package user

import "regexp"

// 账号格式规则，注册、改名、登录共用。
//
// 用户名下限取 4 位：v1 的常量声明为 6，但正则实际允许 4 位，此处以正则为准。
const (
	UsernameMinLength = 4
	UsernameMaxLength = 20
	PasswordMinLength = 6
	PasswordMaxLength = 20
)

var (
	// 字母开头，其余为字母、数字或下划线
	usernamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{3,19}$`)
	emailPattern    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$`)
)

// ValidUsername 校验用户名格式。
func ValidUsername(s string) bool { return usernamePattern.MatchString(s) }

// ValidEmail 校验邮箱格式。
func ValidEmail(s string) bool { return emailPattern.MatchString(s) }

// ValidPassword 校验密码长度。
func ValidPassword(s string) bool {
	n := len([]rune(s))
	return n >= PasswordMinLength && n <= PasswordMaxLength
}

// ValidAccount 校验登录账号，用户名或邮箱合法其一即可。
func ValidAccount(s string) bool { return ValidUsername(s) || ValidEmail(s) }
