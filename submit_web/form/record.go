package form

// SubmitForm 用户提交代码的表单
type SubmitForm struct {
	UID        int32  `json:"u_id" binding:"required"`
	QID        int32  `json:"q_id" binding:"required"`
	Lang       string `json:"lang" binding:"required,oneof=go c++ c python"`
	SubmitCode string `json:"submit_code" binding:"required,max=10000"`
}
