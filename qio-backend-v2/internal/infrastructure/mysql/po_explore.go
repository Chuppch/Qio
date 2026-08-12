package mysql

import (
	"strconv"

	"github.com/Chuppch/Qio/qio-backend-v2/internal/domain/explore"
)

// questionPO 映射 questions 表。
//
// 本表的 create_user / update_user 是 varchar(50)，因此嵌入 auditFieldsStrUser。
type questionPO struct {
	auditFieldsStrUser
	SetID         int64  `gorm:"column:set_id"`
	SetSequenceID int    `gorm:"column:set_sequence_id"`
	Content       string `gorm:"column:content"`
	OptionA       string `gorm:"column:option_a"`
	OptionB       string `gorm:"column:option_b"`
	OptionC       string `gorm:"column:option_c"`
	OptionD       string `gorm:"column:option_d"`
	// CorrectAnswer 在表中是 char，取值为 A / B / C / D。
	CorrectAnswer string `gorm:"column:correct_answer"`
	Explanation   string `gorm:"column:explanation"`
	Remark        string `gorm:"column:remark"`
}

func (questionPO) TableName() string { return "questions" }

// toDomain 把四个选项列折叠为映射。
//
// 出题场景下 correct_answer 与 explanation 不在查询列中，此处保持零值，
// 由 explore.Question.Answered 区分。
func (p questionPO) toDomain() *explore.Question {
	return &explore.Question{
		ID:       p.ID,
		SetID:    p.SetID,
		Sequence: p.SetSequenceID,
		Content:  p.Content,
		Choices: map[explore.Option]string{
			explore.OptionA: p.OptionA,
			explore.OptionB: p.OptionB,
			explore.OptionC: p.OptionC,
			explore.OptionD: p.OptionD,
		},
		CorrectAnswer: explore.Option(p.CorrectAnswer),
		Explanation:   p.Explanation,
	}
}

// questionUserStatusPO 映射 question_user_status 表。
//
// 表结构把十个题库的进度平铺成十个列（question_set_1_id ~ question_set_10_id），
// 增删题库需要改表。这是一处结构性技术债，v2 应改为
// (user_id, set_id, progress) 的纵表；迁移前 PO 如实映射现状。
type questionUserStatusPO struct {
	auditFieldsStrUser
	UserID        int64  `gorm:"column:user_id"`
	QuestionSet1  int    `gorm:"column:question_set_1_id"`
	QuestionSet2  int    `gorm:"column:question_set_2_id"`
	QuestionSet3  int    `gorm:"column:question_set_3_id"`
	QuestionSet4  int    `gorm:"column:question_set_4_id"`
	QuestionSet5  int    `gorm:"column:question_set_5_id"`
	QuestionSet6  int    `gorm:"column:question_set_6_id"`
	QuestionSet7  int    `gorm:"column:question_set_7_id"`
	QuestionSet8  int    `gorm:"column:question_set_8_id"`
	QuestionSet9  int    `gorm:"column:question_set_9_id"`
	QuestionSet10 int    `gorm:"column:question_set_10_id"`
	Remark        string `gorm:"column:remark"`
}

func (questionUserStatusPO) TableName() string { return "question_user_status" }

// questionSetColumns 按套号顺序列出十个进度列名。
//
// v1 靠反射拼 "getQuestionSet" + i 来读这些列，拼错只在运行期暴露。这里改为显式
// 列表，与下面的 setters / getters 一一对应，长度由 SetCount 约束。
var questionSetColumns = [explore.SetCount]string{
	"question_set_1_id",
	"question_set_2_id",
	"question_set_3_id",
	"question_set_4_id",
	"question_set_5_id",
	"question_set_6_id",
	"question_set_7_id",
	"question_set_8_id",
	"question_set_9_id",
	"question_set_10_id",
}

// setFields 按套号顺序返回十个进度字段的指针，供折叠与展开复用。
func (p *questionUserStatusPO) setFields() [explore.SetCount]*int {
	return [explore.SetCount]*int{
		&p.QuestionSet1,
		&p.QuestionSet2,
		&p.QuestionSet3,
		&p.QuestionSet4,
		&p.QuestionSet5,
		&p.QuestionSet6,
		&p.QuestionSet7,
		&p.QuestionSet8,
		&p.QuestionSet9,
		&p.QuestionSet10,
	}
}

// toDomain 把十个平铺的进度列折叠为映射。
func (p questionUserStatusPO) toDomain() *explore.Progress {
	sets := make(map[int64]int, explore.SetCount)
	fields := p.setFields()
	for i, f := range fields {
		sets[int64(i+1)] = *f
	}

	return &explore.Progress{
		ID:        p.ID,
		UserID:    p.UserID,
		Sets:      sets,
		CreatedAt: p.CreateTime,
		UpdatedAt: p.UpdateTime,
	}
}

// questionUserStatusPOFrom 把领域进度展开回十个平铺列。
//
// create_user / update_user 在表中是 NOT NULL varchar，v1 由 MyBatis-Plus 的
// 自动填充写入当前用户 ID，这里显式填。
func questionUserStatusPOFrom(p *explore.Progress) questionUserStatusPO {
	uid := strconv.FormatInt(p.UserID, 10)

	po := questionUserStatusPO{
		auditFieldsStrUser: auditFieldsStrUser{
			ID:         p.ID,
			CreateUser: uid,
			UpdateUser: uid,
		},
		UserID: p.UserID,
	}

	fields := po.setFields()
	for i, f := range fields {
		*f = p.Sets[int64(i+1)]
	}
	return po
}

// setColumnValues 返回进度列到取值的映射，供更新时按列写入。
func (p *questionUserStatusPO) setColumnValues() map[string]any {
	out := make(map[string]any, explore.SetCount)
	fields := p.setFields()
	for i, col := range questionSetColumns {
		out[col] = *fields[i]
	}
	return out
}
