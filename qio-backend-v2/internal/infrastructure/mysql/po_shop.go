package mysql

import (
	"strconv"

	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/shop"
)

// 商城相关表。
//
// v1 把信纸、字体、字体颜色、印章、功能卡各自建表，结构高度相似（名称、
// 预览图、价格）。v2 在领域层合并为 shop 域，但 PO 层如实保留分表现状，
// 表合并需要单独的数据迁移任务。

// paperPO 映射 paper 表。
//
// 注意：paper 表没有主键，且 id 列可为 NULL。这是 v1 schema 的缺陷，
// GORM 的更新与删除操作依赖主键；本域只读，暂不受影响。
type paperPO struct {
	ID           int64  `gorm:"column:id;primaryKey"`
	Name         string `gorm:"column:name"`
	PreviewImage string `gorm:"column:preview_image"`
	FilePath     string `gorm:"column:file_path"`

	// 以下排版参数在表中均为 varchar，实际存的是数值，
	// 供信件渲染时定位文本位置。
	FontSize            string `gorm:"column:font_size"`
	TranslateX          string `gorm:"column:translate_x"`
	TranslateY          string `gorm:"column:translate_y"`
	RecipientTranslateX string `gorm:"column:recipient_translate_x"`
	RecipientTranslateY string `gorm:"column:recipient_translate_y"`
	SenderTranslateX    string `gorm:"column:sender_translate_x"`
	SenderTranslateY    string `gorm:"column:sender_translate_y"`

	Price int `gorm:"column:price"`

	// Type 取值与 letter.letter_type 一致：1 为侨批，2 为普通信纸。
	Type int `gorm:"column:type"`
}

func (paperPO) TableName() string { return "paper" }

func (p paperPO) toDomain() shop.Paper {
	return shop.Paper{
		Item: shop.Item{
			ID:           p.ID,
			Type:         shop.ItemPaper,
			Name:         p.Name,
			PreviewImage: p.PreviewImage,
			Price:        p.Price,
		},
		FilePath: p.FilePath,
		Layout:   p.Type,
		Typography: shop.Typography{
			FontSize:            parseFloatDefault(p.FontSize, 0),
			TranslateX:          parseFloatDefault(p.TranslateX, 0),
			TranslateY:          parseFloatDefault(p.TranslateY, 0),
			RecipientTranslateX: parseFloatDefault(p.RecipientTranslateX, 0),
			RecipientTranslateY: parseFloatDefault(p.RecipientTranslateY, 0),
			SenderTranslateX:    parseFloatDefault(p.SenderTranslateX, 0),
			SenderTranslateY:    parseFloatDefault(p.SenderTranslateY, 0),
		},
	}
}

// fontPO 映射 font 表。
type fontPO struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Name         string `gorm:"column:name"`
	PreviewImage string `gorm:"column:preview_image"`
	FilePath     string `gorm:"column:file_path"`
	Price        int    `gorm:"column:price"`
}

func (fontPO) TableName() string { return "font" }

func (p fontPO) toDomain() shop.Font {
	return shop.Font{
		Item: shop.Item{
			ID:           p.ID,
			Type:         shop.ItemFont,
			Name:         p.Name,
			PreviewImage: p.PreviewImage,
			Price:        p.Price,
		},
		FilePath: p.FilePath,
	}
}

// fontColorPO 映射 font_color 表。
//
// hex_code 上有唯一约束。表中没有 name 列，展示名称取自 description。
type fontColorPO struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement"`
	HexCode      string `gorm:"column:hex_code"`
	RGBValue     string `gorm:"column:rgb_value"`
	Description  string `gorm:"column:description"`
	PreviewImage string `gorm:"column:preview_image"`
	Price        int    `gorm:"column:price"`
}

func (fontColorPO) TableName() string { return "font_color" }

func (p fontColorPO) toDomain() shop.FontColor {
	return shop.FontColor{
		Item: shop.Item{
			ID:           p.ID,
			Type:         shop.ItemFontColor,
			Name:         p.Description,
			Description:  p.Description,
			PreviewImage: p.PreviewImage,
			Price:        p.Price,
		},
		HexCode:  p.HexCode,
		RGBValue: p.RGBValue,
	}
}

// signetPO 映射 signet 表。
//
// v1 中印章只作为注册赠品，没有购买接口，因此 shop.Repository 未提供其查询方法；
// PO 保留以便将来接入。
type signetPO struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Name         string `gorm:"column:name"`
	PreviewImage string `gorm:"column:preview_image"`
	FilePath     string `gorm:"column:file_path"`
	Price        int    `gorm:"column:price"`
}

func (signetPO) TableName() string { return "signet" }

func (p signetPO) toDomain() shop.Signet {
	return shop.Signet{
		Item: shop.Item{
			ID:           p.ID,
			Type:         shop.ItemSignet,
			Name:         p.Name,
			PreviewImage: p.PreviewImage,
			Price:        p.Price,
		},
		FilePath: p.FilePath,
	}
}

// functionCardPO 映射 function_card 表。
//
// ReduceTime 与 SpeedRate 在表中是 varchar，分别存分钟数与倍率，
// 使用道具时会写入 letter 表的同名列。
type functionCardPO struct {
	ID              int64  `gorm:"column:id;primaryKey;autoIncrement"`
	CardType        int    `gorm:"column:card_type"`
	CardName        string `gorm:"column:card_name"`
	CardDesc        string `gorm:"column:card_desc"`
	CardPreviewLink string `gorm:"column:card_preview_link"`
	CardStatus      int    `gorm:"column:card_status"`
	ReduceTime      string `gorm:"column:reduce_time"`
	SpeedRate       string `gorm:"column:speed_rate"`
	Price           int    `gorm:"column:price"`
	Remark          string `gorm:"column:remark"`
}

func (functionCardPO) TableName() string { return "function_card" }

func (p functionCardPO) toDomain() shop.FunctionCard {
	return shop.FunctionCard{
		Item: shop.Item{
			ID:           p.ID,
			Type:         shop.ItemFunctionCard,
			Name:         p.CardName,
			Description:  p.CardDesc,
			PreviewImage: p.CardPreviewLink,
			Price:        p.Price,
		},
		CardType:      shop.CardType(p.CardType),
		Enabled:       p.CardStatus == 1,
		ReduceMinutes: atoiDefault(p.ReduceTime, 0),
		SpeedRate:     parseFloatDefault(p.SpeedRate, 1),
		Remark:        p.Remark,
	}
}

// fontPaperPO 映射 font_paper 表，记录字体与信纸组合下的适配字数。
//
// 注意：本表主键无自增，插入时需要显式指定 ID；本域只读。
type fontPaperPO struct {
	ID        int64 `gorm:"column:id;primaryKey"`
	PaperID   int64 `gorm:"column:paper_id"`
	FontID    int64 `gorm:"column:font_id"`
	FitNumber int64 `gorm:"column:fit_number"`
}

func (fontPaperPO) TableName() string { return "font_paper" }

func (p fontPaperPO) toDomain() shop.PaperFontFit {
	return shop.PaperFontFit{
		PaperID:   p.PaperID,
		FontID:    p.FontID,
		FitNumber: int(p.FitNumber),
	}
}

// commodityPO 映射 commodity 表，用于运营位展示的外部文创商品。
//
// Price 在表中是 varchar，因为展示的是价格文案而非可运算的金额，故不做数值转换。
type commodityPO struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	Price       string `gorm:"column:price"`
	Image       string `gorm:"column:image"`
	Marketing   string `gorm:"column:marketing"`
	Link        string `gorm:"column:link"`
}

func (commodityPO) TableName() string { return "commodity" }

func (p commodityPO) toDomain() shop.Commodity {
	return shop.Commodity{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Image:       p.Image,
		Marketing:   p.Marketing,
		Link:        p.Link,
	}
}

// parseFloatDefault 解析 varchar 中的数值，失败时返回兜底值。
//
// v1 的排版参数与加速倍率都存成字符串，存量表中确实存在空值与脏数据；
// 这里选择宽松处理而不报错，与 v1 由前端自行解析的行为保持一致。
func parseFloatDefault(s string, fallback float64) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fallback
	}
	return v
}
