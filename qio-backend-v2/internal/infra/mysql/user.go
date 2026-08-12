package mysql

import "time"

// userPO 映射 user 表。
type userPO struct {
	auditFields

	Username string `gorm:"column:username"`
	Nickname string `gorm:"column:nickname"`
	Email    string `gorm:"column:email"`
	Sex      string `gorm:"column:sex"`
	Avatar   string `gorm:"column:avatar"`

	// Password 列存的是密码摘要。
	Password string `gorm:"column:password"`

	// Status、DelFlag 在表中是 varchar，取值 "0" / "1"。
	Status  string `gorm:"column:status"`
	DelFlag string `gorm:"column:del_flag"`

	LoginIP   string    `gorm:"column:login_ip"`
	LoginDate time.Time `gorm:"column:login_date"`

	Money int64 `gorm:"column:money"`

	// 以下六列均为 JSON，存的是数组。
	// Addresses 元素结构见 addressJSON，其余见 ownedItemJSON。
	Addresses     []byte `gorm:"column:addresses"`
	Fonts         []byte `gorm:"column:fonts"`
	Papers        []byte `gorm:"column:papers"`
	Signets       []byte `gorm:"column:signets"`
	FontColors    []byte `gorm:"column:font_colors"`
	FunctionCards []byte `gorm:"column:function_cards"`

	Remark string `gorm:"column:remark"`
}

func (userPO) TableName() string { return "user" }

// avatarPO 映射 avatar 表。
//
// 头像字典表，无审计字段。
type avatarPO struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:name"`
	URL  string `gorm:"column:url"`
}

func (avatarPO) TableName() string { return "avatar" }

// TODO: 实现 userPO 与 user.User 的互转，包含六个 JSON 列的解析与
// Status / DelFlag 的字符串转换。
