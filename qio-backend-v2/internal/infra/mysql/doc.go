// Package mysql 提供 GORM 连接管理、持久化对象（PO）与各业务域 repository 接口的实现。
//
// # PO 约定
//
// PO 是数据库表的映射对象，一张表一个类型，命名为 <表名>PO 且不导出。
// 不导出是刻意的：业务域与传输层在编译期就拿不到 PO 类型，越层直接操作数据库
// 在语言层面不可能发生。这一点比 Java 的 public PO 更严格。
//
// PO 与领域模型之间通过本包内的转换函数互转：
//
//	func (p letterPO) toDomain() letter.Letter
//	func letterPOFrom(l letter.Letter) letterPO
//
// # 文件划分
//
// 按业务域组织，一个域一个文件；跨域共用的地理字典单独成文件：
//
//	letter.go    letter
//	user.go      user、avatar
//	friend.go    friend、friend_request
//	bottle.go    bottle
//	shop.go      paper、font、font_color、signet、function_card、font_paper、commodity
//	explore.go   questions、question_user_status
//	country.go   country
//
// # 当前 schema 的已知问题
//
// PO 如实映射 v1 现有表结构，以便平滑接入存量数据。下列问题在 PO 注释中逐条
// 标注，属于待偿还的技术债，需要单独的迁移任务处理，不在建模阶段擅自改动：
//
//   - paper 表没有主键，且 id 可为 NULL
//   - font_paper 表主键无自增
//   - create_user / update_user 在不同表中混用 bigint 与 varchar(50)
//   - user、letter、friend、bottle、friend_request 用 JSON 列存关系数据
//   - letter.speed_rate、letter.reduce_time 用 varchar 存数值
//   - id 列在不同表中混用 bigint、int、bigint unsigned
package mysql
