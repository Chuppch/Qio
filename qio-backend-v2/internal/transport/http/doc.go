// Package http 存放 HTTP 处理器与路由注册，对应 Agent Service 的 trigger/http。
//
// 文件约定：
//   - router.go   全部路由集中注册，一处即可看清接口全貌
//   - <域>.go     该域的处理器，与业务域一一对应
//
// 处理器职责：解析并校验请求 DTO、调用 app 或域 service、把领域模型转成响应 DTO、
// 映射错误码。不写业务规则，不直接访问数据库。
//
// 命名提示：本包名与标准库 net/http 同名，外部包同时引用两者时需要起别名。
package http
