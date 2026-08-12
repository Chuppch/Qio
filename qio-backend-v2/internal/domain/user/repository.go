package user

import (
	"context"
	"time"
)

// Repository 是用户账号的数据访问接口，实现由 internal/infra/mysql 提供。
//
// 方法集覆盖 v1 UserServiceImpl 中对 userMapper 的全部调用场景。
type Repository interface {
	// Create 注册一个新用户，成功后回填 ID。
	Create(ctx context.Context, u *User) error

	// FindByID 按主键查询，不存在时返回 ErrNotFound。
	FindByID(ctx context.Context, id int64) (*User, error)

	// FindByIDs 按主键批量查询，用于好友列表等需要批量补齐用户信息的场景。
	//
	// 返回结果不保证与入参顺序一致，也不保证数量相同——不存在的 ID 会被跳过。
	FindByIDs(ctx context.Context, ids []int64) ([]*User, error)

	// FindByEmail 按邮箱查询，用于登录、重置密码与收信人关联。
	FindByEmail(ctx context.Context, email string) (*User, error)

	// FindByUsername 按用户名查询，用于登录与重名校验。
	FindByUsername(ctx context.Context, username string) (*User, error)

	// Update 按主键更新用户的可变字段，含余额、地址簿与背包。
	Update(ctx context.Context, u *User) error

	// UpdateByEmail 按邮箱更新用户，用于重置密码——此时只有邮箱没有主键。
	UpdateByEmail(ctx context.Context, u *User) error

	// UpdateMoney 按增量调整余额，delta 可为负。
	//
	// 单独提供该方法而不复用 Update：余额在发信、读信、购买、签到发奖多处变动，
	// 走「读整个 User 再整体写回」在并发下会覆盖其他字段的修改。
	UpdateMoney(ctx context.Context, userID int64, delta int64) error

	// EmailTaken 判断邮箱是否已被占用。
	EmailTaken(ctx context.Context, email string) (bool, error)

	// UsernameTaken 判断用户名是否已被占用。
	UsernameTaken(ctx context.Context, username string) (bool, error)
}

// VerifyCodeRepository 是验证码的读写接口，实现由 internal/infra/redis 提供。
//
// 覆盖三种验证码：登录图形验证码（键为随机 uuid）、注册邮箱验证码与重置密码
// 邮箱验证码（键为邮箱地址）。v1 中三者共用同一套裸键，v2 由实现层负责加前缀
// 隔离，接口只区分用途。
type VerifyCodeRepository interface {
	// SaveCaptcha 保存图形验证码，key 为随机生成的 uuid。
	SaveCaptcha(ctx context.Context, key, code string, ttl time.Duration) error

	// TakeCaptcha 取出并删除图形验证码，不存在或已过期时返回 ErrNotFound。
	TakeCaptcha(ctx context.Context, key string) (string, error)

	// SaveEmailCode 保存邮箱验证码，scene 区分注册与重置密码两种用途。
	SaveEmailCode(ctx context.Context, scene CodeScene, email, code string, ttl time.Duration) error

	// FindEmailCode 查询邮箱验证码，不存在或已过期时返回 ErrNotFound。
	FindEmailCode(ctx context.Context, scene CodeScene, email string) (string, error)

	// DeleteEmailCode 删除邮箱验证码，校验通过后调用。
	DeleteEmailCode(ctx context.Context, scene CodeScene, email string) error

	// EmailCodePending 判断该邮箱是否已有未过期的验证码，用于限制重复发送。
	EmailCodePending(ctx context.Context, scene CodeScene, email string) (bool, error)
}

// TaskRepository 是每日任务与签到的数据访问接口。
//
// 实现落在 internal/infra/redis：任务配置与签到记录在 v1 中都存于 Redis，
// 没有对应数据表。
type TaskRepository interface {
	// ListTasks 返回指定用户在指定日期的任务及其完成状态。
	//
	// 该用户当日无记录时，实现需从任务模板复制一份后返回。
	ListTasks(ctx context.Context, userID int64, date string) ([]Task, error)

	// MarkTaskFinished 标记一条任务已完成。
	MarkTaskFinished(ctx context.Context, userID int64, date string, taskID int64) error

	// MarkTaskClaimed 标记一条任务的奖励已领取。
	MarkTaskClaimed(ctx context.Context, userID int64, date string, taskID int64) error

	// Sign 记录一次签到，返回签到后的累计天数。
	Sign(ctx context.Context, userID int64, date string) (days int, err error)

	// Signed 判断指定日期是否已签到。
	Signed(ctx context.Context, userID int64, date string) (bool, error)

	// SignedDays 返回指定月份内已签到的日期，用于渲染签到日历。
	//
	// month 形如 202608，返回的日期形如 20260811。
	SignedDays(ctx context.Context, userID int64, month string) ([]string, error)

	// ListSignAwards 返回指定日期生效的签到奖励配置。
	ListSignAwards(ctx context.Context, date string) ([]SignAward, error)
}
