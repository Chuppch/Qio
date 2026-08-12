package mysql

// bottlePO 映射 bottle 表。
//
// 本表的 create_user / update_user 是 varchar(50) 而非 bigint，因此嵌入
// auditFieldsStrUser。
type bottlePO struct {
	auditFieldsStrUser

	UserID   int64  `gorm:"column:user_id"`
	NickName string `gorm:"column:nick_name"`
	Email    string `gorm:"column:email"`

	// SenderAddress 是 JSON 列，存单个 addressJSON。
	SenderAddress []byte `gorm:"column:sender_address"`

	Content string `gorm:"column:content"`

	// Picked 对应 is_picked，表中类型为 tinyint(1)。
	Picked bool `gorm:"column:is_picked"`

	BottleURL string `gorm:"column:bottle_url"`

	Remark string `gorm:"column:remark"`
}

func (bottlePO) TableName() string { return "bottle" }

// TODO: 实现与 bottle.Bottle 的互转。
