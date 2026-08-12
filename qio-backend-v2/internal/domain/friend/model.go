package friend

import "time"

// RequestStatus 是好友申请的处理状态。
//
// 取值沿用 v1 FriendConstants，避免存量数据迁移。v1 的常量表末尾留了一行
// 「已送达」注释但没有对应常量，说明曾计划第四个状态而未实现，这里不予保留。
type RequestStatus int

const (
	RequestPending  RequestStatus = 0 // 待处理
	RequestAccepted RequestStatus = 1 // 已同意
	RequestRejected RequestStatus = 2 // 已拒绝
)

// Pending 表示申请仍待处理。
func (s RequestStatus) Pending() bool { return s == RequestPending }

// Friend 是一条好友关系记录。
//
// 单向存储：OwnerID 是这条记录归属的用户，UserID 是对方。互为好友时表中存在
// 两条方向相反的记录，由 v1 的 addFriend 一次插入两行建立。
//
// Name / Sex / Email / Avatar 是对方用户资料的快照，v1 在读取好友列表时会用
// 用户表的最新值回写这几列（见 docs/TODO-migration.md 中的读路径写操作）。
type Friend struct {
	ID int64

	// UserID 是好友本人的用户 ID。
	//
	// v1 的字段注释标为「非必需项」，表中也允许为空，因此可能为 0。
	UserID int64

	// OwnerID 对应表中的 owning_id。
	//
	// v1 的 Java 字段名是 OwningId，首字母大写，靠 MyBatis-Plus 的驼峰转换
	// 推导出列名。这里改用 OwnerID，列名映射由仓储层负责。
	OwnerID int64

	Name   string
	Sex    string
	Email  string
	Avatar string

	// Addresses 是该好友的收件地址簿，至多一条 IsDefault 为 true。
	Addresses []Address

	// Remark 是好友备注名。
	//
	// v1 复用了审计基类的 remark 列承载这一业务语义，建立好友关系时默认填
	// 对方昵称。
	Remark string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// OwnedBy 表示该记录归属指定用户。
//
// v1 的多数好友操作都先按 (id, owning_id) 查询来隐式完成归属校验，
// 但 updateFriendRemark 漏掉了这一步（见 docs/TODO-migration.md）。
func (f *Friend) OwnedBy(userID int64) bool { return f.OwnerID == userID }

// FindAddress 按地址 ID 查找地址簿中的地址，不存在时返回零值与 false。
func (f *Friend) FindAddress(addressID int64) (Address, bool) {
	for _, a := range f.Addresses {
		if a.ID == addressID {
			return a, true
		}
	}
	return Address{}, false
}

// DefaultAddress 返回默认地址，不存在时返回零值与 false。
func (f *Friend) DefaultAddress() (Address, bool) {
	for _, a := range f.Addresses {
		if a.IsDefault {
			return a, true
		}
	}
	return Address{}, false
}

// SetDefaultAddress 把指定地址设为默认，其余地址取消默认。
//
// 对应 v1 setFriendDefaultAddress 中的遍历赋值。地址不存在时返回
// ErrAddressNotFound，与 v1 的先校验后赋值顺序一致。
func (f *Friend) SetDefaultAddress(addressID int64) error {
	if _, ok := f.FindAddress(addressID); !ok {
		return ErrAddressNotFound
	}
	for i := range f.Addresses {
		f.Addresses[i].IsDefault = f.Addresses[i].ID == addressID
	}
	return nil
}

// RemoveAddress 从地址簿中删除指定地址。
//
// 对应 v1 deleteFriendAddress：地址不存在返回 ErrAddressNotFound，
// 默认地址不允许删除，返回 ErrDefaultAddressUndeletable。
func (f *Friend) RemoveAddress(addressID int64) error {
	target, ok := f.FindAddress(addressID)
	if !ok {
		return ErrAddressNotFound
	}
	if target.IsDefault {
		return ErrDefaultAddressUndeletable
	}

	kept := make([]Address, 0, len(f.Addresses))
	for _, a := range f.Addresses {
		if a.ID != addressID {
			kept = append(kept, a)
		}
	}
	f.Addresses = kept
	return nil
}

// EnsureAddress 确保指定地址存在于地址簿中，返回补全 ID 后的地址，
// 以及地址簿是否发生了变化。
//
// 迁移自 v1 sendLetterPre 中那段内联的地址处理逻辑（原注释写着「顺带做个地址
// 处理你再更新」）。行为逐条对应：
//
//   - 地址簿为空：入参补上 ID 1 并置为默认，随后加入
//   - 地址簿非空：按格式化地址或 ID 命中已有地址，命中则复用其 ID 且不改动地址簿
//   - 未命中且入参无 ID：取末位地址 ID 加一
//   - 未命中且入参已有 ID：沿用该 ID，v1 不检测冲突，这里同样不检测
//
// v1 用 null 表示「入参无 ID」，Go 侧以 0 表达。因此按 ID 比对时跳过 0，
// 与 v1 中 equals(null) 恒为 false 的行为一致。
func (f *Friend) EnsureAddress(addr Address) (Address, bool) {
	if len(f.Addresses) == 0 {
		addr.ID = 1
		addr.IsDefault = true
		f.Addresses = append(f.Addresses, addr)
		return addr, true
	}

	for _, existing := range f.Addresses {
		sameText := existing.FormattedAddress == addr.FormattedAddress
		sameID := addr.ID != 0 && existing.ID == addr.ID
		if sameText || sameID {
			addr.ID = existing.ID
			return addr, false
		}
	}

	if addr.ID == 0 {
		addr.ID = f.Addresses[len(f.Addresses)-1].ID + 1
	}
	f.Addresses = append(f.Addresses, addr)
	return addr, true
}

// Address 是好友地址簿中的一条收件地址。
//
// 与 user.Address、bottle.Address 结构相同但语义不同：这里是「我记录的对方
// 收件地址」，可由发信流程自动补充。刻意不共享类型，避免任一侧调整波及对方的
// 存储格式。
type Address struct {
	ID               int64
	CountryID        int64
	FormattedAddress string
	Longitude        float64
	Latitude         float64
	IsDefault        bool
}

// Request 是一条好友申请。
//
// v1 中申请只在两个场景产生：捞到漂流瓶后主动申请，以及读到陌生人来信时自动
// 生成。两者分别由 BottleID、LetterID 标识来源。申请被处理后只改状态，不删除记录。
type Request struct {
	ID int64

	SenderID   int64
	ReceiverID int64

	Status RequestStatus

	// GiveAddress 是申请人提供给对方的地址。
	//
	// 同意后写入接收方那条好友记录的地址簿；对称方向的地址取自来源漂流瓶或
	// 信件的寄件地址。
	GiveAddress Address

	Content string

	// BottleID、LetterID 标识申请来源，二者互斥。
	//
	// 表中两列均可为空，v1 靠判空来区分来源。Go 侧以 0 表示「无」——两张来源表的
	// 主键都自增且从 1 起，不会与 0 冲突。
	BottleID int64
	LetterID int64

	// Remark 对应表中的 remark 列，好友申请场景下 v1 未使用。
	Remark string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// ReceivedBy 表示当前用户是该申请的接收方，即有权处理它。
func (r *Request) ReceivedBy(userID int64) bool { return r.ReceiverID == userID }

// FromBottle 表示该申请由漂流瓶发起。
func (r *Request) FromBottle() bool { return r.BottleID != 0 }

// FromLetter 表示该申请由信件发起。
func (r *Request) FromLetter() bool { return r.LetterID != 0 }

// Accept 同意申请。
//
// 已被处理过的申请返回 ErrRequestUnavailable，对应 v1 becomeFriend 中
// status != 0 时抛出的异常。
func (r *Request) Accept() error { return r.transit(RequestAccepted) }

// Reject 拒绝申请，语义同 Accept。
func (r *Request) Reject() error { return r.transit(RequestRejected) }

func (r *Request) transit(to RequestStatus) error {
	if !r.Status.Pending() {
		return ErrRequestUnavailable
	}
	r.Status = to
	return nil
}
