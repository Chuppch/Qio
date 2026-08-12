package user

import "time"

// AccountStatus 是账号的启用状态。
type AccountStatus int

const (
	StatusActive   AccountStatus = 0 // 正常
	StatusDisabled AccountStatus = 1 // 停用
)

// Enabled 表示账号未被停用。
func (s AccountStatus) Enabled() bool { return s == StatusActive }

// ItemType 是道具类型。取值与签到奖励、商城共用同一套编码。
type ItemType int

const (
	ItemMoney        ItemType = 1 // 猪仔钱
	ItemFunctionCard ItemType = 2 // 功能卡
	ItemFont         ItemType = 3 // 字体
	ItemFontColor    ItemType = 4 // 字体颜色
	ItemPaper        ItemType = 5 // 信纸
	ItemCollection   ItemType = 6 // 其他收藏品，含印章
)

// User 是用户聚合根。
type User struct {
	ID       int64
	Username string
	Nickname string
	Email    string
	Sex      string
	Avatar   string

	// PasswordHash 只存摘要，不存明文
	PasswordHash string

	Status  AccountStatus
	Deleted bool

	// Money 是猪仔钱余额
	Money int64

	// Addresses 至多一条 IsDefault 为 true
	Addresses []Address

	OwnedItems []OwnedItem

	LastLoginIP string
	LastLoginAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Active 表示账号可正常使用，即未停用且未删除。
func (u *User) Active() bool { return u.Status.Enabled() && !u.Deleted }

// Affordable 表示余额足以支付指定金额。
func (u *User) Affordable(amount int64) bool {
	return amount >= 0 && u.Money >= amount
}

// DefaultAddress 返回默认地址，不存在时返回零值与 false。
func (u *User) DefaultAddress() (Address, bool) {
	for _, a := range u.Addresses {
		if a.IsDefault {
			return a, true
		}
	}
	return Address{}, false
}

// Owns 表示用户已拥有指定道具。
func (u *User) Owns(t ItemType, itemID int64) bool {
	for _, it := range u.OwnedItems {
		if it.Type == t && it.ItemID == itemID {
			return true
		}
	}
	return false
}

// CodeScene 区分邮箱验证码的用途。
type CodeScene string

const (
	CodeSceneRegister      CodeScene = "register"       // 注册
	CodeSceneResetPassword CodeScene = "reset-password" // 重置密码
)

// Address 是地址簿中的一条地址，随 User 一起存取。
type Address struct {
	ID               int64
	CountryID        int64
	FormattedAddress string
	Longitude        float64
	Latitude         float64
	IsDefault        bool
}

// OwnedItem 是背包中的一件道具，随 User 一起存取。
//
// Count 用于可叠加的道具（功能卡）；字体、信纸等不可叠加的恒为 1。
type OwnedItem struct {
	Type   ItemType
	ItemID int64
	Count  int
}
