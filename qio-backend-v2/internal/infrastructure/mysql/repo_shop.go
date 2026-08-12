package mysql

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/shop"
)

// shopRepository 实现 shop.Repository。
//
// 全部为读操作：购买时的扣费与背包写入发生在 user 表，由 user.Repository 承担。
type shopRepository struct{ db *gorm.DB }

// NewShopRepository 构造商城仓储。
func NewShopRepository(db *gorm.DB) shop.Repository {
	return &shopRepository{db: db}
}

// ---- 字体 ----

func (r *shopRepository) ListFonts(ctx context.Context) ([]shop.Font, error) {
	var pos []fontPO
	if err := r.db.WithContext(ctx).Find(&pos).Error; err != nil {
		return nil, fmt.Errorf("list fonts: %w", err)
	}

	out := make([]shop.Font, 0, len(pos))
	for _, p := range pos {
		out = append(out, p.toDomain())
	}
	return out, nil
}

func (r *shopRepository) FindFont(ctx context.Context, id int64) (shop.Font, error) {
	var po fontPO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return shop.Font{}, shop.ErrFontNotFound
	}
	if err != nil {
		return shop.Font{}, fmt.Errorf("find font %d: %w", id, err)
	}
	return po.toDomain(), nil
}

// ---- 字体颜色 ----

func (r *shopRepository) ListFontColors(ctx context.Context) ([]shop.FontColor, error) {
	var pos []fontColorPO
	if err := r.db.WithContext(ctx).Find(&pos).Error; err != nil {
		return nil, fmt.Errorf("list font colors: %w", err)
	}

	out := make([]shop.FontColor, 0, len(pos))
	for _, p := range pos {
		out = append(out, p.toDomain())
	}
	return out, nil
}

func (r *shopRepository) FindFontColor(ctx context.Context, id int64) (shop.FontColor, error) {
	var po fontColorPO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return shop.FontColor{}, shop.ErrFontColorNotFound
	}
	if err != nil {
		return shop.FontColor{}, fmt.Errorf("find font color %d: %w", id, err)
	}
	return po.toDomain(), nil
}

// ---- 信纸 ----

func (r *shopRepository) ListPapers(ctx context.Context) ([]shop.Paper, error) {
	var pos []paperPO
	if err := r.db.WithContext(ctx).Find(&pos).Error; err != nil {
		return nil, fmt.Errorf("list papers: %w", err)
	}

	out := make([]shop.Paper, 0, len(pos))
	for _, p := range pos {
		out = append(out, p.toDomain())
	}
	return out, nil
}

func (r *shopRepository) FindPaper(ctx context.Context, id int64) (shop.Paper, error) {
	var po paperPO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return shop.Paper{}, shop.ErrPaperNotFound
	}
	if err != nil {
		return shop.Paper{}, fmt.Errorf("find paper %d: %w", id, err)
	}
	return po.toDomain(), nil
}

// ---- 功能卡 ----

func (r *shopRepository) ListFunctionCards(ctx context.Context) ([]shop.FunctionCard, error) {
	var pos []functionCardPO
	if err := r.db.WithContext(ctx).Find(&pos).Error; err != nil {
		return nil, fmt.Errorf("list function cards: %w", err)
	}

	out := make([]shop.FunctionCard, 0, len(pos))
	for _, p := range pos {
		out = append(out, p.toDomain())
	}
	return out, nil
}

// FindFunctionCard 查询功能卡。
//
// 用 Where("id = ?") 而非 First(&po, id)：v1 的功能卡表中存在 id 为 0 的记录
// （注册赠品与「立即送达」特例卡），而 GORM 的主键简写形式会忽略零值条件。
func (r *shopRepository) FindFunctionCard(ctx context.Context, id int64) (shop.FunctionCard, error) {
	var po functionCardPO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return shop.FunctionCard{}, shop.ErrFunctionCardNotFound
	}
	if err != nil {
		return shop.FunctionCard{}, fmt.Errorf("find function card %d: %w", id, err)
	}
	return po.toDomain(), nil
}

// ---- 信纸字体适配 ----

func (r *shopRepository) ListPaperFontFits(ctx context.Context) ([]shop.PaperFontFit, error) {
	var pos []fontPaperPO
	if err := r.db.WithContext(ctx).Find(&pos).Error; err != nil {
		return nil, fmt.Errorf("list paper font fits: %w", err)
	}

	out := make([]shop.PaperFontFit, 0, len(pos))
	for _, p := range pos {
		out = append(out, p.toDomain())
	}
	return out, nil
}

// ---- 运营位商品 ----

// PageCommodities 分页查询外部文创商品。
//
// 只取数不排序：v1 在应用层对结果做随机打乱，该行为保留在 service 中。
func (r *shopRepository) PageCommodities(
	ctx context.Context, page, size int,
) ([]shop.Commodity, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	var total int64
	if err := r.db.WithContext(ctx).Model(&commodityPO{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count commodities: %w", err)
	}
	if total == 0 {
		return nil, 0, nil
	}

	var pos []commodityPO
	err := r.db.WithContext(ctx).
		Offset((page - 1) * size).
		Limit(size).
		Find(&pos).Error
	if err != nil {
		return nil, 0, fmt.Errorf("page commodities: %w", err)
	}

	out := make([]shop.Commodity, 0, len(pos))
	for _, p := range pos {
		out = append(out, p.toDomain())
	}
	return out, total, nil
}
