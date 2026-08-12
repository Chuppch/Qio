package mysql

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/user"
)

// userPO 映射 user 表。
type userPO struct {
	auditFields

	Username string `gorm:"column:username"`
	Nickname string `gorm:"column:nickname"`
	Email    string `gorm:"column:email"`
	Sex      string `gorm:"column:sex"`
	Avatar   string `gorm:"column:avatar"`

	// Password 列存的是密码摘要
	Password string `gorm:"column:password"`

	// Status、DelFlag 在表中是 varchar，取值 "0" / "1"（DelFlag 删除标记为 "2"）
	Status  string `gorm:"column:status"`
	DelFlag string `gorm:"column:del_flag"`

	LoginIP   string    `gorm:"column:login_ip"`
	LoginDate time.Time `gorm:"column:login_date"`

	Money int64 `gorm:"column:money"`

	// 以下六列均为 JSON 数组。
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

// delFlagDeleted 是 del_flag 列表示已删除的取值，沿用 v1 约定。
const delFlagDeleted = "2"

func (p userPO) toDomain() *user.User {
	u := &user.User{
		ID:           p.ID,
		Username:     p.Username,
		Nickname:     p.Nickname,
		Email:        p.Email,
		Sex:          p.Sex,
		Avatar:       p.Avatar,
		PasswordHash: p.Password,
		Status:       user.AccountStatus(atoiDefault(p.Status, 0)),
		Deleted:      p.DelFlag == delFlagDeleted,
		Money:        p.Money,
		LastLoginIP:  p.LoginIP,
		LastLoginAt:  p.LoginDate,
		CreatedAt:    p.CreateTime,
		UpdatedAt:    p.UpdateTime,
	}

	u.Addresses = decodeUserAddresses(p.Addresses)
	u.OwnedItems = p.decodeOwnedItems()
	return u
}

// decodeUserAddresses 解析 addresses JSON 列。
//
// 解析失败时返回空切片而不报错：v1 存量数据中该列可能为 NULL 或空字符串，
// 视作「尚未填写地址」。
func decodeUserAddresses(raw []byte) []user.Address {
	items := decodeAddressList(raw)
	if len(items) == 0 {
		return nil
	}

	out := make([]user.Address, 0, len(items))
	for _, a := range items {
		out = append(out, user.Address{
			ID:               a.ID,
			CountryID:        a.CountryID,
			FormattedAddress: a.FormattedAddress,
			Longitude:        a.Longitude,
			Latitude:         a.Latitude,
			IsDefault:        a.defaultFlag(),
		})
	}
	return out
}

// decodeOwnedItems 把五个道具 JSON 列合并为统一的背包表示。
//
// v1 中每列存的是各自的 VO 数组，字段不完全一致，此处只取 id——背包只需要
// 「拥有哪件道具」，名称与预览图在展示时由 shop 域按 id 查询。
func (p userPO) decodeOwnedItems() []user.OwnedItem {
	sources := []struct {
		raw  []byte
		kind user.ItemType
	}{
		{p.FunctionCards, user.ItemFunctionCard},
		{p.Fonts, user.ItemFont},
		{p.FontColors, user.ItemFontColor},
		{p.Papers, user.ItemPaper},
		{p.Signets, user.ItemCollection},
	}

	var out []user.OwnedItem
	for _, s := range sources {
		for _, it := range decodeOwnedItemList(s.raw) {
			out = append(out, user.OwnedItem{Type: s.kind, ItemID: it.ID, Count: 1})
		}
	}
	return out
}

// userPOFrom 把领域模型转为 PO。
//
// 六个 JSON 列按 v1 的存储格式回写，保证新旧版本可读同一份数据。
func userPOFrom(u *user.User) (userPO, error) {
	addresses, err := encodeUserAddresses(u.Addresses)
	if err != nil {
		return userPO{}, err
	}

	grouped := groupOwnedItems(u.OwnedItems)

	cards, err := encodeOwnedItems(grouped[user.ItemFunctionCard])
	if err != nil {
		return userPO{}, err
	}
	fonts, err := encodeOwnedItems(grouped[user.ItemFont])
	if err != nil {
		return userPO{}, err
	}
	colors, err := encodeOwnedItems(grouped[user.ItemFontColor])
	if err != nil {
		return userPO{}, err
	}
	papers, err := encodeOwnedItems(grouped[user.ItemPaper])
	if err != nil {
		return userPO{}, err
	}
	signets, err := encodeOwnedItems(grouped[user.ItemCollection])
	if err != nil {
		return userPO{}, err
	}

	delFlag := "0"
	if u.Deleted {
		delFlag = delFlagDeleted
	}

	return userPO{
		auditFields: auditFields{
			ID:         u.ID,
			CreateTime: u.CreatedAt,
			UpdateTime: u.UpdatedAt,
		},
		Username:      u.Username,
		Nickname:      u.Nickname,
		Email:         u.Email,
		Sex:           u.Sex,
		Avatar:        u.Avatar,
		Password:      u.PasswordHash,
		Status:        strconv.Itoa(int(u.Status)),
		DelFlag:       delFlag,
		LoginIP:       u.LastLoginIP,
		LoginDate:     u.LastLoginAt,
		Money:         u.Money,
		Addresses:     addresses,
		Fonts:         fonts,
		Papers:        papers,
		Signets:       signets,
		FontColors:    colors,
		FunctionCards: cards,
	}, nil
}

func encodeUserAddresses(list []user.Address) ([]byte, error) {
	items := make([]addressJSON, 0, len(list))
	for _, a := range list {
		items = append(items, addressJSON{
			ID:               a.ID,
			CountryID:        a.CountryID,
			FormattedAddress: a.FormattedAddress,
			Longitude:        a.Longitude,
			Latitude:         a.Latitude,
			IsDefault:        strconv.FormatBool(a.IsDefault),
		})
	}
	return json.Marshal(items)
}

func groupOwnedItems(items []user.OwnedItem) map[user.ItemType][]user.OwnedItem {
	out := make(map[user.ItemType][]user.OwnedItem)
	for _, it := range items {
		out[it.Type] = append(out[it.Type], it)
	}
	return out
}

func encodeOwnedItems(items []user.OwnedItem) ([]byte, error) {
	out := make([]ownedItemJSON, 0, len(items))
	for _, it := range items {
		out = append(out, ownedItemJSON{ID: it.ItemID})
	}
	return json.Marshal(out)
}

func atoiDefault(s string, fallback int) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return n
}
