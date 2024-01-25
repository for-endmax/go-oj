package message

// MsgSend 通过mq发送的记录信息
type MsgSend struct {
	ID         int32  `json:"id"`
	Lang       string `json:"lang,omitempty"`
	SubmitCode string `json:"submit_code,omitempty"`
	//新增
	TimeLimit int32 `json:"time_limit"`
	MemLimit  int32 `json:"mem_limit"`
}

// MsgReply 接收的mq回调消息
type MsgReply struct {
	ID        int32  `json:"id"`
	Status    int32  `json:"status,omitempty"`
	ErrCode   int32  `json:"err_code,omitempty"`
	ErrMsg    string `json:"err_msg,omitempty"`
	TimeUsage int32  `json:"time_usage"`
	MemUsage  int32  `json:"mem_usage"`
}
