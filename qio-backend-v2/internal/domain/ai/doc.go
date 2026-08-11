// Package ai 承载面向用户的 AI 能力：写信辅助、文化问答、角色对话。
//
// 本域不自己实现模型调用，而是 Qio Agent Service 的业务侧门面：负责配额、
// 用量统计、会话归属与审计留痕，再经 internal/infra/agentsvc 转调中台。
// 模型、Prompt、Advisor、MCP、RAG 等中台概念不得泄漏到用户端。
//
// 域内文件约定见 internal/domain 包注释。
package ai
