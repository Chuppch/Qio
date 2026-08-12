package mysql

import "time"

// letterPO 映射 letter 表。
type letterPO struct {
	auditFields

	SenderUserID    int64  `gorm:"column:sender_user_id"`
	SenderName      string `gorm:"column:sender_name"`
	SenderEmail     string `gorm:"column:sender_email"`
	RecipientEmail  string `gorm:"column:recipient_email"`
	RecipientUserID int64  `gorm:"column:recipient_user_id"`
	RecipientName   string `gorm:"column:recipient_name"`

	LetterContent string `gorm:"column:letter_content"`
	LetterLink    string `gorm:"column:letter_link"`
	CoverLink     string `gorm:"column:cover_link"`

	// SenderAddress、RecipientAddress 是 JSON 列，结构见 addressJSON。
	SenderAddress    []byte `gorm:"column:sender_address"`
	RecipientAddress []byte `gorm:"column:recipient_address"`

	ExpectedDeliveryTime time.Time `gorm:"column:expected_delivery_time"`
	DeliveryTime         time.Time `gorm:"column:delivery_time"`

	Status           int   `gorm:"column:status"`
	DeliveryProgress int64 `gorm:"column:delivery_progress"`
	ReadStatus       int   `gorm:"column:read_status"`
	PiggyMoney       int64 `gorm:"column:piggy_money"`

	// LetterType 取值 1 为竖版侨批、2 为横版普通信件，与表注释一致。
	LetterType int `gorm:"column:letter_type"`

	// SpeedRate、ReduceTime 在表中是 varchar，分别存倍率与分钟数。
	// 转换为领域模型时解析为 float64 与 time.Duration。
	SpeedRate  string `gorm:"column:speed_rate"`
	ReduceTime string `gorm:"column:reduce_time"`

	Remark string `gorm:"column:remark"`
}

// TableName 指定表名，避免 GORM 默认的复数化推断。
func (letterPO) TableName() string { return "letter" }

// TODO: 实现 toDomain 与 letterPOFrom，包含 JSON 地址解析与
// SpeedRate / ReduceTime 的字符串转换。
