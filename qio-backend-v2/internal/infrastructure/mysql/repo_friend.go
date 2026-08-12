package mysql

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/friend"
)

// friendRepository 实现 friend.Repository。
type friendRepository struct{ db *gorm.DB }

// NewFriendRepository 构造好友仓储。
func NewFriendRepository(db *gorm.DB) friend.Repository {
	return &friendRepository{db: db}
}

// ---- 好友关系 ----

func (r *friendRepository) ListByOwner(ctx context.Context, ownerID int64) ([]*friend.Friend, error) {
	var pos []friendPO
	if err := r.db.WithContext(ctx).Where("owning_id = ?", ownerID).Find(&pos).Error; err != nil {
		return nil, fmt.Errorf("list friends of owner %d: %w", ownerID, err)
	}

	out := make([]*friend.Friend, 0, len(pos))
	for _, p := range pos {
		out = append(out, p.toDomain())
	}
	return out, nil
}

func (r *friendRepository) FindByID(ctx context.Context, id int64) (*friend.Friend, error) {
	var po friendPO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, friend.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find friend %d: %w", id, err)
	}
	return po.toDomain(), nil
}

func (r *friendRepository) FindByIDAndOwner(ctx context.Context, id, ownerID int64) (*friend.Friend, error) {
	var po friendPO
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Where("owning_id = ?", ownerID).
		First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, friend.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find friend %d of owner %d: %w", id, ownerID, err)
	}
	return po.toDomain(), nil
}

// FindByUserAndOwner 按好友的用户 ID 与归属者查询。
//
// 表上没有 (owning_id, user_id) 唯一约束，重复的好友记录客观上可能存在；
// 此处取第一条，与 v1 selectOne 在有重复时的行为差异记录在
// docs/TODO-migration.md。
func (r *friendRepository) FindByUserAndOwner(ctx context.Context, userID, ownerID int64) (*friend.Friend, error) {
	var po friendPO
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("owning_id = ?", ownerID).
		First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, friend.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find friend of user %d owned by %d: %w", userID, ownerID, err)
	}
	return po.toDomain(), nil
}

func (r *friendRepository) Create(ctx context.Context, f *friend.Friend) error {
	po, err := friendPOFrom(f)
	if err != nil {
		return fmt.Errorf("encode friend: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(&po).Error; err != nil {
		return fmt.Errorf("create friend of owner %d: %w", f.OwnerID, err)
	}

	f.ID = po.ID
	f.CreatedAt = valueOrZero(po.CreateTime)
	f.UpdatedAt = valueOrZero(po.UpdateTime)
	return nil
}

func (r *friendRepository) Update(ctx context.Context, f *friend.Friend) error {
	values, err := friendUpdateValues(f)
	if err != nil {
		return err
	}

	res := r.db.WithContext(ctx).Model(&friendPO{}).Where("id = ?", f.ID).Updates(values)
	if res.Error != nil {
		return fmt.Errorf("update friend %d: %w", f.ID, res.Error)
	}
	if res.RowsAffected == 0 {
		return friend.ErrNotFound
	}
	return nil
}

// UpdateAll 逐条更新好友记录。
//
// v1 用 MyBatis-Plus 的 updateById(Collection) 批量提交，最终同样是每条一句
// UPDATE。这里在一个事务里顺序执行，比 v1 多了原子性——v1 的 getMyFriends
// 没有事务注解。
func (r *friendRepository) UpdateAll(ctx context.Context, fs []*friend.Friend) error {
	if len(fs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, f := range fs {
			values, err := friendUpdateValues(f)
			if err != nil {
				return err
			}
			if err := tx.Model(&friendPO{}).Where("id = ?", f.ID).Updates(values).Error; err != nil {
				return fmt.Errorf("update friend %d: %w", f.ID, err)
			}
		}
		return nil
	})
}

// friendUpdateValues 组装可更新的列。
//
// 用 map 而非结构体：GORM 传结构体会跳过零值，备注清空、地址簿删到空都会写不进去。
func friendUpdateValues(f *friend.Friend) (map[string]any, error) {
	po, err := friendPOFrom(f)
	if err != nil {
		return nil, fmt.Errorf("encode friend: %w", err)
	}

	return map[string]any{
		"name":        po.Name,
		"sex":         po.Sex,
		"email":       po.Email,
		"avatar":      po.Avatar,
		"addresses":   po.Addresses,
		"remark":      po.Remark,
		"update_user": po.UpdateUser,
	}, nil
}

// ---- 好友申请 ----

func (r *friendRepository) ListPendingRequests(ctx context.Context, receiverID int64) ([]*friend.Request, error) {
	var pos []friendRequestPO
	err := r.db.WithContext(ctx).
		Where("receiver_id = ?", receiverID).
		Where("status = ?", int(friend.RequestPending)).
		Find(&pos).Error
	if err != nil {
		return nil, fmt.Errorf("list pending friend requests of user %d: %w", receiverID, err)
	}

	out := make([]*friend.Request, 0, len(pos))
	for _, p := range pos {
		out = append(out, p.toDomain())
	}
	return out, nil
}

func (r *friendRepository) FindRequestByID(ctx context.Context, id int64) (*friend.Request, error) {
	var po friendRequestPO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, friend.ErrRequestUnavailable
	}
	if err != nil {
		return nil, fmt.Errorf("find friend request %d: %w", id, err)
	}
	return po.toDomain(), nil
}

// CountRequests 统计指定方向上的申请数量，不区分状态。
func (r *friendRepository) CountRequests(ctx context.Context, senderID, receiverID int64) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&friendRequestPO{}).
		Where("sender_id = ?", senderID).
		Where("receiver_id = ?", receiverID).
		Count(&n).Error
	if err != nil {
		return 0, fmt.Errorf("count friend requests from %d to %d: %w", senderID, receiverID, err)
	}
	return n, nil
}

func (r *friendRepository) CreateRequest(ctx context.Context, req *friend.Request) error {
	po, err := friendRequestPOFrom(req)
	if err != nil {
		return fmt.Errorf("encode friend request: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(&po).Error; err != nil {
		return fmt.Errorf("create friend request from %d to %d: %w", req.SenderID, req.ReceiverID, err)
	}

	req.ID = po.ID
	req.CreatedAt = valueOrZero(po.CreateTime)
	req.UpdatedAt = valueOrZero(po.UpdateTime)
	return nil
}

// UpdateRequest 更新申请状态。
//
// v1 用 updateById 整行更新，但除状态外没有别的字段会变；这里只写 status 与
// update_user。v1 在调用前手动 setUpdateTime，那个值会被自动填充覆盖掉，
// 因此不迁移该赋值。
func (r *friendRepository) UpdateRequest(ctx context.Context, req *friend.Request) error {
	res := r.db.WithContext(ctx).Model(&friendRequestPO{}).
		Where("id = ?", req.ID).
		Updates(map[string]any{
			"status":      int(req.Status),
			"update_user": req.ReceiverID,
		})
	if res.Error != nil {
		return fmt.Errorf("update friend request %d: %w", req.ID, res.Error)
	}
	if res.RowsAffected == 0 {
		return friend.ErrRequestUnavailable
	}
	return nil
}
