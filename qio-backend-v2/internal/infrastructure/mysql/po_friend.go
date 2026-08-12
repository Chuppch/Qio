package mysql

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/friend"
)

// friendAuditFields 映射 friend 与 friend_request 中允许为 NULL 的审计列。
//
// 这两张 v1 表除主键外的审计字段均未声明 NOT NULL。使用指针承接数据库 NULL，
// 再在 PO/Domain 边界转换为领域零值，避免 database/sql 把 NULL 扫描到标量时报错。
type friendAuditFields struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement"`
	CreateUser *int64     `gorm:"column:create_user"`
	CreateTime *time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateUser *int64     `gorm:"column:update_user"`
	UpdateTime *time.Time `gorm:"column:update_time;autoUpdateTime"`
}

// friendPO 映射 friend 表。
//
// 单向存储：OwningID 为关系归属者，UserID 为对方。互为好友时存在两条记录。
//
// 注意：friend.id 在表中是 int 而非 bigint，与 user.id 的 bigint 不一致。
// auditFields 统一用 int64 承载，读写不受影响。
type friendPO struct {
	friendAuditFields

	UserID   *int64  `gorm:"column:user_id"`
	OwningID *int64  `gorm:"column:owning_id"`
	Name     *string `gorm:"column:name"`
	Sex      *string `gorm:"column:sex"`
	Email    *string `gorm:"column:email"`
	Avatar   *string `gorm:"column:avatar"`

	// Addresses 是 JSON 列，存 addressJSON 数组。
	Addresses []byte `gorm:"column:addresses"`

	Remark *string `gorm:"column:remark"`
}

func (friendPO) TableName() string { return "friend" }

func (p friendPO) toDomain() *friend.Friend {
	return &friend.Friend{
		ID:        p.ID,
		UserID:    valueOrZero(p.UserID),
		OwnerID:   valueOrZero(p.OwningID),
		Name:      valueOrZero(p.Name),
		Sex:       valueOrZero(p.Sex),
		Email:     valueOrZero(p.Email),
		Avatar:    valueOrZero(p.Avatar),
		Addresses: friendAddressesToDomain(p.Addresses),
		Remark:    valueOrZero(p.Remark),
		CreatedAt: valueOrZero(p.CreateTime),
		UpdatedAt: valueOrZero(p.UpdateTime),
	}
}

// friendPOFrom 把领域好友记录转为持久化对象。
//
// create_user / update_user 在本表是 bigint，v1 由 MyBatis-Plus 自动填充为当前
// 用户 ID。这里统一填归属者：好友记录的增删改都由归属者本人触发。
func friendPOFrom(f *friend.Friend) (friendPO, error) {
	addresses, err := encodeFriendAddresses(f.Addresses)
	if err != nil {
		return friendPO{}, err
	}

	return friendPO{
		friendAuditFields: friendAuditFields{
			ID:         f.ID,
			CreateUser: optionalInt64(f.OwnerID),
			CreateTime: optionalTime(f.CreatedAt),
			UpdateUser: optionalInt64(f.OwnerID),
			UpdateTime: optionalTime(f.UpdatedAt),
		},
		UserID:    optionalInt64(f.UserID),
		OwningID:  optionalInt64(f.OwnerID),
		Name:      pointerTo(f.Name),
		Sex:       pointerTo(f.Sex),
		Email:     pointerTo(f.Email),
		Avatar:    pointerTo(f.Avatar),
		Addresses: addresses,
		Remark:    pointerTo(f.Remark),
	}, nil
}

// friendRequestPO 映射 friend_request 表。
type friendRequestPO struct {
	friendAuditFields

	SenderID   *int64 `gorm:"column:sender_id"`
	ReceiverID *int64 `gorm:"column:receiver_id"`
	Status     *int   `gorm:"column:status"`

	// GiveAddress 是 JSON 列，存单个 addressJSON。
	GiveAddress []byte `gorm:"column:give_address"`

	Content *string `gorm:"column:content"`

	// BottleID、LetterID 记录申请来源，二者互斥，均为 0 表示无来源。
	BottleID *int64 `gorm:"column:bottle_id"`
	LetterID *int64 `gorm:"column:letter_id"`

	Remark *string `gorm:"column:remark"`
}

func (friendRequestPO) TableName() string { return "friend_request" }

func (p friendRequestPO) toDomain() *friend.Request {
	return &friend.Request{
		ID:          p.ID,
		SenderID:    valueOrZero(p.SenderID),
		ReceiverID:  valueOrZero(p.ReceiverID),
		Status:      friend.RequestStatus(valueOrZero(p.Status)),
		GiveAddress: friendAddressToDomain(decodeAddress(p.GiveAddress)),
		Content:     valueOrZero(p.Content),
		BottleID:    valueOrZero(p.BottleID),
		LetterID:    valueOrZero(p.LetterID),
		Remark:      valueOrZero(p.Remark),
		CreatedAt:   valueOrZero(p.CreateTime),
		UpdatedAt:   valueOrZero(p.UpdateTime),
	}
}

// friendRequestPOFrom 把领域好友申请转为持久化对象。
//
// 审计字段填申请人：v1 的自动填充取当前登录用户，而申请的创建与状态变更分别
// 由申请人与接收人触发，这里统一记申请人，与 v1 在多数场景下的实际落库值一致。
func friendRequestPOFrom(r *friend.Request) (friendRequestPO, error) {
	addr, err := json.Marshal(friendAddressToJSON(r.GiveAddress))
	if err != nil {
		return friendRequestPO{}, err
	}

	return friendRequestPO{
		friendAuditFields: friendAuditFields{
			ID:         r.ID,
			CreateUser: optionalInt64(r.SenderID),
			CreateTime: optionalTime(r.CreatedAt),
			UpdateUser: optionalInt64(r.SenderID),
			UpdateTime: optionalTime(r.UpdatedAt),
		},
		SenderID:    optionalInt64(r.SenderID),
		ReceiverID:  optionalInt64(r.ReceiverID),
		Status:      pointerTo(int(r.Status)),
		GiveAddress: addr,
		Content:     pointerTo(r.Content),
		BottleID:    optionalInt64(r.BottleID),
		LetterID:    optionalInt64(r.LetterID),
		Remark:      pointerTo(r.Remark),
	}, nil
}

func pointerTo[T any](value T) *T { return &value }

func valueOrZero[T any](value *T) T {
	if value == nil {
		var zero T
		return zero
	}
	return *value
}

func optionalInt64(value int64) *int64 {
	if value == 0 {
		return nil
	}
	return &value
}

func optionalTime(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}
	return &value
}

// ---- 地址列互转 ----

func friendAddressToDomain(a addressJSON) friend.Address {
	return friend.Address{
		ID:               a.ID,
		CountryID:        a.CountryID,
		FormattedAddress: a.FormattedAddress,
		Longitude:        a.Longitude,
		Latitude:         a.Latitude,
		IsDefault:        a.defaultFlag(),
	}
}

func friendAddressToJSON(a friend.Address) addressJSON {
	return addressJSON{
		ID:               a.ID,
		CountryID:        a.CountryID,
		FormattedAddress: a.FormattedAddress,
		Longitude:        a.Longitude,
		Latitude:         a.Latitude,
		IsDefault:        strconv.FormatBool(a.IsDefault),
	}
}

func friendAddressesToDomain(raw []byte) []friend.Address {
	items := decodeAddressList(raw)
	if len(items) == 0 {
		return nil
	}
	out := make([]friend.Address, 0, len(items))
	for _, a := range items {
		out = append(out, friendAddressToDomain(a))
	}
	return out
}

// encodeFriendAddresses 序列化地址簿。
//
// 地址簿为空时写入 JSON 空数组而非 NULL：v1 的 Jackson 处理器对空集合同样
// 序列化为 []，保持存量格式一致。
func encodeFriendAddresses(addresses []friend.Address) ([]byte, error) {
	items := make([]addressJSON, 0, len(addresses))
	for _, a := range addresses {
		items = append(items, friendAddressToJSON(a))
	}
	return json.Marshal(items)
}
