package shop

import "context"

// Repository 是商城的数据访问接口，实现由 internal/infrastructure/mysql 提供。
//
// 本域是纯读的：购买时的扣费与背包写入都发生在 user 表，由 user.Repository 承担。
//
// v1 每类道具各有一张表，字段不完全相同（信纸多七个排版参数、功能卡多加速参数、
// 字体颜色多色值），因此按类型分别提供查询方法而非泛型接口——强行统一会退化成
// map[string]any，丢掉类型安全。
type Repository interface {
	ListFonts(ctx context.Context) ([]Font, error)
	// FindFont 查询字体，不存在时返回 ErrFontNotFound。
	FindFont(ctx context.Context, id int64) (Font, error)

	ListFontColors(ctx context.Context) ([]FontColor, error)
	// FindFontColor 查询字体颜色，不存在时返回 ErrFontColorNotFound。
	FindFontColor(ctx context.Context, id int64) (FontColor, error)

	ListPapers(ctx context.Context) ([]Paper, error)
	// FindPaper 查询信纸，不存在时返回 ErrPaperNotFound。
	//
	// 除购买时校验价格外，写信渲染也需要它携带的排版参数。
	FindPaper(ctx context.Context, id int64) (Paper, error)

	ListFunctionCards(ctx context.Context) ([]FunctionCard, error)
	// FindFunctionCard 查询功能卡，不存在时返回 ErrFunctionCardNotFound。
	//
	// 购买与使用道具都要用它——使用时需要 CardType 与加速参数。
	FindFunctionCard(ctx context.Context, id int64) (FunctionCard, error)

	// ListPaperFontFits 返回全部信纸与字体的适配字数，供写信页一次性加载。
	ListPaperFontFits(ctx context.Context) ([]PaperFontFit, error)

	// PageCommodities 分页查询运营位展示的外部文创商品。
	//
	// 只负责分页取数，v1 中对结果的随机打乱发生在应用层，本方法不做排序。
	PageCommodities(ctx context.Context, page, size int) (items []Commodity, total int64, err error)
}
