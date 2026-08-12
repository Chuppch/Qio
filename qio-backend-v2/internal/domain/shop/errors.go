package shop

import "errors"

// 商城域的错误值。
//
// 逐条对应 v1 的异常与 i18n 消息码，注释中给出出处以便核对。
var (
	// ErrFontNotFound 字体不存在。对应 font.not.exists
	ErrFontNotFound = errors.New("font not found")

	// ErrFontColorNotFound 字体颜色不存在。对应 font.color.not.exists
	ErrFontColorNotFound = errors.New("font color not found")

	// ErrPaperNotFound 信纸不存在。对应 paper.not.exists
	ErrPaperNotFound = errors.New("paper not found")

	// ErrSignetNotFound 印章不存在
	ErrSignetNotFound = errors.New("signet not found")

	// ErrFunctionCardNotFound 功能卡不存在。对应 card.not.exists
	ErrFunctionCardNotFound = errors.New("function card not found")

	// ErrPaperFontFitNotFound 该信纸与字体的组合没有配置适配字数
	ErrPaperFontFitNotFound = errors.New("paper font fit not found")

	// ErrAlreadyOwned 用户已拥有该商品。对应 font.own、font.color.own、
	// paper.own 等一组消息码
	ErrAlreadyOwned = errors.New("item already owned")

	// ErrItemTypeInvalid 商品类型不合法
	ErrItemTypeInvalid = errors.New("invalid item type")
)
