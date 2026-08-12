package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/user"
)

// userTaskRepository 实现 user.TaskRepository。
//
// 数据全部落 Redis，沿用 v1 的存储方式：
//   - 任务模板存一份全局 JSON，用户当日首次访问时复制为个人副本
//   - 个人副本保留至次日，由 taskTTL 控制过期
//   - 签到用月度位图记录，一天一位
type userTaskRepository struct{ rdb *Client }

// NewUserTaskRepository 构造任务与签到仓储。
func NewUserTaskRepository(rdb *Client) user.TaskRepository {
	return &userTaskRepository{rdb: rdb}
}

const (
	// taskTTL 是个人任务副本的存活时间，跨过当日即可回收
	taskTTL = 36 * time.Hour
	// signBitmapTTL 是签到位图的存活时间，保留一年供日历回溯
	signBitmapTTL = 400 * 24 * time.Hour
)

// taskJSON 是任务在 Redis 中的存储结构。
//
// 字段名沿用 v1 序列化 TaskTable 时使用的名称，保证存量数据可读。
type taskJSON struct {
	ID          int64  `json:"id"`
	TaskName    string `json:"taskName"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	Money       int    `json:"money"`
	Link        string `json:"link"`
	Route       string `json:"route"`
}

func (t taskJSON) toDomain() user.Task {
	return user.Task{
		ID:          t.ID,
		Name:        t.TaskName,
		Description: t.Description,
		Status:      user.TaskStatus(t.Status),
		Reward:      t.Money,
		Link:        t.Link,
		Route:       t.Route,
	}
}

// signAwardJSON 是签到奖励配置在 Redis 中的存储结构。
type signAwardJSON struct {
	ID          int64  `json:"id"`
	AwardType   int    `json:"awardType"`
	AwardID     int64  `json:"awardId"`
	PreviewLink string `json:"previewLink"`
	AwardDesc   string `json:"awardDesc"`
	AwardName   string `json:"awardName"`
	AwardNum    int    `json:"awardNum"`
	SignDays    int    `json:"signDays"`
}

func (a signAwardJSON) toDomain() user.SignAward {
	return user.SignAward{
		ID:           a.ID,
		Type:         user.ItemType(a.AwardType),
		ItemID:       a.AwardID,
		Name:         a.AwardName,
		Description:  a.AwardDesc,
		PreviewURL:   a.PreviewLink,
		Count:        a.AwardNum,
		RequiredDays: a.SignDays,
	}
}

// ListTasks 返回用户当日任务。
//
// 个人副本不存在时从模板复制一份——这是 v1 的行为，保持一致。
func (r *userTaskRepository) ListTasks(ctx context.Context, userID int64, date string) ([]user.Task, error) {
	key := taskUserKey(userID, date)

	raw, err := r.rdb.Get(ctx, key).Bytes()
	switch {
	case errors.Is(err, redis.Nil):
		raw, err = r.rdb.Get(ctx, taskTemplateKey()).Bytes()
		if errors.Is(err, redis.Nil) {
			// 模板未配置，视作当日无任务
			return nil, nil
		}
		if err != nil {
			return nil, fmt.Errorf("get task template: %w", err)
		}
		if err := r.rdb.Set(ctx, key, raw, taskTTL).Err(); err != nil {
			return nil, fmt.Errorf("init user tasks: %w", err)
		}
	case err != nil:
		return nil, fmt.Errorf("get user tasks: %w", err)
	}

	var items []taskJSON
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("decode user tasks: %w", err)
	}

	out := make([]user.Task, 0, len(items))
	for _, it := range items {
		out = append(out, it.toDomain())
	}
	return out, nil
}

func (r *userTaskRepository) MarkTaskFinished(ctx context.Context, userID int64, date string, taskID int64) error {
	return r.setTaskStatus(ctx, userID, date, taskID, user.TaskFinished)
}

func (r *userTaskRepository) MarkTaskClaimed(ctx context.Context, userID int64, date string, taskID int64) error {
	return r.setTaskStatus(ctx, userID, date, taskID, user.TaskClaimed)
}

// setTaskStatus 修改个人副本中某条任务的状态。
//
// 读改写整个 JSON 数组是 v1 的做法，此处保持一致。并发下同一用户同时完成两个
// 任务可能丢失一次更新，属于已知问题，记录在 docs/TODO-migration.md。
func (r *userTaskRepository) setTaskStatus(
	ctx context.Context, userID int64, date string, taskID int64, status user.TaskStatus,
) error {
	key := taskUserKey(userID, date)

	raw, err := r.rdb.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return user.ErrTaskNotClaimable
	}
	if err != nil {
		return fmt.Errorf("get user tasks: %w", err)
	}

	var items []taskJSON
	if err := json.Unmarshal(raw, &items); err != nil {
		return fmt.Errorf("decode user tasks: %w", err)
	}

	found := false
	for i := range items {
		if items[i].ID == taskID {
			items[i].Status = int(status)
			found = true
			break
		}
	}
	if !found {
		return user.ErrTaskNotClaimable
	}

	updated, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("encode user tasks: %w", err)
	}

	ttl, err := r.rdb.TTL(ctx, key).Result()
	if err != nil || ttl <= 0 {
		ttl = taskTTL
	}
	if err := r.rdb.Set(ctx, key, updated, ttl).Err(); err != nil {
		return fmt.Errorf("save user tasks: %w", err)
	}
	return nil
}

// Sign 记录一次签到并返回当月累计签到天数。
//
// 用月度位图存储：SETBIT 的返回值是该位原先的值，为 1 说明当天已签过。
func (r *userTaskRepository) Sign(ctx context.Context, userID int64, date string) (int, error) {
	month, offset, err := splitSignDate(date)
	if err != nil {
		return 0, err
	}
	key := signedUserKey(userID, month)

	prev, err := r.rdb.SetBit(ctx, key, offset, 1).Result()
	if err != nil {
		return 0, fmt.Errorf("set sign bit: %w", err)
	}
	if prev == 1 {
		return 0, user.ErrAlreadySigned
	}

	if err := r.rdb.Expire(ctx, key, signBitmapTTL).Err(); err != nil {
		return 0, fmt.Errorf("expire sign bitmap: %w", err)
	}

	days, err := r.rdb.BitCount(ctx, key, nil).Result()
	if err != nil {
		return 0, fmt.Errorf("count sign bits: %w", err)
	}
	return int(days), nil
}

func (r *userTaskRepository) Signed(ctx context.Context, userID int64, date string) (bool, error) {
	month, offset, err := splitSignDate(date)
	if err != nil {
		return false, err
	}

	bit, err := r.rdb.GetBit(ctx, signedUserKey(userID, month), offset).Result()
	if err != nil {
		return false, fmt.Errorf("get sign bit: %w", err)
	}
	return bit == 1, nil
}

// SignedDays 返回指定月份已签到的日期，格式与入参 date 一致（形如 20260811）。
func (r *userTaskRepository) SignedDays(ctx context.Context, userID int64, month string) ([]string, error) {
	if len(month) != 6 {
		return nil, fmt.Errorf("invalid month %q, want yyyyMM", month)
	}

	raw, err := r.rdb.Get(ctx, signedUserKey(userID, month)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get sign bitmap: %w", err)
	}

	var out []string
	for i, b := range raw {
		for bit := 0; bit < 8; bit++ {
			if b&(1<<(7-bit)) == 0 {
				continue
			}
			// 位偏移从 0 起，对应当月 1 号
			day := i*8 + bit + 1
			out = append(out, month+fmt.Sprintf("%02d", day))
		}
	}
	return out, nil
}

func (r *userTaskRepository) ListSignAwards(ctx context.Context, date string) ([]user.SignAward, error) {
	raw, err := r.rdb.Get(ctx, signAwardKey(date)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get sign awards: %w", err)
	}

	var items []signAwardJSON
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("decode sign awards: %w", err)
	}

	out := make([]user.SignAward, 0, len(items))
	for _, it := range items {
		out = append(out, it.toDomain())
	}
	return out, nil
}

// splitSignDate 把 yyyyMMdd 拆成月份键与位偏移。
//
// 位偏移为「日 - 1」，即 1 号对应第 0 位。
func splitSignDate(date string) (month string, offset int64, err error) {
	if len(date) != 8 {
		return "", 0, fmt.Errorf("invalid date %q, want yyyyMMdd", date)
	}
	day, err := strconv.Atoi(date[6:])
	if err != nil || day < 1 || day > 31 {
		return "", 0, fmt.Errorf("invalid day in date %q", date)
	}
	return date[:6], int64(day - 1), nil
}
