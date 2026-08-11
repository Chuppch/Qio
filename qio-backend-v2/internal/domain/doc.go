// Package domain 是业务域的根包。
//
// 每个子包是一个业务域，只承载业务规则与领域模型，不依赖 Web 框架、
// 不依赖具体存储实现。域内文件约定：
//
//   - service.go     业务逻辑
//   - repository.go  数据访问接口定义，实现由 internal/infra 提供
//   - model.go       领域模型
//
// 接口数量多且主题分组明显时，按主题拆成 service_xxx.go，同一个域的所有文件
// 共享同一个 package，不必新建子目录。
//
// 域复杂到单个域目录难以维护时，再在域内拆 application/ 等子包；不提前分层。
//
// 依赖方向：
//
//	transport -> app -> domain <- infra
//
// domain 不得 import transport、app、infra 中的任何包。
package domain
