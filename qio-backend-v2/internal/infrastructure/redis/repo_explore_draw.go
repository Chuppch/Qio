package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/explore"
)

// exploreDrawRepository 实现 explore.DrawRepository。
//
// 抽奖次数只存 Redis，沿用 v1 GameServiceImpl 的存储方式：一个用户一个键，
// 值是剩余次数的十进制文本，每次写入都重置存活时间。
type exploreDrawRepository struct{ rdb *Client }

// NewExploreDrawRepository 构造抽奖次数仓储。
func NewExploreDrawRepository(rdb *Client) explore.DrawRepository {
	return &exploreDrawRepository{rdb: rdb}
}

// drawLimitTTL 是抽奖次数的存活时间，取自 v1 的 Duration.ofDays(1)。
//
// 注意这是「自最后一次写入起 24 小时」，不是「自然日归零」：v1 每次扣减都会
// 重置存活时间，因此持续参与的用户其次数窗口会不断顺延。等价迁移保留该行为，
// 记录在 docs/TODO-migration.md。
const drawLimitTTL = 24 * time.Hour

// FindLimit 查询剩余次数，键不存在时按默认值初始化并写回。
//
// 对应 v1 getFflLimit。
func (r *exploreDrawRepository) FindLimit(ctx context.Context, userID int64) (explore.DrawLimit, error) {
	key := GameFortuneUserKey(userID)

	remaining, err := r.rdb.Get(ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		remaining = explore.DefaultDrawLimit
		if err := r.setLimit(ctx, key, remaining); err != nil {
			return explore.DrawLimit{}, err
		}
		return explore.DrawLimit{UserID: userID, Remaining: remaining}, nil
	}
	if err != nil {
		return explore.DrawLimit{}, fmt.Errorf("get draw limit of user %d: %w", userID, err)
	}

	return explore.DrawLimit{UserID: userID, Remaining: remaining}, nil
}

// Consume 扣减一次抽奖次数。
//
// 对应 v1 winFfl 的前半段：键缺失或次数非正时报错，否则减一后写回。
// v1 的读—判断—写回不是原子操作，同一用户并发抽奖可能少扣次数；等价迁移保留
// 该实现，记录在 docs/TODO-migration.md。
func (r *exploreDrawRepository) Consume(ctx context.Context, userID int64) (explore.DrawLimit, error) {
	key := GameFortuneUserKey(userID)

	remaining, err := r.rdb.Get(ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		// v1 在键缺失时同样按次数耗尽处理，不做初始化
		return explore.DrawLimit{}, explore.ErrDrawLimitExhausted
	}
	if err != nil {
		return explore.DrawLimit{}, fmt.Errorf("get draw limit of user %d: %w", userID, err)
	}
	if remaining <= 0 {
		return explore.DrawLimit{}, explore.ErrDrawLimitExhausted
	}

	remaining--
	if err := r.setLimit(ctx, key, remaining); err != nil {
		return explore.DrawLimit{}, err
	}
	return explore.DrawLimit{UserID: userID, Remaining: remaining}, nil
}

func (r *exploreDrawRepository) setLimit(ctx context.Context, key string, remaining int) error {
	if err := r.rdb.Set(ctx, key, strconv.Itoa(remaining), drawLimitTTL).Err(); err != nil {
		return fmt.Errorf("set draw limit %s: %w", key, err)
	}
	return nil
}
