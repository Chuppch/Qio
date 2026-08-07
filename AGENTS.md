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
| Spring Boot 3.3.3 | Gin (Go 1.22+) |
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

- 前端安装依赖：`cd qio-frontend && npm install`
- 前端本地开发：`cd qio-frontend && npm run serve`
- 前端生产构建：`cd qio-frontend && npm run build`
- 后端编译运行：待 Go 项目初始化后补充
- Agent Service 编译：`cd qio-agent-service && mvn clean package -DskipTests`
- Agent Service 本地启动：`cd qio-agent-service && mvn -pl qio-agent-app -am spring-boot:run`
- Agent Console 安装依赖：`cd qio-agent-console && npm install`
- Agent Console 本地开发：`cd qio-agent-console && npm run dev`
- Agent Console 代码检查：`cd qio-agent-console && npm run lint`

## 目录边界

优先编辑这些目录或文件：

- `qio-frontend/src/` — 前端源码
- `qio-backend/` — 后端源码
- `qio-agent-service/` — Agent 中台后端源码
- `qio-agent-console/src/` — Agent 中台管理端源码
- `docs/` — 项目文档
- `AGENTS.md` — 协作规范

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
- 项目结构待初始化后补充

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

## 验证要求

改动后至少执行与改动匹配的最小验证：

- 前端页面/交互改动：优先运行开发服务器确认
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

提交信息使用 cz-emoji shortcode 风格：

- Header 格式：`:emoji: subject` 或 `:emoji: (scope) subject`
- 如果仓库内 `emoji-commit` skill 可用，提交时优先参考该技能约定

## 沟通输出

- 工作中提供简短进度更新
- 结束时说明：改了什么、验证了什么、是否有剩余风险或可选下一步
