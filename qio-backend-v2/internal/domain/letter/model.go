package letter

import "time"

// Status 是信件的投递状态。
//
// 取值沿用 v1，避免存量数据迁移。
type Status int

const (
	StatusNotSent   Status = 0 // 未发送
	StatusSent      Status = 1 // 已发送
	StatusTransit   Status = 2 // 传递中
	StatusDelivered Status = 3 // 已送达
)

// ReadStatus 是信件的阅读状态。
type ReadStatus int

const (
	Unread ReadStatus = 0 // 未读
	Read   ReadStatus = 1 // 已读
)

// Layout 是信件的排版方向。
type Layout int

const (
	LayoutHorizontal Layout = 1 // 横版
	LayoutVertical   Layout = 2 // 竖版字体
)

// ProgressScale 是投递进度的分母，进度用万分制表示。
//
// 例如 5000 表示 50%，ProgressScale 表示已送达。
const ProgressScale int64 = 10000

// Letter 是信件领域模型。
//
// 字段与 v1 的 letter 表语义对齐，但 ReduceTime 与 SpeedRate 在 v1 中被存成
// 字符串，这里改为强类型：ReduceTime 用 time.Duration，SpeedRate 用 float64。
// 迁移时需要在仓储层做一次转换。
type Letter struct {
	ID         int64
	Status     Status
	ReadStatus ReadStatus
	Layout     Layout

	// 投递进度，万分制，取值范围 [0, ProgressScale]
	DeliveryProgress int64

	CreatedAt time.Time
	// ExpectedDeliveryAt 是未使用任何加速道具时的预计送达时间，创建后不再变更，
	// 作为进度计算的基准总时长。
	ExpectedDeliveryAt time.Time
	// DeliveryAt 是叠加加速与减时后的实际预计送达时间，会随道具使用而变化。
	DeliveryAt time.Time

	// ReduceTime 是道具直接减免的时长
	ReduceTime time.Duration
	// SpeedRate 是加速倍率，1 表示不加速
	SpeedRate float64
}

// InTransit 表示信件是否处于传递中。
//
// 只有传递中的信件需要重新计算投递进度。
func (l *Letter) InTransit() bool {
	return l.Status == StatusTransit
}

// IsRead 表示信件是否已被收件人读过。
func (l *Letter) IsRead() bool {
	return l.ReadStatus == Read
}
