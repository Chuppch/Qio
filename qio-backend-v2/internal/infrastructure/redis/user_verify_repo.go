package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/user"
)

// userVerifyCodeRepository 实现 user.VerifyCodeRepository。
type userVerifyCodeRepository struct{ rdb *Client }

// NewUserVerifyCodeRepository 构造验证码仓储。
func NewUserVerifyCodeRepository(rdb *Client) user.VerifyCodeRepository {
	return &userVerifyCodeRepository{rdb: rdb}
}

func (r *userVerifyCodeRepository) SaveCaptcha(ctx context.Context, key, code string, ttl time.Duration) error {
	if err := r.rdb.Set(ctx, captchaKey(key), code, ttl).Err(); err != nil {
		return fmt.Errorf("save captcha: %w", err)
	}
	return nil
}

// TakeCaptcha 取出并删除图形验证码。
//
// 用 GETDEL 保证「取出即失效」是原子的，避免同一个验证码被并发复用。
func (r *userVerifyCodeRepository) TakeCaptcha(ctx context.Context, key string) (string, error) {
	code, err := r.rdb.GetDel(ctx, captchaKey(key)).Result()
	if errors.Is(err, redis.Nil) {
		return "", user.ErrCodeExpired
	}
	if err != nil {
		return "", fmt.Errorf("take captcha: %w", err)
	}
	return code, nil
}

func (r *userVerifyCodeRepository) SaveEmailCode(
	ctx context.Context, scene user.CodeScene, email, code string, ttl time.Duration,
) error {
	if err := r.rdb.Set(ctx, emailCodeKey(scene, email), code, ttl).Err(); err != nil {
		return fmt.Errorf("save email code: %w", err)
	}
	return nil
}

func (r *userVerifyCodeRepository) FindEmailCode(
	ctx context.Context, scene user.CodeScene, email string,
) (string, error) {
	code, err := r.rdb.Get(ctx, emailCodeKey(scene, email)).Result()
	if errors.Is(err, redis.Nil) {
		return "", user.ErrCodeExpired
	}
	if err != nil {
		return "", fmt.Errorf("find email code: %w", err)
	}
	return code, nil
}

func (r *userVerifyCodeRepository) DeleteEmailCode(
	ctx context.Context, scene user.CodeScene, email string,
) error {
	if err := r.rdb.Del(ctx, emailCodeKey(scene, email)).Err(); err != nil {
		return fmt.Errorf("delete email code: %w", err)
	}
	return nil
}

func (r *userVerifyCodeRepository) EmailCodePending(
	ctx context.Context, scene user.CodeScene, email string,
) (bool, error) {
	n, err := r.rdb.Exists(ctx, emailCodeKey(scene, email)).Result()
	if err != nil {
		return false, fmt.Errorf("check email code: %w", err)
	}
	return n > 0, nil
}
