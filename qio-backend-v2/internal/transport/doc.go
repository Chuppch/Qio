// Package transport 是入站传输层的根包，对应 Agent Service 的 qio-agent-trigger 模块。
//
// 子包划分：
//   - http        HTTP 处理器与路由注册
//   - middleware  HTTP 中间件
//   - dto         对外接口契约（请求体、响应体）
//
// 职责边界：只做协议转换（DTO 与领域模型互转）、参数校验、错误码映射。
// 业务规则一律不写在本层。跨域动作调用 internal/app，单域动作调用对应域的 service。
package transport
