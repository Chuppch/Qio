// Package main 是 Qio 业务后端的启动入口，对应 Agent Service 的 qio-agent-app 模块。
//
// 这里只做三件事：加载配置、组装依赖、启动 HTTP 服务。
// 任何业务逻辑都不应写在本包内。
package main

func main() {
	// TODO: 加载配置
	// TODO: 初始化基础设施（mysql / redis / rabbitmq / minio / mail / agentsvc）
	// TODO: 组装各域 service 与 app 用例，注入 infra 实现
	// TODO: 注册路由与中间件，启动 HTTP 服务并支持优雅关闭
}
