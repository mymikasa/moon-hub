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

type ProfileResp struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Birthday int64  `json:"birthday"`
	AboutMe  string `json:"about_me"`
	Phone    string `json:"phone"`
}

type UpdateProfileReq struct {
	Nickname string `json:"nickname"`
	Birthday int64  `json:"birthday"`
	AboutMe  string `json:"about_me"`
	Phone    string `json:"phone"`
}
