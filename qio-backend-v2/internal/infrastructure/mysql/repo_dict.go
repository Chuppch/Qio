package mysql

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/dict"
)

// dictRepository 实现 dict.Repository。
type dictRepository struct{ db *gorm.DB }

// NewDictRepository 构造字典仓储。
func NewDictRepository(db *gorm.DB) dict.Repository {
	return &dictRepository{db: db}
}

func (r *dictRepository) ListAvatars(ctx context.Context) ([]dict.Avatar, error) {
	var pos []avatarPO
	if err := r.db.WithContext(ctx).Find(&pos).Error; err != nil {
		return nil, fmt.Errorf("list avatars: %w", err)
	}

	out := make([]dict.Avatar, 0, len(pos))
	for _, p := range pos {
		out = append(out, p.toDomain())
	}
	return out, nil
}

func (r *dictRepository) FindAvatar(ctx context.Context, id int64) (dict.Avatar, error) {
	var po avatarPO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dict.Avatar{}, dict.ErrAvatarNotFound
	}
	if err != nil {
		return dict.Avatar{}, fmt.Errorf("find avatar %d: %w", id, err)
	}
	return po.toDomain(), nil
}

func (r *dictRepository) ListCountries(ctx context.Context) ([]dict.Country, error) {
	var pos []countryPO
	if err := r.db.WithContext(ctx).Find(&pos).Error; err != nil {
		return nil, fmt.Errorf("list countries: %w", err)
	}

	out := make([]dict.Country, 0, len(pos))
	for _, p := range pos {
		out = append(out, p.toDomain())
	}
	return out, nil
}

func (r *dictRepository) FindCountry(ctx context.Context, id int64) (dict.Country, error) {
	var po countryPO
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dict.Country{}, dict.ErrCountryNotFound
	}
	if err != nil {
		return dict.Country{}, fmt.Errorf("find country %d: %w", id, err)
	}
	return po.toDomain(), nil
}
