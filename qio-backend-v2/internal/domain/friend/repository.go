package friend

import "context"

// Repository 是好友域的数据访问接口，实现由 internal/infrastructure/mysql 提供。
//
// 方法集对应 v1 中对 friendMapper 与 friendRequestMapper 的全部生效调用，
// 这些调用散落在 FriendServiceImpl、UserServiceImpl、LetterServiceImpl、
// BottleServiceImpl 四个类里。FriendServiceImpl 中两段被 /* */ 注释掉的
// sendFriendRequest 不计入。
//
// 好友与好友申请是两个聚合，但同属好友域、同落 MySQL，且「同意申请」会同时写
// 两者，因此收在同一个接口内。
//
// v1 从未删除过 friend 或 friend_request 记录，也没有分页与投影，
// 因此这里不提供对应方法。
type Repository interface {
	// ---- 好友关系 ----

	// ListByOwner 查询指定用户的全部好友。
	//
	// 对应 v1 getMyFriends 中 owning_id 的等值查询。无排序、无分页。
	ListByOwner(ctx context.Context, ownerID int64) ([]*Friend, error)

	// FindByID 按主键查询好友记录，不存在时返回 ErrNotFound。
	//
	// 对应 v1 updateFriendRemark 中的 selectById。该处未校验归属，
	// 归属校验应由调用方用 Friend.OwnedBy 补上。
	FindByID(ctx context.Context, id int64) (*Friend, error)

	// FindByIDAndOwner 按主键与归属者查询好友记录，不存在时返回 ErrNotFound。
	//
	// 对应 v1 getFriendAddress、setFriendDefaultAddress、deleteFriendAddress
	// 中 (id, owning_id) 的组合查询。
	FindByIDAndOwner(ctx context.Context, id, ownerID int64) (*Friend, error)

	// FindByUserAndOwner 按好友的用户 ID 与归属者查询好友记录，
	// 不存在时返回 ErrNotFound。
	//
	// 对应 v1 sendLetterPre 中「检查收件人是否已是好友」的查询。
	// 注意查的是 user_id 而非主键。
	FindByUserAndOwner(ctx context.Context, userID, ownerID int64) (*Friend, error)

	// Create 新增一条好友关系，成功后回填 ID。
	//
	// v1 建立双向关系时连续调用两次，两次之间不在事务内，
	// 且表上没有 (owning_id, user_id) 唯一约束，重复调用会产生重复记录。
	Create(ctx context.Context, f *Friend) error

	// Update 更新好友记录的资料快照、地址簿与备注。
	Update(ctx context.Context, f *Friend) error

	// UpdateAll 批量更新好友记录，语义同 Update。
	//
	// 对应 v1 getMyFriends 末尾的 updateById(List)——一次读操作触发 N 条 UPDATE。
	// 该调用的问题记录在 docs/TODO-migration.md，此处只提供等价能力。
	UpdateAll(ctx context.Context, fs []*Friend) error

	// ---- 好友申请 ----

	// ListPendingRequests 查询指定用户收到的待处理申请。
	//
	// 对应 v1 ProcessingFriendRequests 中 (receiver_id, status = 0) 的查询。
	ListPendingRequests(ctx context.Context, receiverID int64) ([]*Request, error)

	// FindRequestByID 按主键查询好友申请，不存在时返回 ErrRequestUnavailable。
	//
	// v1 把「申请不存在」与「申请已处理」用同一条消息呈现，因此这里也复用
	// 同一个错误值。
	FindRequestByID(ctx context.Context, id int64) (*Request, error)

	// CountRequests 统计两个用户之间指定方向的申请数量。
	//
	// 对应 v1 readLetter 中的查重。不带状态条件，即历史上存在过任意状态的
	// 同向申请都会被计入，从而不再自动新建。
	CountRequests(ctx context.Context, senderID, receiverID int64) (int64, error)

	// CreateRequest 新增一条好友申请，成功后回填 ID。
	//
	// 若 Request.CreatedAt 非零则按该值落库，对应 v1 readLetter 用信件投递
	// 时间覆盖创建时间的做法。
	CreateRequest(ctx context.Context, r *Request) error

	// UpdateRequest 更新好友申请的处理状态。
	UpdateRequest(ctx context.Context, r *Request) error
}
