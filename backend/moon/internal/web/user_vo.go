package web

type SignUpReq struct {
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Nickname        string `json:"nickname" binding:"required"`
}

type LoginJWTReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
