package mysql

import (
	"encoding/json"
	"strconv"

	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/bottle"
)

// bottlePO 映射 bottle 表。
//
// 本表的 create_user / update_user 是 varchar(50) 而非 bigint，因此嵌入
// auditFieldsStrUser。v1 复用 update_user 记录「谁捞的」，转换时映射为
// 领域模型的 PickedBy。
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

func (p bottlePO) toDomain() *bottle.Bottle {
	addr := decodeAddress(p.SenderAddress)

	return &bottle.Bottle{
		ID:       p.ID,
		UserID:   p.UserID,
		Nickname: p.NickName,
		Email:    p.Email,
		SenderAddr: bottle.Address{
			ID:               addr.ID,
			CountryID:        addr.CountryID,
			FormattedAddress: addr.FormattedAddress,
			Longitude:        addr.Longitude,
			Latitude:         addr.Latitude,
			IsDefault:        addr.defaultFlag(),
		},
		Content:   p.Content,
		ImageURL:  p.BottleURL,
		Picked:    p.Picked,
		PickedBy:  atoiDefault64(p.UpdateUser, 0),
		Remark:    p.Remark,
		CreatedAt: p.CreateTime,
		UpdatedAt: p.UpdateTime,
	}
}

func bottlePOFrom(b *bottle.Bottle) (bottlePO, error) {
	addr, err := json.Marshal(addressJSON{
		ID:               b.SenderAddr.ID,
		CountryID:        b.SenderAddr.CountryID,
		FormattedAddress: b.SenderAddr.FormattedAddress,
		Longitude:        b.SenderAddr.Longitude,
		Latitude:         b.SenderAddr.Latitude,
		IsDefault:        strconv.FormatBool(b.SenderAddr.IsDefault),
	})
	if err != nil {
		return bottlePO{}, err
	}

	updateUser := ""
	if b.PickedBy != 0 {
		updateUser = strconv.FormatInt(b.PickedBy, 10)
	}

	return bottlePO{
		auditFieldsStrUser: auditFieldsStrUser{
			ID:         b.ID,
			CreateUser: strconv.FormatInt(b.UserID, 10),
			CreateTime: b.CreatedAt,
			UpdateUser: updateUser,
			UpdateTime: b.UpdatedAt,
		},
		UserID:        b.UserID,
		NickName:      b.Nickname,
		Email:         b.Email,
		SenderAddress: addr,
		Content:       b.Content,
		Picked:        b.Picked,
		BottleURL:     b.ImageURL,
		Remark:        b.Remark,
	}, nil
}

// atoiDefault64 解析 varchar 列中的用户 ID。
//
// bottle、questions、question_user_status 三张表把 create_user / update_user
// 定义为 varchar，实际存的是数字 ID；空值与脏数据按兜底值处理。
func atoiDefault64(s string, fallback int64) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fallback
	}
	return v
}
