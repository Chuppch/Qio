package dict

import "context"

// Repository 是字典数据的访问接口，实现由 internal/infrastructure/mysql 提供。
//
// 列表方法供展示使用；单条查询在找不到时返回 ErrAvatarNotFound 或
// ErrCountryNotFound，调用方以此判断 ID 是否有效，因此不额外提供 Exists 方法。
type Repository interface {
	ListAvatars(ctx context.Context) ([]Avatar, error)
	FindAvatar(ctx context.Context, id int64) (Avatar, error)

	ListCountries(ctx context.Context) ([]Country, error)
	FindCountry(ctx context.Context, id int64) (Country, error)
}
