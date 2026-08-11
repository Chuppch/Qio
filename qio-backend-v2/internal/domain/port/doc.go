// Package port 存放跨域的能力接口（出站端口）。
//
// 这里只定义接口，不含业务规则，也不含任何实现。实现全部放在 internal/infra
// 下对应的适配器包中。业务域与 app 层依赖本包的接口，从而与外部系统解耦。
//
// 对应 Agent Service 中的 domain/agent/adapter/port。
//
// 说明：各业务域自己的数据访问接口（repository）仍留在域内，不放本包。
// 本包只收纳被多个域共用、且指向外部系统的能力，如通知投递、对象存储。
package port
