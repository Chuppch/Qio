# 需求文档：架构迁移与中间件升级

## 1. 目标与范围

### 目标

完成 Qio v2.0.0 的技术底座切换，将前后端从 v1 技术栈迁移至目标技术栈，同时将中间件升级至当前稳定版本。

### 范围

包含：
- 前端项目从 Vue 2 迁移至 Next.js (React 18 + TypeScript)
- 后端项目从 Java/Spring Boot 迁移至 Go/Gin
- 中间件版本升级（MySQL、Redis、RabbitMQ、Minio）
- 开发环境 Docker Compose 配置更新

不包含：
- 业务功能新增
- UI 重设计（仅做基础页面骨架迁移）
- 生产环境部署方案（后续单独出文档）

## 2. 现状描述

### 技术栈

- 前端：Vue 2 + Element UI + Webpack（26 个页面，~17,800 行）
- 后端：Spring Boot 3.3.3 + MyBatis-Plus + Java（196 个文件，~17,400 行）
- 详细模块分析见 `docs/architecture/module-analysis.md`

### 中间件版本

| 中间件 | 当前版本 |
|--------|----------|
| MySQL | 8.3 |
| Redis | 7 |
| RabbitMQ | 3.x |
| Minio | 未锁定版本 |

## 3. 目标架构

### 中间件目标版本

| 中间件 | 目标版本 |
|--------|----------|
| MySQL | 8.4 (LTS) |
| Redis | 8 |
| RabbitMQ | 4 |
| Minio | latest（关注社区动态，必要时评估替代方案） | // todo

### 目标架构结构

// todo

## 4. 迁移步骤

### Phase 1：后端骨架搭建

1. 初始化 Go module，搭建 Gin 项目结构
2. 配置加载（Viper 读取 YAML）
3. 接入 MySQL (GORM) + Redis (go-redis)
4. 实现 JWT 中间件
5. 完成用户模块（注册、登录、鉴权）作为验证

**产出物**：后端能启动，能注册登录，JWT 鉴权通过

### Phase 2：后端核心业务迁移

1. 信件模块（写信、寄信、收信）
2. WebSocket 通信（信件状态推送、AI 流式输出）
3. RabbitMQ 接入（AI 消息队列）
4. Minio 文件上传
5. 漂流瓶、好友、商店、游戏等剩余模块

**产出物**：所有 v1 API 在 Go 后端均有对应实现

### Phase 3：前端骨架搭建

1. 初始化 Next.js 项目 + TypeScript + Tailwind + shadcn/ui
2. 配置路由结构（App Router）
3. 封装 API 请求层（fetch / axios）
4. 实现登录注册页面作为验证

**产出物**：前端能启动，能完成登录注册闭环

### Phase 4：前端页面迁移

1. 首页展示模块
2. 信件模块（写信、收信、漂流瓶）
3. AI 对话（WebSocket + 流式渲染）
4. 游戏、商店、营销等剩余页面
5. 个人中心

**产出物**：所有 v1 页面在 React 前端均有对应实现

### Phase 5：中间件升级 & 联调

1. 更新 docker-compose.yml 中间件版本
2. 验证新版本中间件与 Go 后端兼容性
3. 前后端联调，跑通完整业务流程
4. 清理 v1 旧代码（或归档到独立分支）

**产出物**：新版中间件 + Go 后端 + React 前端完整跑通

## 5. 验收标准

### 阶段验收

| 阶段 | 验收条件 |
|------|----------|
| Phase 1 | `go build` 通过，服务启动，登录接口返回 JWT |
| Phase 2 | 所有 v1 接口有对应 Go 实现，Postman/curl 验证通过 |
| Phase 3 | `npm run dev` 启动，登录注册页面可交互 |
| Phase 4 | 所有页面可访问，UI 基本还原 |
| Phase 5 | `docker compose up` 一键启动，完整业务流程跑通 |

### 最终验收

- 前后端能通过 docker compose 一键启动
- 用户注册 → 登录 → 写信 → 寄信 → 收信完整流程通过
- AI 对话流式输出正常
- 漂流瓶收发正常
- 无 v1 技术栈残留依赖

## 6. 风险与待定项

| 风险 | 影响 | 应对 |
|------|------|------|
| 信件渲染方案未定（Graphics2D → ?） | 阻塞信件模块完成 | Phase 2 中优先调研，可选 Go `gg` 库或前端 Canvas 渲染 |
| Minio 官方仓库已归档 | 长期维护风险 | 短期继续使用，观察 openminio 社区活跃度，必要时切换云 OSS |
| v1 数据库 schema 可能需要调整 | 影响 GORM model 定义 | Phase 1 先 1:1 映射现有表结构，后续优化单独迭代 |
| WebSocket 从 Spring 迁移到 Go | 实现方式差异大 | 使用 gorilla/websocket 或 nhooyr/websocket，Phase 2 重点验证 |
