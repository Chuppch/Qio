package bottle

import "context"

// Repository 是漂流瓶的数据访问接口，实现由 internal/infrastructure/mysql 提供。
//
// 方法集对应 v1 BottleServiceImpl 中对 bottleMapper 的全部调用。
type Repository interface {
	// Create 投放一个漂流瓶，成功后回填 ID。
	Create(ctx context.Context, b *Bottle) error

	// FindByID 按主键查询，不存在时返回 ErrNotFound。
	FindByID(ctx context.Context, id int64) (*Bottle, error)

	// ListAvailable 查询指定用户可捞的漂流瓶。
	//
	// 条件为未被捞起、非该用户投放、且非该用户曾捞起过。v1 是全表捞出后在应用层
	// 随机取一个，本方法保持同样语义，随机选择留在 service。
	ListAvailable(ctx context.Context, userID int64) ([]*Bottle, error)

	// FindLatestPicked 查询指定用户最近捞起的一个漂流瓶，无记录时返回 ErrNotFound。
	//
	// 用于「扔回」与「通过漂流瓶加好友」——两者都作用于当前手上那个瓶子。
	FindLatestPicked(ctx context.Context, userID int64) (*Bottle, error)

	// Update 更新漂流瓶的可捞状态与捞起者。
	Update(ctx context.Context, b *Bottle) error
}
