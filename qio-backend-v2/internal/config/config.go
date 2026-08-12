package config

import "time"

// 配置结构定义。
//
// 字段与 config/config.example.yaml 对应。敏感项（密码、密钥、AccessKey）
// 不写默认值，运行时从环境变量或本地不入库的配置文件注入。

// Config 是应用的完整配置。
type Config struct {
	Server       Server       `mapstructure:"server"`
	Log          Log          `mapstructure:"log"`
	MySQL        MySQL        `mapstructure:"mysql"`
	Redis        Redis        `mapstructure:"redis"`
	RabbitMQ     RabbitMQ     `mapstructure:"rabbitmq"`
	MinIO        MinIO        `mapstructure:"minio"`
	Mail         Mail         `mapstructure:"mail"`
	JWT          JWT          `mapstructure:"jwt"`
	AgentService AgentService `mapstructure:"agent_service"`
}

// Server 是 HTTP 服务配置。
type Server struct {
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// Log 是日志配置。
type Log struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// MySQL 是数据库连接配置。
type MySQL struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Database        string        `mapstructure:"database"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// Redis 是缓存连接配置。
type Redis struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// RabbitMQ 是消息队列连接配置。
type RabbitMQ struct {
	URL string `mapstructure:"url"`
}

// MinIO 是对象存储连接配置。
type MinIO struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	UseSSL    bool   `mapstructure:"use_ssl"`
}

// Mail 是邮件投递配置。
type Mail struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Nickname string `mapstructure:"nickname"`
}

// JWT 是令牌签发配置。
//
// 对应 v1 的 qiaopi.jwt 配置段。v1 另有一组被注释掉的 admin 配置，未启用，
// 不予迁移。
type JWT struct {
	Secret    string        `mapstructure:"secret"`
	ExpiresIn time.Duration `mapstructure:"expires_in"`
}

// AgentService 是 Qio Agent Service 的出站调用配置。
type AgentService struct {
	BaseURL string        `mapstructure:"base_url"`
	Timeout time.Duration `mapstructure:"timeout"`
}

// TODO: 实现 Load，按 profile 读取 config/config.<env>.yaml 并叠加环境变量覆盖。
