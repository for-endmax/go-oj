package response

// QuestionResponse 题目信息
type QuestionResponse struct {
	ID      int32  `json:"id"`
	Seq     int32  `json:"seq"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

// QuestionBrief 题目简要信息
type QuestionBrief struct {
	ID   int32  `json:"id"`
	Seq  int32  `json:"seq"`
	Name string `json:"name"`
}

// TestInfo 测试数据
type TestInfo struct {
	QID       int32  `json:"q_id"`
	TimeLimit int32  `json:"time_limit"`
	MemLimit  int32  `json:"mem_limit"`
	Input     string `json:"input"`
	Output    string `json:"output"`
}
