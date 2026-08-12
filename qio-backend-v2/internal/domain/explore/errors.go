package explore

import "errors"

// 文化探索域的错误值。
//
// 对应 v1 QuestionException 与 GameServiceImpl 中抛出的 UserException 消息码。
// v1 有几处直接抛 RuntimeException 与 IllegalArgumentException，未纳入统一的
// 异常体系，这里一并收敛为域内错误值。
var (
	// ErrQuestionNotFound 题目不存在。
	// 对应 v1 submitAnswers 中的 RuntimeException("题目未找到")
	ErrQuestionNotFound = errors.New("question not found")

	// ErrProgressNotFound 用户尚无答题进度记录。
	//
	// 这不是异常路径：v1 据此判定用户首次答题，随即插入一条初始记录。
	ErrProgressNotFound = errors.New("question progress not found")

	// ErrInvalidSet 套号越界。
	// 对应 v1 updateUserQuestionStatus 中 switch default 抛出的 IllegalArgumentException
	ErrInvalidSet = errors.New("invalid question set")

	// ErrGetSetFailed 读取答题进度失败。
	// 对应消息码 question.getSetId.failed
	ErrGetSetFailed = errors.New("get question set failed")

	// ErrDrawLimitExhausted 当日抽奖次数已用尽。
	// 对应消息码 ffl.limit.error
	ErrDrawLimitExhausted = errors.New("draw limit exhausted")
)
