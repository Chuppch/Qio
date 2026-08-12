package response

// 统一响应结构。
//
// 迁移自 v1 的 AjaxResult，但不沿用其实现方式：
//
//   - v1 是 `AjaxResult extends HashMap<String, Object>`，字段靠字符串键写入，
//     没有类型约束，也无法生成准确的接口文档。这里改为泛型结构体。
//   - v1 的 code 直接复用 HTTP 状态码，并额外引入了非标准的 601 表示警告。
//     业务错误码与 HTTP 状态码混用会让客户端难以区分「传输失败」和「业务拒绝」，
//     这里保留 v1 的取值以兼容现有前端，但语义上定义为业务码，与 HTTP 状态码解耦：
//     响应体里的 Code 表示业务结果，HTTP 状态码由 transport 层单独决定。

// Code 是业务结果码。
type Code int

const (
	CodeSuccess Code = 200 // 成功
	CodeError   Code = 500 // 失败
	CodeWarn    Code = 601 // 警告，操作未执行但不视为异常
)

// Body 是所有接口的统一响应体。
type Body[T any] struct {
	Code Code   `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data,omitempty"`
}

// Page 是分页数据的载荷，配合 Body 使用。
type Page[T any] struct {
	Total int64 `json:"total"`
	List  []T   `json:"list"`
}

// OK 构造成功响应。
func OK[T any](data T) Body[T] {
	return Body[T]{Code: CodeSuccess, Msg: "操作成功", Data: data}
}

// OKWithMsg 构造带自定义提示的成功响应。
func OKWithMsg[T any](msg string, data T) Body[T] {
	return Body[T]{Code: CodeSuccess, Msg: msg, Data: data}
}

// Error 构造失败响应。
func Error(msg string) Body[any] {
	return Body[any]{Code: CodeError, Msg: msg}
}

// Warn 构造警告响应。
func Warn(msg string) Body[any] {
	return Body[any]{Code: CodeWarn, Msg: msg}
}

// Success 表示该响应是否为成功结果。
func (b Body[T]) Success() bool {
	return b.Code == CodeSuccess
}
