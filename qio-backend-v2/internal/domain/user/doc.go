// Package user 承载用户域：注册登录、账号信息、地址簿、签到与每日任务、
// 虚拟货币余额、背包。
//
// 文件划分：
//
//   - model.go       User 聚合根、其值对象 Address 与 OwnedItem
//   - growth.go      Task、SignAward，独立于 User 的聚合，落 Redis
//   - validation.go  账号格式校验，无接收者的包级函数
//   - errors.go      本域的错误值，逐条对应 v1 的异常类
//   - repository.go  三个数据访问接口
//
// 三类内容不在本域：
//
//   - 头像与国家字典：多域共用，见 internal/domain/dict
//   - 好友关系：v1 把五个好友方法放在 UserServiceImpl 中直接操作 friend 表，
//     v2 归 internal/domain/friend
//   - 用户统计（收信数、好友数）：跨域聚合，由 internal/app 组合 letter 与
//     friend 两个域的计数得出
package user
