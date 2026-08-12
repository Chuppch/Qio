package mysql

import "encoding/json"

// JSON 列的存储结构。
//
// v1 把地址与用户背包以 JSON 形式存在 letter、user、friend、bottle、
// friend_request 多张表中，序列化的直接是 Java 的 entity 与 vo 类，因此这些
// 字段名成了存储格式的一部分，不可随意改动。
//
// 把存储结构显式定义在 infra 内，与领域模型解耦：领域模型改字段名或改类型
// 都不影响存量数据的读取。

// addressJSON 是地址在 JSON 列中的存储结构。
type addressJSON struct {
	ID               int64   `json:"id"`
	CountryID        int64   `json:"countryId"`
	FormattedAddress string  `json:"formattedAddress"`
	Longitude        float64 `json:"longitude"`
	Latitude         float64 `json:"latitude"`
	// IsDefault 在 v1 中被存成字符串 "true" / "false"，不是布尔值
	IsDefault string `json:"isDefault"`
}

// defaultFlag 把 v1 的字符串标记转为布尔值。
func (a addressJSON) defaultFlag() bool { return a.IsDefault == "true" }

// ownedItemJSON 是用户背包中一件道具在 JSON 列中的存储结构。
//
// 对应 user 表的 fonts、papers、signets、font_colors、function_cards 五列。
// v1 在这些列里存的是各自的 VO 对象，字段随 VO 定义而异，此处只取共有的 id
// 与展示字段。
//
// 注意：这五列的实际内容因 VO 而异，接入存量数据前需逐列核对真实 JSON 结构。
type ownedItemJSON struct {
	ID           int64  `json:"id"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	PreviewImage string `json:"previewImage,omitempty"`
}

// decodeAddressList 解析地址数组列。
//
// 列为 NULL、空字符串或格式非法时返回 nil，视作「尚未填写」，不作为错误处理——
// v1 存量数据中这些情况普遍存在。
func decodeAddressList(raw []byte) []addressJSON {
	if len(raw) == 0 {
		return nil
	}
	var out []addressJSON
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil
	}
	return out
}

// decodeAddress 解析单个地址列，用于 bottle.sender_address 与
// friend_request.give_address 这类存单对象的列。
//
// 列为 NULL 或格式非法时返回零值，视作「未填写地址」。
func decodeAddress(raw []byte) addressJSON {
	if len(raw) == 0 {
		return addressJSON{}
	}
	var out addressJSON
	if err := json.Unmarshal(raw, &out); err != nil {
		return addressJSON{}
	}
	return out
}

// decodeOwnedItemList 解析道具数组列，语义同 decodeAddressList。
func decodeOwnedItemList(raw []byte) []ownedItemJSON {
	if len(raw) == 0 {
		return nil
	}
	var out []ownedItemJSON
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil
	}
	return out
}
