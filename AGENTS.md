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

## 常用命令

- 前端安装依赖：`cd qio-frontend && npm install`
- 前端本地开发：`cd qio-frontend && npm run serve`
- 前端生产构建：`cd qio-frontend && npm run build`
- 后端编译运行：待 Go 项目初始化后补充

## 目录边界

优先编辑这些目录或文件：

- `qio-frontend/src/` — 前端源码
- `qio-backend/` — 后端源码
- `docs/` — 项目文档
- `AGENTS.md` — 协作规范

请注意以下边界：

- 不要手改构建产物目录（`dist/`、`target/`、`build/`）
- 不要无故修改 `.env.*` 文件内容；如确有需要，只处理和当前任务直接相关的项
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
