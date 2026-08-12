package dict

import "errors"

var (
	// ErrAvatarNotFound 头像不存在
	ErrAvatarNotFound = errors.New("avatar not found")

	// ErrCountryNotFound 国家不存在。对应 v1 的 user.address.country.not.exists
	ErrCountryNotFound = errors.New("country not found")
)
