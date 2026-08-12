// Package mysql 提供 GORM 连接管理、持久化对象（PO）与各业务域 repository 接口的实现。
//
// # 文件命名
//
//	db.go        连接与连接池
//	base.go      共享的审计字段
//	json.go      JSON 列的存储结构与解析
//	po_<域>.go    表映射与 PO ↔ 领域模型转换
//	repo_<域>.go  repository 接口实现，SQL 只出现在这类文件中
//
// 用前缀而非后缀，使 po_* 与 repo_* 在按名排序时各自聚成一组。
//
// # PO 约定
//
// PO 是数据库表的映射对象，一张表一个类型，命名为 <表名>PO 且不导出。
// 不导出是刻意的：业务域与传输层在编译期就拿不到 PO 类型，越层直接操作数据库
// 在语言层面不可能发生。这一点比 Java 的 public PO 更严格。
//
// 一个 repo 对应一个聚合，不是一个表。字典表（avatar、country）没有聚合根身份，
// 共用 repo_dict.go；shop 域的七类道具同属商品聚合，将来共用 repo_shop.go。
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
//   - user.status、user.del_flag 用 varchar 存枚举
//   - id 列在不同表中混用 bigint、int、bigint unsigned
package mysql
