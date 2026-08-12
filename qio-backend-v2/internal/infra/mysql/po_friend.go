package mysql

// friendPO 映射 friend 表。
//
// 单向存储：OwningID 为关系归属者，UserID 为对方。互为好友时存在两条记录。
//
// 注意：friend.id 在表中是 int 而非 bigint，与 user.id 的 bigint 不一致。
// auditFields 统一用 int64 承载，读写不受影响。
type friendPO struct {
	auditFields

	UserID   int64  `gorm:"column:user_id"`
	OwningID int64  `gorm:"column:owning_id"`
	Name     string `gorm:"column:name"`
	Sex      string `gorm:"column:sex"`
	Email    string `gorm:"column:email"`
	Avatar   string `gorm:"column:avatar"`

	// Addresses 是 JSON 列，存 addressJSON 数组。
	Addresses []byte `gorm:"column:addresses"`

	Remark string `gorm:"column:remark"`
}

func (friendPO) TableName() string { return "friend" }

// friendRequestPO 映射 friend_request 表。
type friendRequestPO struct {
	auditFields

	SenderID   int64 `gorm:"column:sender_id"`
	ReceiverID int64 `gorm:"column:receiver_id"`
	Status     int   `gorm:"column:status"`

	// GiveAddress 是 JSON 列，存单个 addressJSON。
	GiveAddress []byte `gorm:"column:give_address"`

	Content string `gorm:"column:content"`

	// BottleID、LetterID 记录申请来源，二者互斥，均为 0 表示主动添加。
	BottleID int64 `gorm:"column:bottle_id"`
	LetterID int64 `gorm:"column:letter_id"`

	Remark string `gorm:"column:remark"`
}

func (friendRequestPO) TableName() string { return "friend_request" }

// TODO: 实现与 friend.Friend、friend.Request 的互转。
