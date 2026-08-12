package friend

import "errors"

// 好友域的错误值。
//
// 对应 v1 的 FriendException、FriendNotExistsException 与其消息码。
var (
	// ErrNotFound 好友关系不存在。
	// 对应 FriendNotExistsException
	ErrNotFound = errors.New("friend not found")

	// ErrNoPermission 当前用户无权处理该好友申请。
	// 对应 friend.friendRequest.no.permission
	ErrNoPermission = errors.New("no permission on friend request")

	// ErrRequestUnavailable 好友申请不存在或已被处理。
	//
	// v1 把「不存在」与「已处理」合并为同一条消息 friend.friendRequest.empty，
	// 这里保持合并，不细分两种情形。
	ErrRequestUnavailable = errors.New("friend request unavailable")

	// ErrAddressNotFound 好友地址簿中不存在该地址。
	// 对应 friend.address.not.exists
	ErrAddressNotFound = errors.New("friend address not found")

	// ErrDefaultAddressUndeletable 默认地址不允许删除。
	// 对应 friend.address.default.not.delete
	ErrDefaultAddressUndeletable = errors.New("default friend address is undeletable")

	// ErrAddFailed 建立好友关系失败。
	// 对应 friend.add.failed
	ErrAddFailed = errors.New("add friend failed")

	// ErrCreateRequestFailed 创建好友申请失败。
	// 对应 friend.create.Request.failed
	ErrCreateRequestFailed = errors.New("create friend request failed")

	// ErrSourceNotFound 好友申请的来源漂流瓶或信件不存在。
	//
	// 对应 friend.bottle.not.exists。v1 的判空写成 `bottle == null || letter == null`，
	// 而两个变量都被初始化为空对象、只有走进对应分支才可能为 null，实际语义混乱，
	// 记录在 docs/TODO-migration.md。
	ErrSourceNotFound = errors.New("friend request source not found")
)
