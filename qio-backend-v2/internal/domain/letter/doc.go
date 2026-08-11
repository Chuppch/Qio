// Package letter 承载信件域：信件本身的规则与投递状态流转、收信箱、已读标记。
//
// 本域是 Qio 的核心业务，预计最先变复杂。跨域动作（发信涉及扣费、封面生成、
// 好友关系判断、邮件通知）放在 internal/app 编排，本域只保留信件自身的规则。
//
// 域内文件约定见 internal/domain 包注释。
package letter
