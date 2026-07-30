# 中间件配置

项目依赖的外部中间件及其配置说明。所有中间件推荐使用 Docker 部署。

## 端口总览

| 中间件 | 端口 | 用途 |
|--------|------|------|
| MySQL | 3306 | 主数据库 |
| Redis | 6379 | 缓存、Session、分布式锁 |
| RabbitMQ | 5672 (AMQP) / 15672 (管理面板) | AI 消息队列 |
| Minio | 9000 (API) / 9001 (控制台) | 对象存储（图片、字体等） |

## MySQL

```yaml
image: mysql:8.4
port: 3306
database: qiaopi
username: root
password: root（开发环境）
charset: utf8mb4
timezone: Asia/Shanghai
```

初始化数据位于 `qio-backend/init_qiaopi/` 目录下，导入 SQL 文件即可。

## Redis

```yaml
image: redis:8
port: 6379
password: root（开发环境）
database: 1
```

用于缓存用户信息、验证码、信件送达倒计时等。初始化数据也在 `init_qiaopi/` 中。

## RabbitMQ

```yaml
image: rabbitmq:4-management
port: 5672（AMQP）/ 15672（Web管理面板）
virtual-host: /
username: root
password: root（开发环境）
```

用于 AI 对话的异步消息处理，前端通过 WebSocket 接收流式响应，后端通过 RabbitMQ 消费 AI 模型返回的消息。

## Minio

```yaml
image: minio/minio
port: 9000（API）/ 9001（控制台）
access-key: 自行配置
secret-key: 自行配置
bucket: qiaopi
base-path: qiaopi-images/
```

用于存储用户上传的图片、信件渲染结果、头像等文件资源。

## 外部服务

### 智谱 ChatGLM

- 申请地址：https://bigmodel.cn/usercenter/apikeys
- 配置项：`api-key`
- 用途：AI 对话、文化问答

### QQ 邮箱 SMTP

- 开通地址：https://wx.mail.qq.com/list/readtemplate?name=app_intro.html#/agreement/authorizationCode
- 配置项：`username`（QQ 邮箱）、`password`（授权码）
- 用途：验证码发送、通知邮件

## 密钥配置

| 配置项 | 要求 | 用途 |
|--------|------|------|
| JWT secret-key | 32 字符以上 | 用户登录令牌签名 |
| AES secret-key | 32 字符 | 答题功能题目加密 |
