package response

// RecordInfoResponse 记录信息
type RecordInfoResponse struct {
	ID         int32  `json:"id"`
	UID        int32  `json:"u_id"`
	QID        int32  `json:"q_id"`
	Lang       string `json:"lang"`
	Status     int32  `json:"status"`
	ErrCode    int32  `json:"err_code"`
	ErrMsg     string `json:"err_msg"`
	TimeLimit  int32  `json:"time_limit"`
	MemLimit   int32  `json:"mem_limit"`
	SubmitCode string `json:"submit_code"`
	MemUsage   int32  `json:"mem_usage"`
	TimeUsage  int32  `json:"time_usage"`
}

// RecordInfoListResponse 用户的所有记录
type RecordInfoListResponse struct {
	Total int32                `json:"total"`
	Data  []RecordInfoResponse `json:"data"`
}
