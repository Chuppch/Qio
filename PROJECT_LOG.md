# Qio Project Log

本文件记录 Qio 仓库中具有长期价值的项目演进信息，包括架构调整、项目边界、跨项目依赖、接口与数据迁移决策，以及阶段性实施状态。

它与 `AGENTS.md` 同级，但职责不同：

- `AGENTS.md`：记录当前有效的开发规范、项目边界和协作要求。
- `PROJECT_LOG.md`：记录项目为何演进到当前状态，以及各项目之间需要同步关注的变化。
- Git 提交历史：记录具体代码修改，不在本文件重复维护普通代码变更。

## 维护约定

- 仅记录架构、技术栈、接口、数据库、消息协议、项目边界和重要功能迁移等长期变化。
- 普通空值修复、格式调整、局部重命名等内容不单独记录。
- 涉及多个项目的修改，应说明上游、下游和兼容影响。
- 每条记录应包含日期、状态、背景、调整内容、影响范围和后续工作。
- 状态统一使用：`规划中`、`进行中`、`已完成`、`已暂停`、`已废弃`。
- 相关方案发生变化时，追加新记录说明，不覆盖历史结论。
- 不在本文件中记录密钥、密码、私人数据和生产环境敏感地址。

## 项目全景基线

| 项目 | 定位 | 当前状态 | 主要技术栈 |
|---|---|---|---|
| `qio-frontend/` | Qio v1 用户前端，作为迁移参考保留 | 维护模式 | Vue 2、Element UI |
| `qio-backend/` | Qio v1 业务后端，作为迁移参考保留 | 维护模式 | Java、Spring Boot |
| `qio-frontend-v2/` | Qio v2 用户前端 | 尚未初始化 | Next.js、React、TypeScript |
| `qio-backend-v2/` | Qio v2 核心业务后端 | 骨架阶段 | Go、Gin、GORM |
| `qio-agent-service/` | Qio Agent 运行时与管理服务 | 持续完善 | Java 17、Spring Boot、Spring AI |
| `qio-agent-console/` | Agent 内部配置与运营控制台 | 持续完善 | React、TypeScript、FlowGram |

## 当前系统边界

推荐的用户侧调用链：

```text
Qio 用户端
  -> Qio 业务后端
  -> Qio Agent Service
  -> 模型服务 / RAG / MCP / 外部工具
```

推荐的管理侧调用链：

```text
Qio Agent Console
  -> Qio Agent Service
```

边界约定：

- Qio 前端不直接调用 Agent Service。
- 用户、信件、订单、好友等业务数据归 Qio 业务后端管理。
- Agent Service 管理模型、Prompt、Advisor、RAG、MCP 和 Agent 执行能力。
- Agent Console 仅面向内部管理员，不向普通用户暴露。
- Qio v2 采用单体架构优先，不为技术形式提前拆分微服务。

---

## 2026-08-11：建立项目演进日志基线

**状态：已完成**

### 背景

Qio 已同时包含 v1 前后端、v2 重构项目和独立 Agent 中台。多个项目之间存在接口、数据和调用边界依赖，仅依靠 Git 历史和单个项目说明，难以让开发者与 Agent 快速理解整体演进过程。

### 调整内容

- 在仓库根目录建立 `PROJECT_LOG.md`。
- 以当前仓库状态作为后续演进记录的统一基线。
- 明确稳定规范归属 `AGENTS.md`，演进历史归属 `PROJECT_LOG.md`。

### 影响范围

- 后续跨项目、架构、接口和数据迁移调整应同步更新本文件。
- 各子项目仍可保留自身 README，不在本文件复制具体开发说明。

### 后续工作

- 在后续提交中持续追加真实演进记录。
- 当单文件内容过大时，再拆分为 `docs/project-evolution/` 下的分项目记录。

---

## 2026-08-11：Qio v2 技术栈与业务后端架构基线

**状态：进行中**

### 背景

Qio v2 将在保留侨批文化核心体验的基础上，对 v1 前后端进行系统性重构，降低历史技术债，并为新的品牌定位、写信主流程和 Agent 能力接入提供更清晰的工程基础。

### 已确定方案

- 用户前端从 Vue 2 迁移到 Next.js、React 和 TypeScript。
- 业务后端从 Java/Spring Boot 迁移到 Go/Gin。
- `qio-backend-v2` 采用领域化单体和端口适配器设计。
- 后端按 `user`、`letter`、`bottle`、`friend`、`shop`、`explore`、`ai` 等业务域组织。
- 入站协议放在 `internal/transport`，核心业务放在 `internal/domain`，基础设施实现放在 `internal/infra`。
- 跨多个业务域的用例放在 `internal/app`，避免领域 Service 重新演变为巨型服务。

### 当前完成度

- `qio-backend-v2` 已完成目录和职责骨架。
- 启动、路由、领域模型、Repository、基础设施和业务用例尚未实现。
- 当前骨架不能代表业务闭环已经跑通。

### 跨项目影响

- `qio-frontend-v2` 后续应以 Go 后端定义的业务 API 为唯一用户侧接口。
- `qio-backend-v2` 负责用户鉴权、业务编排、限流、审计和 Agent 输出加工。
- `qio-agent-service` 保持独立 Java 服务，不随业务后端一起迁移到 Go。

### 下一阶段建议

- 优先围绕“写信”完成一条纵向业务闭环。
- 在实现其他业务域前，先确定认证、错误模型、事务边界和 API 版本规范。
- 根据真实业务需求逐步实现目录内容，避免继续扩充空骨架。

---

## 2026-08-11：Agent 中台接入边界与运行时状态

**状态：进行中**

### 已确定方案

- Agent 中台由 `qio-agent-service` 和 `qio-agent-console` 组成。
- Qio 面向用户的 AI 能力由业务后端统一转接，不允许前端直接调用 Agent Service。
- Agent Service 已具备 Fixed、Auto、Flow 三类执行策略及动态 Bean 装配基础。

### 当前实现状态

- Agent Service 已补齐部分执行责任树 Bean 注册和空值保护。
- Agent 执行引擎存在，但当前没有面向 Qio 业务后端的正式执行入口。
- 数据库配置变化后，Spring 运行时 Bean 不会自动同步刷新。
- 当前不增加手动装配 HTTP 接口，也不在应用启动时自动全量装配。

### 后续方向

- 配置保存并成功提交事务后，通过领域事件触发受影响 Bean 的定向热更新。
- 正式接入 Qio 后端时，再确定 Agent 执行协议、鉴权、超时、错误模型和流式返回格式。
- 优先考虑按需装配或版本化运行时注册，避免替换 Bean 时影响正在执行的请求。
