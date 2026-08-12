package mysql

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/bottle"
)

// bottleRepository 实现 bottle.Repository。
type bottleRepository struct{ db *gorm.DB }

// NewBottleRepository 构造漂流瓶仓储。
func NewBottleRepository(db *gorm.DB) bottle.Repository {
	return &bottleRepository{db: db}
}

func (r *bottleRepository) Create(ctx context.Context, b *bottle.Bottle) error {
	po, err := bottlePOFrom(b)
	if err != nil {
		return fmt.Errorf("encode bottle: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(&po).Error; err != nil {
		return fmt.Errorf("create bottle: %w", err)
	}

	b.ID = po.ID
	b.CreatedAt = po.CreateTime
	b.UpdatedAt = po.UpdateTime
	return nil
}

func (r *bottleRepository) FindByID(ctx context.Context, id int64) (*bottle.Bottle, error) {
	var po bottlePO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, bottle.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find bottle %d: %w", id, err)
	}
	return po.toDomain(), nil
}

// ListAvailable 查询可捞的漂流瓶。
//
// 条件与 v1 的 getNotIsPickedBottles 一致：未被捞起、非本人投放、非本人捞过。
// update_user 在表中是 varchar，因此比较时把用户 ID 转成字符串。
//
// 沿用 v1 的做法全量取出，由 service 随机选一个。数据量增长后应改为
// ORDER BY RAND() LIMIT 1 或按主键区间随机，记录在 docs/TODO-migration.md。
func (r *bottleRepository) ListAvailable(ctx context.Context, userID int64) ([]*bottle.Bottle, error) {
	uid := strconv.FormatInt(userID, 10)

	var pos []bottlePO
	err := r.db.WithContext(ctx).
		Where("is_picked = ?", false).
		Where("user_id <> ?", userID).
		Where("update_user IS NULL OR update_user <> ?", uid).
		Find(&pos).Error
	if err != nil {
		return nil, fmt.Errorf("list available bottles: %w", err)
	}

	out := make([]*bottle.Bottle, 0, len(pos))
	for _, p := range pos {
		out = append(out, p.toDomain())
	}
	return out, nil
}

// FindLatestPicked 查询用户最近捞起的一个漂流瓶。
//
// 对应 v1 的 getMostRecentBottleByUserId：按 update_user 匹配、update_time 倒序取一条。
func (r *bottleRepository) FindLatestPicked(ctx context.Context, userID int64) (*bottle.Bottle, error) {
	var po bottlePO
	err := r.db.WithContext(ctx).
		Where("update_user = ?", strconv.FormatInt(userID, 10)).
		Order("update_time DESC").
		First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, bottle.ErrNotPickedByUser
	}
	if err != nil {
		return nil, fmt.Errorf("find latest picked bottle of user %d: %w", userID, err)
	}
	return po.toDomain(), nil
}

func (r *bottleRepository) Update(ctx context.Context, b *bottle.Bottle) error {
	po, err := bottlePOFrom(b)
	if err != nil {
		return fmt.Errorf("encode bottle: %w", err)
	}

	res := r.db.WithContext(ctx).Model(&bottlePO{}).Where("id = ?", b.ID).
		Updates(map[string]any{
			"is_picked":      po.Picked,
			"update_user":    po.UpdateUser,
			"content":        po.Content,
			"sender_address": po.SenderAddress,
			"bottle_url":     po.BottleURL,
			"remark":         po.Remark,
		})
	if res.Error != nil {
		return fmt.Errorf("update bottle %d: %w", b.ID, res.Error)
	}
	if res.RowsAffected == 0 {
		return bottle.ErrNotFound
	}
	return nil
}
