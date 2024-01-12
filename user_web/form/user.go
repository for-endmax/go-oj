package form

// LoginForm 登录表单
type LoginForm struct {
	NickName string `json:"nickname" binding:"required,min=3,max=10"`
	Password string `json:"password" binding:"required,min=6,max=15"`
}

// SignupForm 注册表单
type SignupForm struct {
	NickName string `json:"nickname" binding:"required,min=3,max=10"`
	Password string `json:"password" binding:"required,min=6,max=15"`
	Gender   string `json:"gender" binding:"required,oneof=male female"`
}
