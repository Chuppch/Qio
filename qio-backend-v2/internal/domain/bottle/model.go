package bottle

import "time"

// Bottle 是一个漂流瓶。
//
// 状态很简单：投放后未被捞，被他人捞起后置为已捞，捞起者可以扔回使其重新可捞。
// 因此 Picked 是唯一的状态字段，不需要状态枚举。
type Bottle struct {
	ID int64

	// 投放者。Nickname 与 Email 是投放时的快照，
	// 使捞瓶子的人无需再查用户表即可展示与建立联系。
	UserID   int64
	Nickname string
	Email    string

	// SenderAddr 是投放者地址，捞起方同意加好友后会写入其好友记录。
	SenderAddr Address

	Content string

	// ImageURL 是瓶子的渲染图
	ImageURL string

	// Picked 表示是否已被他人捞起
	Picked bool

	// PickedBy 是捞起者的用户 ID。
	//
	// v1 复用审计列 update_user 承载这一语义——捞起与扔回都是更新操作，
	// 更新人即捞起者。等价迁移保留该做法，仓储层负责与 update_user 列互转。
	PickedBy int64

	Remark string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Available 表示该瓶子当前可被捞起。
func (b *Bottle) Available() bool { return !b.Picked }

// PickableBy 表示指定用户可以捞起该瓶子。
//
// 既不能捞自己投放的，也不能捞自己曾经捞过的——后者是 v1 用 update_user 排除的行为。
func (b *Bottle) PickableBy(userID int64) bool {
	return b.Available() && b.UserID != userID && b.PickedBy != userID
}

// Pick 把瓶子标记为被指定用户捞起。
func (b *Bottle) Pick(userID int64) {
	b.Picked = true
	b.PickedBy = userID
}

// ThrowBack 把瓶子扔回海里，重新可被他人捞起。
//
// 不清除 PickedBy：v1 保留 update_user，使同一个人不会重复捞到自己扔回的瓶子。
func (b *Bottle) ThrowBack() { b.Picked = false }

// Address 是投放者的地址。
type Address struct {
	ID               int64
	CountryID        int64
	FormattedAddress string
	Longitude        float64
	Latitude         float64
	IsDefault        bool
}
