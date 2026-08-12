package bottle

import "errors"

// 漂流瓶域的错误值，对应 v1 的 BottleException 与其消息码。
var (
	// ErrNotFound 漂流瓶不存在
	ErrNotFound = errors.New("bottle not found")

	// ErrNoneAvailable 当前没有可捞的漂流瓶。
	// 对应 bottle.Database.Bottles.empty
	ErrNoneAvailable = errors.New("no available bottle")

	// ErrNotPickedByUser 该用户当前没有捞起中的漂流瓶，无法扔回或加好友。
	// 对应 bottle.not.accord.condition
	ErrNotPickedByUser = errors.New("no bottle picked by user")

	// ErrGenerateFailed 漂流瓶生成失败
	ErrGenerateFailed = errors.New("bottle generate failed")
)
