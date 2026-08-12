package mysql

// questionPO 映射 questions 表。
//
// 本表的 create_user / update_user 是 varchar(50)，因此嵌入 auditFieldsStrUser。
type questionPO struct {
	auditFieldsStrUser

	SetID         int64 `gorm:"column:set_id"`
	SetSequenceID int   `gorm:"column:set_sequence_id"`

	Content string `gorm:"column:content"`
	OptionA string `gorm:"column:option_a"`
	OptionB string `gorm:"column:option_b"`
	OptionC string `gorm:"column:option_c"`
	OptionD string `gorm:"column:option_d"`

	// CorrectAnswer 在表中是 char，取值为 A / B / C / D。
	CorrectAnswer string `gorm:"column:correct_answer"`
	Explanation   string `gorm:"column:explanation"`

	Remark string `gorm:"column:remark"`
}

func (questionPO) TableName() string { return "questions" }

// questionUserStatusPO 映射 question_user_status 表。
//
// 表结构把十个题库的进度平铺成十个列（question_set_1_id ~ question_set_10_id），
// 增删题库需要改表。这是一处结构性技术债，v2 应改为
// (user_id, set_id, progress) 的纵表；迁移前 PO 如实映射现状。
type questionUserStatusPO struct {
	auditFieldsStrUser

	UserID int64 `gorm:"column:user_id"`

	QuestionSet1  int `gorm:"column:question_set_1_id"`
	QuestionSet2  int `gorm:"column:question_set_2_id"`
	QuestionSet3  int `gorm:"column:question_set_3_id"`
	QuestionSet4  int `gorm:"column:question_set_4_id"`
	QuestionSet5  int `gorm:"column:question_set_5_id"`
	QuestionSet6  int `gorm:"column:question_set_6_id"`
	QuestionSet7  int `gorm:"column:question_set_7_id"`
	QuestionSet8  int `gorm:"column:question_set_8_id"`
	QuestionSet9  int `gorm:"column:question_set_9_id"`
	QuestionSet10 int `gorm:"column:question_set_10_id"`

	Remark string `gorm:"column:remark"`
}

func (questionUserStatusPO) TableName() string { return "question_user_status" }

// TODO: 实现与 explore 域模型的互转，十个题库列需要折叠为切片。
