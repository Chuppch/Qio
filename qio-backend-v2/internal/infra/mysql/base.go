package mysql

import "time"

// 审计字段。
//
// v1 的 BaseEntity 让全部 20 个实体继承同一套字段，包括那些表里并没有这些列的
// 实体。这里改为按需嵌入，并且因为 v1 的 schema 在不同表中用了不同的
// create_user / update_user 类型，拆成两个变体如实映射。
//
// 两个变体的存在本身就是 schema 不一致的证据，后续统一为 bigint 之后应合并为一个。

// auditFields 适用于 create_user / update_user 为 bigint 的表：
// letter、user、friend、friend_request。
type auditFields struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement"`
	CreateUser int64     `gorm:"column:create_user"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateUser int64     `gorm:"column:update_user"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

// auditFieldsStrUser 适用于 create_user / update_user 为 varchar(50) 的表：
// bottle、questions、question_user_status。
//
// 存的仍然是用户 ID，只是列类型被定义成了字符串。
type auditFieldsStrUser struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement"`
	CreateUser string    `gorm:"column:create_user"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateUser string    `gorm:"column:update_user"`
	UpdateTime time.Time `gorm:"column:update_time"`
}
