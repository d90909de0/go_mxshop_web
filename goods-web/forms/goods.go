package forms

type PassWordLoginForm struct {
	Mobile   string `json:"mobile" form:"mobile" binding:"required,mobile"`
	PassWord string `json:"passWord" form:"passWord" binding:"required"`
}

type CreateUserForm struct {
	Mobile   string `json:"mobile" form:"mobile" binding:"required,mobile"`
	PassWord string `json:"passWord" form:"passWord" binding:"required"`
	NickName string `json:"nickName" form:"nickName" binding:"required"`
}
