// Package dto 定义对外接口契约：请求体、响应体与查询参数。
//
// 对应 Agent Service 的 qio-agent-api/dto。
//
// 存在的意义是把接口契约与领域模型解耦：领域模型可以自由演进而不破坏已发布的
// 接口，领域模型中的内部字段也不会意外泄漏给客户端。
//
// 约定：
//   - 一个业务域一个文件
//   - 命名统一为 XxxRequest / XxxResponse / XxxQuery
//   - 只放数据结构与校验标签，不写业务逻辑
//   - DTO 与领域模型的转换写在 transport/http 的处理器中
package dto
