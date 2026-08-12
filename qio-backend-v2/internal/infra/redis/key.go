package redis

import "strconv"

// 缓存键构造。
//
// 迁移自 v1 的 CacheConstant，但改变了使用方式：v1 导出键前缀常量，由各个
// Service 自行 `PREFIX + userId` 拼接，键的正确性依赖调用方自觉，且一处改名
// 要全局搜索。这里改为不导出的前缀 + 导出的构造函数，业务侧拿不到裸前缀，
// 只能通过函数生成键，拼错的可能性由编译器消除。
//
// 键的命名体系沿用 v1，避免灰度期间新旧版本读到不同的键。

const (
	prefixUserInfo          = "cache:user:user-info:"
	prefixUserFunctionCards = "cache:user:function-cards:"
	prefixUserFriends       = "cache:user:friends:"
	prefixUserAddresses     = "cache:user:addresses:"
	prefixUserRepository    = "cache:user:repository:"
	prefixUserWriteLetter   = "cache:user:write-letter:"
	prefixUserReceiveLetter = "cache:user:receive-letter:"

	keyAvatarList  = "cache:avatars:list"
	keyCountryList = "cache:countries:list"
	keyWordLimit   = "cache:word:limit"

	keyShopFont         = "cache:shop:font"
	keyShopFontColor    = "cache:shop:font-color"
	keyShopPaper        = "cache:shop:paper"
	keyShopFunctionCard = "cache:shop:function-card"

	prefixSignToday  = "sign:today:"
	prefixSignAward  = "sign:award:"
	keySignCurrent   = "sign:current"
	prefixSignSigned = "sign:signed:"
	infixSignUser    = ":user-"

	prefixGameFortuneUser = "game:ffl:user:"

	keyTaskDetails = "task:taskDetails"
)

// ---- 用户维度 ----

// UserInfoKey 用户基本信息缓存键。
func UserInfoKey(userID int64) string { return prefixUserInfo + itoa(userID) }

// UserFunctionCardsKey 用户功能卡缓存键。
func UserFunctionCardsKey(userID int64) string { return prefixUserFunctionCards + itoa(userID) }

// UserFriendsKey 用户好友列表缓存键。
func UserFriendsKey(userID int64) string { return prefixUserFriends + itoa(userID) }

// UserAddressesKey 用户地址簿缓存键。
func UserAddressesKey(userID int64) string { return prefixUserAddresses + itoa(userID) }

// UserRepositoryKey 用户仓库缓存键。
func UserRepositoryKey(userID int64) string { return prefixUserRepository + itoa(userID) }

// UserWriteLetterKey 用户已寄出信件列表缓存键。
func UserWriteLetterKey(userID int64) string { return prefixUserWriteLetter + itoa(userID) }

// UserReceiveLetterKey 用户收到信件列表缓存键。
func UserReceiveLetterKey(userID int64) string { return prefixUserReceiveLetter + itoa(userID) }

// ---- 全局字典 ----

// AvatarListKey 头像列表缓存键。
func AvatarListKey() string { return keyAvatarList }

// CountryListKey 国家列表缓存键。
func CountryListKey() string { return keyCountryList }

// WordLimitKey 敏感词限制缓存键。
func WordLimitKey() string { return keyWordLimit }

// ---- 商城 ----

// ShopFontKey 字体商品缓存键。
func ShopFontKey() string { return keyShopFont }

// ShopFontColorKey 字体颜色商品缓存键。
func ShopFontColorKey() string { return keyShopFontColor }

// ShopPaperKey 信纸商品缓存键。
func ShopPaperKey() string { return keyShopPaper }

// ShopFunctionCardKey 功能卡商品缓存键。
func ShopFunctionCardKey() string { return keyShopFunctionCard }

// ---- 签到 ----

// SignTodayKey 指定日期的签到汇总缓存键，date 形如 20260811。
func SignTodayKey(date string) string { return prefixSignToday + date }

// SignAwardKey 指定日期的签到奖励缓存键。
func SignAwardKey(date string) string { return prefixSignAward + date }

// SignCurrentKey 当前签到周期缓存键。
func SignCurrentKey() string { return keySignCurrent }

// SignedKey 用户在指定日期的签到标记缓存键。
func SignedKey(date string, userID int64) string {
	return prefixSignSigned + date + infixSignUser + itoa(userID)
}

// ---- 玩法 ----

// GameFortuneUserKey 用户抽奖次数缓存键。
func GameFortuneUserKey(userID int64) string { return prefixGameFortuneUser + itoa(userID) }

// TaskDetailsKey 任务配置缓存键。
func TaskDetailsKey() string { return keyTaskDetails }

func itoa(v int64) string { return strconv.FormatInt(v, 10) }
