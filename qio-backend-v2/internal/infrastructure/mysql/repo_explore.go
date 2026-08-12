package mysql

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/explore"
)

// exploreRepository 实现 explore.Repository。
type exploreRepository struct{ db *gorm.DB }

// NewExploreRepository 构造文化探索仓储。
func NewExploreRepository(db *gorm.DB) explore.Repository {
	return &exploreRepository{db: db}
}

// questionListColumns 是出题时下发的列。
//
// 与 v1 startAnswer 的 select 保持一致：不含 correct_answer 与 explanation，
// 避免答案随出题响应一起下发。set_sequence_id 同样不在 v1 的列表中。
var questionListColumns = []string{
	"id", "set_id", "content",
	"option_a", "option_b", "option_c", "option_d",
}

// ListQuestions 查询整套题目，含正确答案。
//
// 套号不存在时返回空切片，不返回错误——v1 答完最后一套会去查越界的下一套，
// 依赖的正是空结果。
func (r *exploreRepository) ListQuestions(ctx context.Context, setID int64) ([]*explore.Question, error) {
	var pos []questionPO
	if err := r.db.WithContext(ctx).Where("set_id = ?", setID).Find(&pos).Error; err != nil {
		return nil, fmt.Errorf("list questions of set %d: %w", setID, err)
	}
	return questionsToDomain(pos), nil
}

// ListQuestionsWithoutAnswer 查询整套题目，但不带正确答案。
//
// v1 查完后用 Collections.shuffle 打乱顺序，那属于出题策略，留在 service。
func (r *exploreRepository) ListQuestionsWithoutAnswer(
	ctx context.Context, setID int64,
) ([]*explore.Question, error) {
	var pos []questionPO
	err := r.db.WithContext(ctx).
		Select(questionListColumns).
		Where("set_id = ?", setID).
		Find(&pos).Error
	if err != nil {
		return nil, fmt.Errorf("list questions of set %d without answer: %w", setID, err)
	}
	return questionsToDomain(pos), nil
}

func (r *exploreRepository) FindQuestion(ctx context.Context, id int64) (*explore.Question, error) {
	var po questionPO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, explore.ErrQuestionNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find question %d: %w", id, err)
	}
	return po.toDomain(), nil
}

// FindProgress 查询用户答题进度。
//
// user_id 上有唯一索引，因此取一条即可。无记录不是异常，调用方据此判定首次答题。
func (r *exploreRepository) FindProgress(ctx context.Context, userID int64) (*explore.Progress, error) {
	var po questionUserStatusPO
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, explore.ErrProgressNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find question progress of user %d: %w", userID, err)
	}
	return po.toDomain(), nil
}

func (r *exploreRepository) CreateProgress(ctx context.Context, p *explore.Progress) error {
	po := questionUserStatusPOFrom(p)

	if err := r.db.WithContext(ctx).Create(&po).Error; err != nil {
		return fmt.Errorf("create question progress of user %d: %w", p.UserID, err)
	}

	p.ID = po.ID
	p.CreatedAt = po.CreateTime
	p.UpdatedAt = po.UpdateTime
	return nil
}

// SaveProgress 保存各套题的完成状态。
//
// v1 用 updateById 整行更新，这里只写十个进度列与 update_user：用 map 而非结构体，
// 使值为 0 的列也能落库——GORM 传结构体会跳过零值，无法把已完成改回未完成。
func (r *exploreRepository) SaveProgress(ctx context.Context, p *explore.Progress) error {
	po := questionUserStatusPOFrom(p)

	values := po.setColumnValues()
	values["update_user"] = po.UpdateUser

	res := r.db.WithContext(ctx).Model(&questionUserStatusPO{}).
		Where("id = ?", p.ID).
		Updates(values)
	if res.Error != nil {
		return fmt.Errorf("save question progress %d: %w", p.ID, res.Error)
	}
	if res.RowsAffected == 0 {
		return explore.ErrProgressNotFound
	}
	return nil
}

func questionsToDomain(pos []questionPO) []*explore.Question {
	out := make([]*explore.Question, 0, len(pos))
	for _, p := range pos {
		out = append(out, p.toDomain())
	}
	return out
}
