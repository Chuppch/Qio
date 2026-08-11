// Package agentsvc 是 Qio Agent Service 的出站 HTTP 客户端。
//
// 只负责协议层面的请求构造、流式响应解析、超时与重试。业务侧的鉴权、
// 配额与审计在 internal/ai 中完成。
//
// Agent Service 的地址与凭据来自配置，不在代码中写死。
package agentsvc
