package explore

import (
	"strings"
	"time"
)

// SetCount 是题库的套数。
//
// v1 把进度平铺成 question_set_1_id ~ question_set_10_id 十个列，套数因此被
// 固定在表结构里，增删一套题就要改表。这个常量如实反映该约束，
// 待表改为 (user_id, set_id, state) 纵表后即可移除。
const SetCount = 10

// Option 是选项标识，取值为 A / B / C / D。
//
// v1 用 char 存 correct_answer、用 String 收前端提交的答案，两侧都是裸字符串。
// 这里收敛为具名类型，让「选项」这一概念在签名上可见。
type Option string

const (
	OptionA Option = "A"
	OptionB Option = "B"
	OptionC Option = "C"
	OptionD Option = "D"
)

// AllOptions 按题面展示顺序列出全部选项。
var AllOptions = []Option{OptionA, OptionB, OptionC, OptionD}

// Question 是一道侨批文化题。
type Question struct {
	ID    int64
	SetID int64

	// Sequence 是题目在套内的序号，对应 set_sequence_id。
	Sequence int

	Content string

	// Choices 是各选项的文案。
	//
	// v1 用 option_a ~ option_d 四个列存储，四个列在业务上是同一件东西的四份拷贝，
	// 任何遍历都要写死四次。这里折叠为映射，列与键的互转由仓储层负责。
	Choices map[Option]string

	// CorrectAnswer 是正确选项。
	//
	// 出题时不下发该字段（见 Repository.ListQuestionsWithoutAnswer），此时为空串。
	CorrectAnswer Option

	Explanation string
}

// Choice 返回指定选项的文案，选项不存在时返回零值与 false。
func (q *Question) Choice(o Option) (string, bool) {
	text, ok := q.Choices[o]
	return text, ok
}

// Answered 表示该题携带了正确答案，即不是出题用的脱敏副本。
func (q *Question) Answered() bool { return q.CorrectAnswer != "" }

// Correct 判定作答是否正确，忽略大小写。
//
// v1 用 equalsIgnoreCase 比较，保持一致。未携带正确答案的脱敏副本恒为 false。
func (q *Question) Correct(answer Option) bool {
	if !q.Answered() {
		return false
	}
	return strings.EqualFold(string(answer), string(q.CorrectAnswer))
}

// Progress 是用户的答题进度。
//
// v1 中对应 QuestionUserStatus。一个用户至多一条记录，由 user_id 上的唯一索引保证。
type Progress struct {
	ID     int64
	UserID int64

	// Sets 记录每套题的完成状态，键为套号（自 1 起），值沿用 v1 的 0 未完成 / 1 已完成。
	//
	// 列名叫 question_set_N_id 但存的并不是 ID，v1 只往里写 1。等价迁移保留取值，
	// 语义上的名不副实记录在 docs/TODO-migration.md。
	Sets map[int64]int

	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewProgress 构造一条初始进度，所有套题均未完成。
//
// 对应 v1 userLoginPage 中「首次答题插入一条只带 user_id 的记录」，
// 其余列依赖数据库默认值 0；这里显式填零，不依赖建表语句。
func NewProgress(userID int64) *Progress {
	sets := make(map[int64]int, SetCount)
	for i := int64(1); i <= SetCount; i++ {
		sets[i] = 0
	}
	return &Progress{UserID: userID, Sets: sets}
}

// ValidSet 表示套号在有效范围内。
func ValidSet(setID int64) bool { return setID >= 1 && setID <= SetCount }

// Completed 表示指定套题已完成。
func (p *Progress) Completed(setID int64) bool { return p.Sets[setID] == 1 }

// Complete 把指定套题标记为已完成。
//
// 套号越界时返回 ErrInvalidSet，对应 v1 switch 的 default 分支抛
// IllegalArgumentException。
func (p *Progress) Complete(setID int64) error {
	if !ValidSet(setID) {
		return ErrInvalidSet
	}
	if p.Sets == nil {
		p.Sets = make(map[int64]int, SetCount)
	}
	p.Sets[setID] = 1
	return nil
}

// States 按套号顺序返回长度为 SetCount 的完成状态切片。
//
// 对应 v1 userLoginPage 的返回值——它靠反射逐个调用 getQuestionSetN 拼出这个列表，
// 拼错方法名只会在运行期报错。这里由映射直接投影，越界不再是运行期风险。
func (p *Progress) States() []int {
	out := make([]int, SetCount)
	for i := 0; i < SetCount; i++ {
		out[i] = p.Sets[int64(i+1)]
	}
	return out
}

// NextSet 返回指定套题之后的下一套套号。
//
// 对应 v1 handleCompleteAnswers 中的 setId + 1。v1 不校验上界，答完第 10 套会去查
// set_id = 11，查不到就返回空列表；这里保留该行为，仅由调用方决定如何呈现。
func NextSet(setID int64) int64 { return setID + 1 }

// DrawLimit 是用户当日剩余的抽奖次数。
//
// v1 中对应 GameServiceImpl 的「发财树」玩法，次数只落 Redis，没有数据表。
type DrawLimit struct {
	UserID    int64
	Remaining int
}

// DefaultDrawLimit 是每日初始抽奖次数，取自 v1 GameServiceImpl 的硬编码值。
const DefaultDrawLimit = 10

// DrawReward 是单次抽奖发放的猪仔钱，取自 v1 GameServiceImpl 的硬编码值。
const DrawReward = 10

// Exhausted 表示当日次数已用尽。
func (d DrawLimit) Exhausted() bool { return d.Remaining <= 0 }
