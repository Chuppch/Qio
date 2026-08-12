package middleware

// JWT 认证。
//
// claim 名称迁移自 v1 的 JwtClaimsConstant，保持一致以便灰度期间新旧版本
// 签发的 token 可以互认。
//
// v1 的 JwtProperties 支持 userSecretKey / userTtl / userTokenName 三项配置
// （另有一组被注释掉的 admin 配置，未启用，不迁移）。v2 中密钥与有效期由
// internal/config 提供，本包只负责签发与校验。

// JWT 载荷中的 claim 名称。
const (
	ClaimUserID   = "userId"
	ClaimEmail    = "email"
	ClaimUsername = "username"
	ClaimName     = "name"
)

// TODO: 实现签发与校验，密钥与有效期从 config 注入，不在本包内写死。
