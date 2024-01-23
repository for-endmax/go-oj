package form

// AddQuestionForm 添加问题表单
type AddQuestionForm struct {
	Seq     int32  `json:"seq" binding:"required"`
	Name    string `json:"name" binding:"required,min=3,max=10"`
	Content string `json:"content" binding:"required"`
}

// DelQuestionForm 删除问题表单
type DelQuestionForm struct {
	ID int32 `json:"id" binding:"required"`
}

// UpdateQuestionForm 修改问题表单
type UpdateQuestionForm struct {
	ID      int32  `json:"id" binding:"required"`
	Seq     int32  `json:"seq" binding:"required"`
	Name    string `json:"name" binding:"required,min=3,max=10"`
	Content string `json:"content" binding:"required"`
}

// AddTestForm 添加测试数据表单
type AddTestForm struct {
	QID    int32  `json:"q_id" binding:"required"`
	Input  string `json:"input" binding:"required"`
	Output string `json:"output" binding:"required"`
}

// UpdateTestForm 修改测试数据表单
type UpdateTestForm struct {
	ID     int32  `json:"id" binding:"required"`
	QID    int32  `json:"q_id" `
	Input  string `json:"input" `
	Output string `json:"output" `
}

// DelTestForm 删除测试数据表单
type DelTestForm struct {
	ID int32 `json:"id" binding:"required"`
}
