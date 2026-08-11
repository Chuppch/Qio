# qio-backend-v2

Qio（侨缘信使）v2 业务后端，Go + Gin 实现，替代 v1 的 Java/Spring Boot 版本（`qio-backend/`）。

当前状态：仅完成目录骨架，尚无业务实现。

## 模块路径

```
github.com/Chuppch/Qio/qio-backend-v2
```

## 目录结构

```
qio-backend-v2/
├── cmd/
│   └── server/                启动入口：加载配置、组装依赖、启动服务
├── internal/
│   ├── transport/             入站传输层
│   │   ├── http/              处理器与路由（router.go 集中注册路由）
│   │   ├── middleware/        认证、跨域、日志、限流、恢复
│   │   └── dto/               对外接口契约（Request / Response）
│   │
│   ├── app/                   跨域用例编排（一个用例一个文件）
│   │
│   ├── domain/                业务域，不依赖 Web 框架与存储实现
│   │   ├── user/              认证、账号、地址簿、签到任务
│   │   ├── letter/            信件规则与投递状态流转
│   │   ├── bottle/            漂流瓶
│   │   ├── friend/            好友关系与申请
│   │   ├── shop/              信纸/字体/明信片/功能卡、运营位
│   │   ├── explore/           答题、抽奖
│   │   ├── ai/                AI 能力（Agent Service 的业务侧门面）
│   │   └── port/              跨域能力接口（出站端口）
│   │       ├── notification/  通知投递
│   │       └── storage/       对象存储
│   │
│   ├── infra/                 基础设施实现
│   │   ├── mysql/             GORM 连接与各域 repository 实现
│   │   ├── redis/             缓存客户端与实现
│   │   ├── rabbitmq/          消息生产与消费
│   │   ├── minio/             storage 端口实现
│   │   ├── mail/              notification 端口实现
│   │   └── agentsvc/          Agent Service 出站客户端
│   │
│   └── config/                配置结构与加载
├── pkg/
│   └── response/              统一响应结构（可被仓库内其他 Go 项目复用）
├── config/                    配置文件（仅样例入库）
├── migrations/                数据库迁移脚本
├── Makefile
└── go.mod
```

## 分层与依赖方向

```
transport ──> app ──> domain <── infra
     │                  ▲
     └──────────────────┘
        （单域动作可直接调用域 service）
```

规则：

- `domain` 不得 import `transport`、`app`、`infra` 中的任何包。
- `domain` 只声明接口（域内 `repository.go`、跨域 `port/`），实现全在 `infra`。
- 跨多个域的业务动作放 `internal/app`，不塞进某一个域的 service。
- `transport` 只做协议转换、参数校验、错误码映射，不写业务规则。
- DTO 与领域模型分离，转换在 `transport/http` 的处理器中完成。
- 不新建 `utils`、`common`、`helpers` 这类无语义包。

## 域内文件约定

```
domain/<域>/
├── doc.go          包说明
├── service.go      业务逻辑
├── repository.go   数据访问接口定义
└── model.go        领域模型
```

接口多且主题分组明显时，按主题拆成 `service_auth.go`、`service_task.go` 等。同一域的所有文件共享同一个 package，不必新建子目录。域复杂到单目录难以维护时，再在域内拆 `application/` 等子包，不提前分层。

## 与 Agent Service 的结构映射

两个后端采用同构分层，便于团队在项目间切换：

| qio-agent-service | qio-backend-v2 |
|-------------------|----------------|
| `qio-agent-trigger/http/` | `internal/transport/http/` |
| `qio-agent-api/dto/` | `internal/transport/dto/` |
| `qio-agent-domain/` | `internal/domain/` |
| `qio-agent-domain/adapter/port/` | `internal/domain/port/` |
| `qio-agent-infrastructure/` | `internal/infra/` |
| `qio-agent-app/`（启动装配） | `cmd/server/` |
| `qio-agent-types/` | 不设；错误定义与常量归属各自的域或功能包 |

差异说明：Java 用 Maven 多模块在编译期强制依赖方向，Go 单 module 下目录只是约定。如需强制，后续在 `golangci-lint` 中加 depguard 规则，禁止 `internal/domain/**` 引用 `internal/infra/**` 与 `internal/transport/**`。

## 与 Agent Service 的调用边界

前端不直连 `qio-agent-service`。面向用户的 AI 能力由本服务的 `internal/domain/ai` 统一鉴权、限流、审计后，经 `internal/infra/agentsvc` 转调中台。模型、Prompt、Advisor、MCP、RAG 等中台概念不暴露到用户端。

## v1 接口归属对照

| v1 Controller | v2 处理器 | v2 业务域 |
|---------------|-----------|-----------|
| `UserController` | `transport/http/user.go` | `domain/user`、`domain/friend` |
| `LetterController` | `transport/http/letter.go` | `domain/letter` |
| `BottleController` | `transport/http/bottle.go`、`friend.go` | `domain/bottle`、`domain/friend` |
| `QuestionController`、`GameController` | `transport/http/explore.go` | `domain/explore` |
| `FontController`、`PaperController`、`CardController`、`MarketingController` | `transport/http/shop.go` | `domain/shop` |
| `CommonController`（上传） | `transport/http/upload.go` | `domain/port/storage` |
| `ChatService`（无 Controller） | `transport/http/ai.go` | `domain/ai` |

## 环境准备

需要 Go 1.26+。

```bash
go mod tidy
make build
```

## 常用命令

```bash
make help    # 查看全部命令
make run     # 本地启动
make fmt     # 格式化
make vet     # 静态检查
make test    # 运行测试
```

## 配置

复制 `config/config.example.yaml` 为 `config/config.dev.yaml` 后填入实际值。派生配置已被 `.gitignore` 排除，敏感项建议通过环境变量注入，不要提交真实凭据。
