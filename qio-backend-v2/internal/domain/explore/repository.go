package explore

import "context"

// Repository 是答题玩法的数据访问接口，实现由 internal/infrastructure/mysql 提供。
//
// 方法集对应 v1 QuestionServiceImpl 中对 questionsMapper 与
// questionUserStatusMapper 的全部调用。已注释掉的 genQuestion 一段不计入。
type Repository interface {
	// ListQuestions 查询整套题目，含正确答案与解析。
	//
	// 对应 v1 handleCompleteAnswers 与 allAnswerToFront 中按 set_id 的全列查询。
	// 套号不存在时返回空切片而非错误——v1 答完最后一套会去查越界的下一套，
	// 依赖的正是空结果。
	ListQuestions(ctx context.Context, setID int64) ([]*Question, error)

	// ListQuestionsWithoutAnswer 查询整套题目，但不带正确答案与解析。
	//
	// 对应 v1 startAnswer 中的列投影，用于出题时避免答案随响应下发。
	// 返回的 Question.CorrectAnswer 为空串。
	ListQuestionsWithoutAnswer(ctx context.Context, setID int64) ([]*Question, error)

	// FindQuestion 按主键查询单题，不存在时返回 ErrQuestionNotFound。
	FindQuestion(ctx context.Context, id int64) (*Question, error)

	// FindProgress 查询用户的答题进度，无记录时返回 ErrProgressNotFound。
	FindProgress(ctx context.Context, userID int64) (*Progress, error)

	// CreateProgress 创建初始答题进度，成功后回填 ID。
	CreateProgress(ctx context.Context, p *Progress) error

	// SaveProgress 保存答题进度中各套题的完成状态。
	SaveProgress(ctx context.Context, p *Progress) error
}

// DrawRepository 是抽奖次数的数据访问接口，实现由 internal/infrastructure/redis 提供。
//
// v1 的抽奖次数只存 Redis，没有对应数据表，因此与 Repository 分开定义：
// 两者的存储介质不同，生命周期也不同（次数按天过期，进度长期保留）。
type DrawRepository interface {
	// FindLimit 查询用户当日剩余抽奖次数。
	//
	// 键不存在时按 DefaultDrawLimit 初始化并写回，与 v1 getFflLimit 一致。
	FindLimit(ctx context.Context, userID int64) (DrawLimit, error)

	// Consume 扣减一次抽奖次数并返回扣减后的剩余值。
	//
	// 次数已用尽时返回 ErrDrawLimitExhausted。
	Consume(ctx context.Context, userID int64) (DrawLimit, error)
}
