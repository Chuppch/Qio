package mysql

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/user"
)

// userRepository 实现 user.Repository。
type userRepository struct{ db *gorm.DB }

// NewUserRepository 构造用户仓储。
func NewUserRepository(db *gorm.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	po, err := userPOFrom(u)
	if err != nil {
		return fmt.Errorf("encode user: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(&po).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	// 回填自增主键与审计时间
	u.ID = po.ID
	u.CreatedAt = po.CreateTime
	u.UpdatedAt = po.UpdateTime
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	return r.findOne(ctx, "id = ?", id)
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return r.findOne(ctx, "email = ?", email)
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	return r.findOne(ctx, "username = ?", username)
}

func (r *userRepository) findOne(ctx context.Context, cond string, args ...any) (*user.User, error) {
	var po userPO
	err := r.db.WithContext(ctx).Where(cond, args...).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, user.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by %s: %w", cond, err)
	}
	return po.toDomain(), nil
}

// FindByIDs 批量查询。不存在的 ID 被跳过，返回顺序由数据库决定。
func (r *userRepository) FindByIDs(ctx context.Context, ids []int64) ([]*user.User, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var pos []userPO
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&pos).Error; err != nil {
		return nil, fmt.Errorf("find users by ids: %w", err)
	}

	out := make([]*user.User, 0, len(pos))
	for _, p := range pos {
		out = append(out, p.toDomain())
	}
	return out, nil
}

func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	po, err := userPOFrom(u)
	if err != nil {
		return fmt.Errorf("encode user: %w", err)
	}

	res := r.db.WithContext(ctx).Model(&userPO{}).Where("id = ?", u.ID).Updates(updatableUserFields(po))
	if res.Error != nil {
		return fmt.Errorf("update user %d: %w", u.ID, res.Error)
	}
	if res.RowsAffected == 0 {
		return user.ErrNotFound
	}
	return nil
}

// UpdateByEmail 按邮箱更新，用于重置密码——此时只有邮箱没有主键。
func (r *userRepository) UpdateByEmail(ctx context.Context, u *user.User) error {
	po, err := userPOFrom(u)
	if err != nil {
		return fmt.Errorf("encode user: %w", err)
	}

	res := r.db.WithContext(ctx).Model(&userPO{}).Where("email = ?", u.Email).Updates(updatableUserFields(po))
	if res.Error != nil {
		return fmt.Errorf("update user by email: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return user.ErrNotFound
	}
	return nil
}

// UpdateMoney 以增量方式调整余额。
//
// 用 SQL 表达式而非「读出再写回」，避免并发下互相覆盖。delta 为负时附加
// money + delta >= 0 条件，余额不足则不更新任何行并返回 ErrNotEnoughMoney。
func (r *userRepository) UpdateMoney(ctx context.Context, userID int64, delta int64) error {
	if delta == 0 {
		return nil
	}

	q := r.db.WithContext(ctx).Model(&userPO{}).Where("id = ?", userID)
	if delta < 0 {
		q = q.Where("money + ? >= 0", delta)
	}

	res := q.UpdateColumn("money", gorm.Expr("money + ?", delta))
	if res.Error != nil {
		return fmt.Errorf("update money of user %d: %w", userID, res.Error)
	}
	if res.RowsAffected == 0 {
		if delta < 0 {
			return user.ErrNotEnoughMoney
		}
		return user.ErrNotFound
	}
	return nil
}

func (r *userRepository) EmailTaken(ctx context.Context, email string) (bool, error) {
	return r.exists(ctx, "email = ?", email)
}

func (r *userRepository) UsernameTaken(ctx context.Context, username string) (bool, error) {
	return r.exists(ctx, "username = ?", username)
}

func (r *userRepository) exists(ctx context.Context, cond string, args ...any) (bool, error) {
	var n int64
	if err := r.db.WithContext(ctx).Model(&userPO{}).Where(cond, args...).Count(&n).Error; err != nil {
		return false, fmt.Errorf("count user by %s: %w", cond, err)
	}
	return n > 0, nil
}

// updatableUserFields 列出允许更新的列。
//
// 显式枚举而不用整对象 Updates：GORM 对结构体入参会跳过零值字段，导致
// 「把余额改成 0」这类操作静默失效；用 map 可以确保零值也被写入。
// 主键与创建时间不在其中，避免被意外改写。
func updatableUserFields(po userPO) map[string]any {
	return map[string]any{
		"username":       po.Username,
		"nickname":       po.Nickname,
		"email":          po.Email,
		"sex":            po.Sex,
		"avatar":         po.Avatar,
		"password":       po.Password,
		"status":         po.Status,
		"del_flag":       po.DelFlag,
		"login_ip":       po.LoginIP,
		"login_date":     po.LoginDate,
		"money":          po.Money,
		"addresses":      po.Addresses,
		"fonts":          po.Fonts,
		"papers":         po.Papers,
		"signets":        po.Signets,
		"font_colors":    po.FontColors,
		"function_cards": po.FunctionCards,
	}
}
