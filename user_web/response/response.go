package response

type UserResponse struct {
	ID       int32  `json:"id"`
	Nickname string `json:"nickname"`
	Gender   string `json:"gender"`
	Role     int32  `json:"role"`
}
