package user

// 成长体系：每日任务与签到奖励。
//
// Task 与 SignAward 是独立于 User 的聚合——它们有自己的仓储（TaskRepository），
// 能被单独加载和保存，因此不放在 model.go 中。
//
// 两者在 v1 均无持久化表：任务模板与完成状态在 Redis，签到奖励由运营配置，
// 仓储实现见 internal/infrastructure/redis。

// TaskStatus 是每日任务的完成状态。
type TaskStatus int

const (
	TaskPending  TaskStatus = 0 // 未完成
	TaskFinished TaskStatus = 1 // 已完成未领取
	TaskClaimed  TaskStatus = 2 // 已完成已领取
)

// Claimable 表示任务奖励可以领取。
func (s TaskStatus) Claimable() bool { return s == TaskFinished }

// Task 是一条每日任务。
type Task struct {
	ID          int64
	Name        string
	Description string
	Status      TaskStatus

	// Reward 是完成后获得的猪仔钱
	Reward int

	// Link 与 Route 供前端跳转到任务对应的页面
	Link  string
	Route string
}

// SignAward 是累计签到达到指定天数后可获得的奖励。
type SignAward struct {
	ID          int64
	Type        ItemType
	ItemID      int64
	Name        string
	Description string
	PreviewURL  string

	// Count 是奖励数量，猪仔钱类型时表示金额
	Count int

	// RequiredDays 是触发该奖励所需的累计签到天数
	RequiredDays int
}
