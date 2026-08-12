package mysql

// addressJSON 是地址在 JSON 列中的存储结构。
//
// v1 把地址以 JSON 形式存在 letter.sender_address、letter.recipient_address、
// user.addresses、friend.addresses、bottle.sender_address、
// friend_request.give_address 六处。字段名沿用 v1 序列化时使用的名称，
// 不可随意改动，否则读不出存量数据。
//
// v1 直接把 com.qiaopi.entity.Address（以及 vo 包下的道具 VO）序列化进 JSON 列，
// 导致 Java 类的字段名成了存储格式的一部分。这里把存储结构显式定义在 infra 内，
// 与领域模型解耦：领域模型改字段名不影响存量数据。
type addressJSON struct {
	ID               int64   `json:"id"`
	CountryID        int64   `json:"countryId"`
	FormattedAddress string  `json:"formattedAddress"`
	Longitude        float64 `json:"longitude"`
	Latitude         float64 `json:"latitude"`
	// IsDefault 在 v1 中被存成字符串 "true" / "false"，不是布尔值。
	IsDefault string `json:"isDefault"`
}

// ownedItemJSON 是用户背包中一件道具在 JSON 列中的存储结构。
//
// 对应 user.fonts、user.papers、user.signets、user.font_colors、
// user.function_cards 五个 JSON 列。v1 在这些列里存的是各自的 VO 对象，
// 字段随 VO 定义而不同，此处取其公共部分。
//
// 注意：这些列的实际内容因 VO 而异，接入存量数据前需要逐列核对真实 JSON 结构。
type ownedItemJSON struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	PreviewImage string `json:"previewImage"`
}
