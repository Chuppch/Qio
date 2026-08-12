package mysql

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Chuppch/Qio/qio-backend-v2/internal/config"
)

// Open 建立 MySQL 连接并配置连接池。
//
// 不使用 GORM 的自动迁移：表结构由 migrations/ 下的 SQL 脚本管理，
// 避免运行时改表。
func Open(cfg config.MySQL) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
		// 表名已由各 PO 的 TableName 显式指定，禁用复数化推断
		NamingStrategy: nil,
	})
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	pool, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}

	if cfg.MaxOpenConns > 0 {
		pool.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		pool.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		pool.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	pool.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}
