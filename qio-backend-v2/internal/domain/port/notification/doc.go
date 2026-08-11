// Package notification 提供通知能力：邮件、站内消息等。
//
// 本包只定义能力接口与模板组织方式，具体投递实现放在 internal/infra 下。
// 邮件正文模板不应以字符串拼接的方式写在业务代码里。
package notification
