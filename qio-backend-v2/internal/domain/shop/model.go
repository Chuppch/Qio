package shop

// ItemType 是商品类型。
//
// 取值与 user.ItemType 一致，使购买、签到奖励、背包三处共用同一套编码。
// 两个域各自持有一份定义而不共享类型，避免 shop 与 user 互相依赖。
type ItemType int

const (
	ItemFunctionCard ItemType = 2 // 功能卡
	ItemFont         ItemType = 3 // 字体
	ItemFontColor    ItemType = 4 // 字体颜色
	ItemPaper        ItemType = 5 // 信纸
	ItemSignet       ItemType = 6 // 印章，v1 中归入「其他收藏品」
)

// Valid 表示取值是可购买的商品类型。
func (t ItemType) Valid() bool {
	switch t {
	case ItemFunctionCard, ItemFont, ItemFontColor, ItemPaper, ItemSignet:
		return true
	default:
		return false
	}
}

// Item 是商品的公共属性。
//
// v1 把信纸、字体、字体颜色、印章、功能卡各建一张表，结构高度相似。这里抽出
// 公共部分，各类型特有的属性由专门类型嵌入 Item 承载。
type Item struct {
	ID           int64
	Type         ItemType
	Name         string
	Description  string
	PreviewImage string
	Price        int
}

// Font 是字体商品。
type Font struct {
	Item
	// FilePath 是字体文件在对象存储中的路径，渲染信件时加载
	FilePath string
}

// FontColor 是字体颜色商品。
//
// HexCode 在库中有唯一约束。
type FontColor struct {
	Item
	HexCode  string
	RGBValue string
}

// Signet 是印章商品。
//
// v1 中印章只在注册时赠送，没有购买接口。
type Signet struct {
	Item
	FilePath string
}

// Paper 是信纸商品。
//
// 除商品属性外还携带整套排版参数，供信件渲染时定位文本。v1 把这些参数存成
// varchar，实际是数值，这里改为强类型，字符串解析由仓储层负责。
type Paper struct {
	Item

	// FilePath 是信纸底图在对象存储中的路径
	FilePath string

	// Layout 决定该信纸适用的款式：1 侨批，2 普通信纸。取值与 letter.Layout 一致。
	Layout int

	Typography Typography
}

// Typography 是信纸上的文本排版参数，单位为渲染画布像素。
type Typography struct {
	FontSize float64

	// 正文起始偏移
	TranslateX float64
	TranslateY float64

	// 收件人姓名偏移
	RecipientTranslateX float64
	RecipientTranslateY float64

	// 寄件人姓名偏移
	SenderTranslateX float64
	SenderTranslateY float64
}

// CardType 是功能卡的效果类型。
type CardType int

const (
	CardAccelerate CardType = 1 // 加速卡，提升投递倍率
	CardReduceTime CardType = 2 // 减时卡，直接减免投递时长
)

// FunctionCard 是功能卡商品。
//
// 功能卡作用于信件投递：ReduceMinutes 直接减免投递时长，SpeedRate 提升投递倍率。
// v1 把这两个值存成 varchar，这里改为数值类型。
type FunctionCard struct {
	Item

	CardType CardType

	// Enabled 对应 v1 的 card_status
	Enabled bool

	// ReduceMinutes 是可减免的投递时长，单位分钟
	ReduceMinutes int
	// SpeedRate 是投递加速倍率，1 表示不加速
	SpeedRate float64

	Remark string
}

// PaperFontFit 是信纸与字体组合下的适配字数。
//
// 迁移自 v1 的 font_paper 表，写信时用它校验正文长度是否超出信纸容量。
type PaperFontFit struct {
	PaperID   int64
	FontID    int64
	FitNumber int
}

// Commodity 是运营位展示的外部文创商品。
//
// 不参与站内购买流程，只做展示与外链跳转，因此不复用 Item。Price 保持字符串，
// 展示的是价格文案而非可运算的金额。
type Commodity struct {
	ID          int64
	Name        string
	Description string
	Price       string
	Image       string
	Marketing   string
	Link        string
}
