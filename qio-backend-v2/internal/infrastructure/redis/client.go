package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Chuppch/Qio/qio-backend-v2/internal/config"
)

// Client 是本包对 go-redis 客户端的别名，供构造函数签名使用。
type Client = redis.Client

// Open 建立 Redis 连接并做一次连通性探测。
func Open(cfg config.Redis) (*Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return c, nil
}
