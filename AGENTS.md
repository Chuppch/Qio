# Qio — 侨缘信使 v2.0.0

## 项目定位

Qio 是基于原 Qiaopi v1.x 的第二次全面迭代升级（v2.0.0），在保留侨批文化互动体验核心玩法的基础上，对整个前后端项目进行系统性重构，包括语言与技术栈的迁移升级、架构重新设计、新功能需求开发、页面与业务流程优化，最终交付一个全面焕新的 v2 版本。

- 前端从 Vue 2 迁移至 React（Next.js）
- 后端从 Java/Spring Boot 迁移至 Go（Gin）
- 构建工具、中间件客户端库、ORM 等全面替换为 Go/React 生态对应方案
- 单体架构优先，API 规范化、多环境配置、容器化部署
- 在原有写信、漂流瓶、AI 对话、答题探索基础上扩展新功能
- UI 重设计与业务流程优化，解决 v1 遗留技术债

## 技术栈迁移

### 前端

| v1 | v2 |
|----|-----|
| Vue 2 + Vuex + Vue Router | Next.js (React 18 + TypeScript + App Router) |
| Element UI | shadcn/ui + Tailwind CSS + framer-motion |
| Webpack (vue-cli) | Next.js 内置 (Turbopack) |
| 无 SSR | SSG + CSR 混合渲染 |

### 后端

| v1 | v2 |
|----|-----|
| Spring Boot 3.3.3 | Gin (Go 1.26+) |
| MyBatis-Plus | GORM |
| Druid 连接池 | database/sql 内置连接池 |
| Spring Data Redis | go-redis |
| Spring AMQP (RabbitMQ) | amqp091-go |
| Minio Java SDK | minio-go |
| jjwt | golang-jwt |
| 智谱 ChatGLM Java SDK | HTTP API 直调 |
| Spring Mail | gomail |
| Graphics2D (信件渲染) | 待定（gg 库 / 前端 Canvas 渲染） |

### 基础设施（不变）

- MySQL
- Redis
- RabbitMQ
- Minio

## 业务后端 v2（qio-backend-v2）

Go module 路径为 `github.com/Chuppch/Qio/qio-backend-v2`。后端采用与 Agent Service 语义相近的四层结构：`trigger`、`app`、`domain`、`infrastructure`。当前 `user`、`dict` 已开始实现领域模型与存储适配，其余业务域仍以迁移骨架和渐进实现为主。

### 目录结构

```text
qio-backend-v2/
├── cmd/server/                Composition Root：加载配置、组装依赖、启动服务
├── internal/
│   ├── trigger/               入站触发层，对应 Agent Service 的 trigger
│   │   ├── http/              Handler 与路由（router.go 集中注册）
│   │   ├── middleware/        认证、跨域、日志、限流、恢复
│   │   └── dto/               HTTP Request / Response 契约
│   ├── app/                   Application 层，按完整业务用例编排
│   ├── domain/                领域层，不依赖 Web 框架与存储实现
│   │   ├── user/ letter/ bottle/ friend/ shop/ explore/ ai/
│   │   └── port/              跨域能力接口（notification、storage）
│   ├── infrastructure/        出站适配与基础设施实现
│   │   └── mysql/ redis/ rabbitmq/ minio/ mail/ agentsvc/
│   └── config/                配置结构与加载
├── pkg/response/              统一响应结构
├── config/                    配置文件（仅样例入库）
└── migrations/                数据库迁移脚本
```

### 分层与依赖方向

```text
HTTP / Job / Event
       │
       ▼
    trigger ──> app ──> domain <── infrastructure
                  ▲                    │
                  └──── cmd/server ────┘
                         仅负责装配
```

- `trigger` 是入站适配层，只负责协议解析、基础参数校验、身份上下文提取、DTO 转换和错误码映射；不得编写业务规则、直接执行 SQL 或操作 Redis。
- `app` 是 Application 层，表达一个完整用例，例如注册、登录、重置密码、发信、领取任务奖励。它负责协调领域服务、Repository、缓存、通知、鉴权和事务边界，但不保存 HTTP DTO，也不直接依赖 GORM、Redis Client 等实现类型。
- `domain` 是核心业务层，包含聚合、实体、值对象、领域服务、业务错误及所需的 Repository/Port 接口。它不得引用 `trigger`、`app`、`infrastructure` 或具体框架。
- `infrastructure` 是出站适配层，实现 `domain` 声明的 Repository 与 Port，负责 MySQL、Redis、RabbitMQ、Minio、邮件和 Agent Service 调用；它可以依赖 `domain`，业务层不得反向依赖它。
- `cmd/server` 是 Composition Root，可以同时引用各层完成实例化与依赖注入，但只能放配置加载、资源初始化、路由注册、生命周期和优雅关闭，不得承载业务逻辑。
- 正常请求必须走 `trigger -> app -> domain`。不要让 Handler 为了省一个用例文件而直接拼接多个 Repository；简单查询也应由 `app` 暴露清晰的查询用例，保持统一入口。
- 目录约定目前没有编译期强制。如需强制，在 `golangci-lint` 中增加 `depguard` 规则，不通过继续加深目录层级解决。

### 各层文件组织

Application 层按业务域建立目录，同一个 `Service` 的方法可以分散在多个用例文件中：

```text
app/user/
├── service.go            依赖集合与构造函数
├── register.go           注册用例
├── login.go              登录用例
├── reset_password.go     重置密码用例
├── update_profile.go     更新资料用例
├── address.go            地址簿用例
└── growth.go             签到与任务用例
```

- 不为每个文件重复创建 `RegisterService`、`LoginService` 等空壳类型；优先让同一业务域共享一个应用服务及其依赖。
- 用例文件按业务意图命名，不按 Controller 方法或数据库表命名。
- 跨域用例仍归最主要的业务意图，例如“发信”放 `app/letter/send.go`，内部协调用户扣费、信件存储和通知。

### 域内文件约定

```text
domain/<域>/
├── doc.go          包说明
├── model.go        聚合、实体和值对象
├── service.go      不依赖基础设施的领域规则
├── repository.go   该领域所需的数据访问接口
├── errors.go       可由上层识别的领域错误
└── validation.go   纯校验规则（确有需要时）
```

- 接口多且主题分组明显时，按主题拆成 `service_auth.go`、`service_task.go` 等；同一域所有文件共享同一个 package，不新建子目录。
- `repository.go` 只声明业务需要的接口，不出现 GORM Model、Redis Key 或 SQL 参数。
- `trigger/http` 与 `trigger/dto` 按业务域一一对应；Handler 不复用数据库 PO 作为响应 DTO。

基础设施当前按技术实现组织，并通过文件前缀区分持久化对象和仓储实现：

```text
infrastructure/
├── mysql/
│   ├── db.go              连接与连接池
│   ├── po_user.go         user 表映射及 PO/Domain 转换
│   ├── repo_user.go       user.Repository 的 MySQL 实现
│   ├── po_<domain>.go     其他领域的持久化对象
│   └── repo_<domain>.go   其他领域的 Repository 实现
├── redis/
│   ├── client.go
│   ├── user_verify_repo.go
│   └── user_task_repo.go
└── agentsvc/ mail/ minio/ rabbitmq/
```

- PO 只存在于 `infrastructure/mysql`，不得流入 `domain`、`app` 或 `trigger`。
- Repository 实现负责 PO 与领域模型转换，并把 GORM/Redis 错误翻译成领域错误或带上下文的基础设施错误。
- 一个 Repository 面向一个聚合或明确的查询能力，而不是机械地“一张表一个 Repository”。
- SQL、GORM 条件、Redis Key 和第三方 SDK 只允许出现在 `infrastructure`。

### v1 模块归属

| v1 Controller | v2 处理器 | v2 业务域 |
|---------------|-----------|-----------|
| `UserController` | `trigger/http/user.go` | `app/user`、`domain/user`、`domain/friend` |
| `LetterController` | `trigger/http/letter.go` | `app/letter`、`domain/letter` |
| `BottleController` | `trigger/http/bottle.go`、`friend.go` | `app/bottle`、`domain/bottle`、`domain/friend` |
| `QuestionController`、`GameController` | `trigger/http/explore.go` | `app/explore`、`domain/explore` |
| `FontController`、`PaperController`、`CardController`、`MarketingController` | `trigger/http/shop.go` | `app/shop`、`domain/shop` |
| `CommonController`（上传） | `trigger/http/upload.go` | 对应用例、`domain/port/storage` |
| `ChatService` | `trigger/http/ai.go` | `app/ai`、`domain/ai` |

### 与 Agent Service 的结构映射

两个后端保持同构分层，便于在项目间切换：

| qio-agent-service | qio-backend-v2 |
|-------------------|----------------|
| `qio-agent-trigger/http/` | `internal/trigger/http/` |
| `qio-agent-api/dto/` | `internal/trigger/dto/` |
| `qio-agent-domain/` | `internal/domain/` |
| `qio-agent-domain/adapter/port/` | `internal/domain/port/` |
| Agent Domain 的用例/服务编排 | `internal/app/` |
| `qio-agent-infrastructure/` | `internal/infrastructure/` |
| `qio-agent-app/` | `cmd/server/` |
| `qio-agent-types/` | 不设；错误定义与常量归属各自的域或功能包 |

这里统一的是模块职责和语言，而不是照搬 Java 的目录深度。Go 侧保持浅目录；不要新建 `types`、`utils`、`common`、`helpers` 这类无语义包。

### 注意事项

- `internal/trigger/http` 的包名与标准库 `net/http` 同名，外部包同时引用两者时需要起别名。
- `config/` 下只有 `config.example.yaml` 入库，派生出的 `config.dev.yaml` 等已被忽略；敏感项优先用环境变量注入。
- Go 工具链装在 `/opt/homebrew/bin`。如果 shell 报 `command not found: go`，先确认该路径在 `PATH` 中。

## Agent 智能中台

侨缘信使的 Agent 能力由一套独立的中台服务提供，整套能力统称为 **Qio Agent Platform（侨缘信使智能中台）**。它不是另一个面向用户的产品，而是为 Qio 主站提供写信辅助、文化知识问答、RAG 检索、多语言处理、内容运营和工具调用能力。

### 项目组成

| 目录 | 定位 | 主要技术栈 |
|------|------|------------|
| `qio-agent-service/` | Agent 运行时及管理 API | Java 17、Spring Boot 3.4.3、Spring AI、MyBatis、Maven |
| `qio-agent-console/` | Agent 内部配置与运营控制台 | React 18、TypeScript、Rsbuild、FlowGram、Semi UI、styled-components |

两个项目已经采用 Qio 统一命名：Maven 父项目为 `qio-agent-service`，子模块为 `qio-agent-*`，Agent Console 的 npm 包名为 `@qio/agent-console`。为保持现有代码和数据兼容，以下标识暂不迁移：

- Java 根包仍为 `com.chuppch`。
- MySQL schema 仍为 `ai-agent-station`。
- 已发布的 HTTP API 路径保持 `/api/v1/...`。

不要在普通重构中修改这些兼容标识。Java 包迁移、数据库改名或 API 版本升级必须作为独立迁移任务处理，并提供兼容与回滚方案。

### 调用边界

推荐调用链为：

```text
Qio 用户端
  -> Qio 业务后端（用户鉴权、业务编排、数据归属）
  -> Qio Agent Service（模型、Prompt、RAG、MCP、Agent 执行）
  -> 模型服务 / 向量数据库 / 外部工具

Qio Agent Console
  -> 仅供内部管理员配置和运营 Qio Agent Service
```

- 普通用户不应直接访问 Agent Console，也不应在用户界面接触模型、Advisor、MCP 等中台概念。
- Qio 的用户、信件、订单等核心业务数据归 Qio 业务后端管理；Agent Service 只保存 Agent 配置、知识库配置和执行所需数据。
- Qio 前端原则上不直接调用 Agent Service。需要面向用户开放的 AI 能力，应由 Qio 业务后端统一鉴权、限流、审计和包装。
- Agent Service 可以继续作为独立 Java 服务存在，不要求随 Qio v2 业务后端一起迁移到 Go。

### Agent Service 模块

`qio-agent-service/` 是 Maven 多模块工程：

- `qio-agent-api/` — 对外接口契约、DTO 和通用响应模型。
- `qio-agent-app/` — Spring Boot 启动模块与运行配置，入口为 `com.chuppch.Application`。
- `qio-agent-domain/` — Agent 核心领域逻辑，包括装配、调度、RAG，以及 fixed、auto、flow 执行策略。
- `qio-agent-infrastructure/` — MyBatis DAO、Repository 实现和数据库持久化。
- `qio-agent-trigger/` — HTTP 管理接口，以及后续任务、监听器等触发入口。
- `qio-agent-types/` — 跨模块公共类型和常量。

关键基础设施包括：

- MySQL：Agent、模型、Prompt、Advisor、MCP 和管理配置。
- PostgreSQL + pgvector：RAG 文档向量与相似度检索。
- Redis：缓存及运行期辅助状态。
- Spring AI：OpenAI 兼容模型接口、RAG Advisor、MCP 客户端和文档解析。

开发配置默认服务端口为 `8091`。环境连接信息集中在 `qio-agent-app/src/main/resources/application-*.yml`。这些文件可能包含敏感配置：不得在回复、日志、测试快照或文档中复述密钥与密码，也不要继续新增硬编码凭据；调整配置时优先迁移到环境变量或本地不入库配置。

### Agent Console 模块

`qio-agent-console/` 是内部管理端，当前包含：

- Agent 列表和可视化流程编排。
- 模型 API、模型实例、Prompt、Advisor 和 MCP 工具管理。
- RAG 知识库配置与文件上传。
- Agent 配置保存、装配与运行入口。
- 管理员登录和数据仪表盘。

前端请求集中在 `src/services/`，接口常量集中在 `src/config/api.ts`。当前 API 地址仍固定指向本机 `8091`，正式部署前应改为环境变量配置。

Agent Console 的 `build`、`test` 等部分 npm scripts 当前只是占位命令（`exit 0`），不能作为有效验证依据。不要因为命令返回成功就声称构建或测试通过。

## 常用命令

- 前端 v1 安装依赖：`cd qio-frontend && npm install`
- 前端 v1 本地开发：`cd qio-frontend && npm run serve`
- 前端 v1 生产构建：`cd qio-frontend && npm run build`
- 前端 v2：尚未初始化
- 后端 v2 整理依赖：`cd qio-backend-v2 && make tidy`
- 后端 v2 编译：`cd qio-backend-v2 && make build`
- 后端 v2 本地启动：`cd qio-backend-v2 && make run`
- 后端 v2 格式化与静态检查：`cd qio-backend-v2 && make fmt && make vet`
- 后端 v2 测试：`cd qio-backend-v2 && make test`
- 后端 v2 全部命令：`cd qio-backend-v2 && make help`
- Agent Service 编译：`cd qio-agent-service && mvn clean package -DskipTests`
- Agent Service 本地启动：`cd qio-agent-service && mvn -pl qio-agent-app -am spring-boot:run`
- Agent Console 安装依赖：`cd qio-agent-console && npm install`
- Agent Console 本地开发：`cd qio-agent-console && npm run dev`
- Agent Console 代码检查：`cd qio-agent-console && npm run lint`

## 目录边界

优先编辑这些目录或文件：

- `qio-backend-v2/` — v2 业务后端源码（Go）
- `qio-frontend-v2/` — v2 前端源码（Next.js，尚未初始化）
- `qio-agent-service/` — Agent 中台后端源码
- `qio-agent-console/src/` — Agent 中台管理端源码
- `docs/` — 项目文档
- `AGENTS.md` — 协作规范

`qio-backend/` 与 `qio-frontend/` 是 v1 遗留实现，作为迁移参照保留。除明确的 v1 修复任务外，不在其中新增功能。

请注意以下边界：

- 不要手改构建产物目录（`dist/`、`target/`、`build/`）
- 不要无故修改 `.env.*` 文件内容；如确有需要，只处理和当前任务直接相关的项
- 不要将 Agent Service 配置中的密钥、密码、生产地址或其他敏感值复制到新文件、测试数据或沟通输出中
- 不要把 Agent Console 的内部管理能力直接暴露到 Qio 用户端
- 不要修改 `.agents/` 下的 skill 源码，除非明确要求

## 代码风格

### 前端（React / TypeScript）

- 缩进：2 个空格
- 引号：单引号
- 语句结尾：分号
- 最大行宽：120
- 新增代码必须使用 TypeScript / TSX
- 样式以 Tailwind CSS 类名为主，优先在 TSX 中通过 `className` 直接组织
- 只有 Tailwind 难以表达或需覆盖第三方组件样式时，才允许补充 CSS Module / .scss
- 不为同一块 UI 额外维护一整套等价样式

### 后端（Go）

- 遵循 Go 官方格式化（gofmt / goimports）
- 命名遵循 Go 惯例（驼峰导出、小写私有）
- 错误处理显式 `if err != nil` 检查，不吞错误
- 包按业务域或功能命名，禁止 `types`、`utils`、`common`、`helpers`、`base` 这类无语义包名
- 目录层级保持浅，一到两层为宜，不做 `internal/services/user/handlers/http/v1/` 这类嵌套
- 每个包用 `doc.go` 承载包注释，说明职责与文件约定，作用等同于 Java 的 `package-info.java`
- 接口在使用方一侧声明，实现放 `internal/infrastructure`；不在 `domain` 内 import 具体存储或 Web 框架
- 同一个 struct 的方法可分散在多个文件，按业务主题拆分，避免出现 v1 那种巨型 Service
- 详细分层与依赖方向见「业务后端 v2」章节

### Agent Service（Java）

- 使用 Java 17，遵循现有 Maven 多模块和领域分层，不跨层直接访问 DAO。
- Controller/Trigger 负责协议转换，核心执行逻辑放在 Domain，持久化实现放在 Infrastructure。
- 保持显式异常处理和结构化日志，不在日志中输出 Prompt 全文、私人信件、访问令牌或数据库凭据。
- 新增模型、RAG、MCP 能力时优先复用 Spring AI 与现有装配、调度策略，不在 Controller 中直接拼装复杂 Agent。

### Agent Console（React / TypeScript）

- 延续现有 React 18、TypeScript、FlowGram、Semi UI 和 styled-components 技术栈。
- 通用接口调用放在 `src/services/`，API 地址和端点放在 `src/config/`，不要在页面组件中硬编码 URL。
- 可视化节点能力放在 `src/nodes/`，编辑器通用能力放在 `src/components/`、`src/plugins/` 或 `src/hooks/`。
- Agent Console 是内部工具，不要求沿用 Qio 面向用户端的品牌视觉组件体系，但应保持命名和业务概念一致。

## UI 组件约定

本项目前端使用 `shadcn/ui` 作为通用 UI 组件方案。

组件选型顺序如下：

1. 先用 shadcn/ui 已有组件（基于 Radix UI）
2. shadcn/ui 不完全满足时，在业务层做轻量封装
3. 动画交互使用 framer-motion
4. 图标使用 lucide-react

补充约定：

- 非必要不引入第二套通用 UI 组件库
- 自定义样式优先在 shadcn/ui 组件基础上覆盖，而不是从头重写

## 工作方式

- 先阅读相关实现，再进行修改
- 改动应聚焦当前任务，避免顺手清理无关代码
- 尽量小步修改，不顺手扩大 diff
- 只有高风险操作才需要确认（删除数据、覆盖大量文件、影响生产配置）
- 普通读文件、改代码、运行本地命令可直接执行

### 项目演进日志

- `PROJECT_LOG.md` 是仓库级项目演进记录。Agent 在开始跨项目或完整功能迭代前，应先阅读其中与当前任务相关的记录，避免基于过时的项目边界或方案继续开发。
- 只有当一次完整功能迭代已经实现并完成必要验证时，才在 `PROJECT_LOG.md` 追加记录。完整功能迭代是指形成了可独立说明的业务能力、端到端流程、架构调整、接口协议或数据迁移结果，而不是单个文件或单次代码修改。
- 普通缺陷修复、空值判断、格式调整、局部重构、依赖升级、文件重命名、探索分析、代码评审以及尚未完成的中间步骤，不要求更新 `PROJECT_LOG.md`。
- 日志记录应简洁说明日期、状态、背景、最终实现、跨项目影响、验证结果和必要的后续工作；不要逐条复制提交记录，也不要写入密钥、密码或私人数据。
- 如果一次迭代改变了之前的方案，应追加新记录说明变更原因和新结论，不覆盖或改写原有历史记录。

## 验证要求

改动后至少执行与改动匹配的最小验证：

- 前端页面/交互改动：优先运行开发服务器确认
- 后端 v2（Go）改动：至少执行 `make build` 与 `make fmt`，有测试则跑 `make test`
- 后端接口改动：优先编译通过，有条件则跑测试
- Agent Service 改动：至少执行对应 Maven 模块编译；涉及模型、数据库或 RAG 时说明所需外部依赖是否可用
- Agent Console 改动：至少执行 `npm run lint` 并使用开发服务器验证；当前占位的 `npm run build`、`npm test` 不算有效验证
- 构建/配置改动：优先运行构建命令
- 仅文档改动：检查事实是否与仓库现状一致

如果命令失败：

- 说明失败命令和关键报错
- 给出下一步建议
- 不要只汇报"已修改"而不说明验证结果

## 提交规范

提交信息使用 cz-emoji shortcode 风格。**提交前必须先读 `.agents/skills/emoji-commit/SKILL.md`**，以该技能及其 `references/` 下的文档为准。该技能未在 `.kiro/` 下建立软链，不会自动出现在可用技能列表中，需要主动按路径读取。

关键参考文件：

- `.agents/skills/emoji-commit/SKILL.md` — 路由判断与最小硬约束
- `.agents/skills/emoji-commit/references/fex-conventional-commits.md` — 完整 header / body / footer 规范
- `.agents/skills/emoji-commit/references/cz-emoji-types.md` — emoji 类型对照表

最小硬约束（细节仍以技能文档为准）：

- Header 骨架：`:emoji: subject`、`:emoji: ! subject`、`:emoji: (scope) subject`、`:emoji: (scope) ! subject`
- emoji 必须用 shortcode（如 `:bug:`），不得用 Unicode emoji；docs 类用 `:memo:`
- body 可选，存在时用 bullet 风格，与 header 之间保留一个空行
- `AI-Co-Authored-By:` 是必需 footer，且必须是提交信息的最后一行
- 禁止 `Co-authored-by` / `Co-Authored-By` 等其他 co-author trailer
- 提交前先看 `git status --short`，据此判断走单条提交还是分批提交

本项目额外约定（在技能的「推荐写法」之上收紧为强制要求）：

- subject 与 body 一律使用英文，不使用中文
- subject 使用祈使语气，如 `add` / `fix` / `document` / `move`，不用第三人称或过去式
- subject 简短明确，`scope` 使用目录或模块短名，如 `backend-v2`、`agent-service`、`agents`
- 历史提交中存在中文 subject，属于此约定生效前的遗留，不回溯修改

示例：

```text
:tada: (backend-v2) init go project scaffold

- add trigger / app / domain / infrastructure layers with dependency direction
- split business domains by v1 controllers

AI-Co-Authored-By: Kiro
```

## 沟通输出

- 工作中提供简短进度更新
- 结束时说明：改了什么、验证了什么、是否有剩余风险或可选下一步
