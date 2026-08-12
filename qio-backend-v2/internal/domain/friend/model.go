package friend

// RequestStatus 是好友申请的处理状态。
//
// 取值沿用 v1，避免存量数据迁移。
type RequestStatus int

const (
	RequestPending  RequestStatus = 0 // 待处理
	RequestAccepted RequestStatus = 1 // 已同意
	RequestRejected RequestStatus = 2 // 已拒绝
)

// Pending 表示申请是否仍待处理。
func (s RequestStatus) Pending() bool {
	return s == RequestPending
}
