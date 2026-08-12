// Package infra 是基础设施适配层的根包。
//
// 各业务域只声明所需的接口（repository / 能力接口），实现统一放在本包下按技术
// 分类的子包中。依赖方向始终由 infra 指向业务域，不允许反向依赖。
//
// 子包划分：
//   - mysql     GORM 连接与各域 repository 实现
//   - redis     缓存客户端与缓存实现
//   - rabbitmq  消息生产与消费
//   - minio     对象存储实现
//   - mail      邮件投递实现
//   - agentsvc  Qio Agent Service 出站客户端
package infra
