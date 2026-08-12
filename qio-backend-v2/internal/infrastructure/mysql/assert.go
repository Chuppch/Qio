package mysql

import (
	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/bottle"
	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/dict"
	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/shop"
	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/user"
)

// 编译期断言：各仓储实现必须完整覆盖对应的领域接口。
//
// 接口新增方法而实现未跟上时，这里会直接编译失败，不必等到运行期。
var (
	_ user.Repository   = (*userRepository)(nil)
	_ dict.Repository   = (*dictRepository)(nil)
	_ shop.Repository   = (*shopRepository)(nil)
	_ bottle.Repository = (*bottleRepository)(nil)
)
