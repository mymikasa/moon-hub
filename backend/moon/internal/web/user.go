package web

import (
	"fmt"
	"net/http"
	"time"

	"moon/internal/errs"
	"moon/internal/service"
	ijwt "moon/internal/web/jwt"
	"moon/pkg/ginx"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	bizLogin             = "login"
)

type UserHandler struct {
	ijwt.Handler
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            service.UserService
	// codeSvc        service.CodeService
}

func NewUserHandler(svc service.UserService,
	hdl ijwt.Handler,
) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
		Handler:        hdl,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.POST("/signup", ginx.WrapBody(h.SignUp))
	ug.POST("/login", ginx.WrapBody(h.LoginJWT))
	ug.POST("/logout", h.LogoutJWT)
	ug.GET("/refresh_token", h.RefreshToken)
	ug.GET("/profile", h.Profile)
	ug.PUT("/profile", ginx.WrapBody(h.UpdateProfile))
}

func (h *UserHandler) SignUp(ctx *gin.Context, req SignUpReq) (ginx.Result, error) {
	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		return ginx.Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统错误",
		}, err
	}
	if !isEmail {
		return ginx.Result{
			Code: errs.UserInvalidInput,
			Msg:  "非法邮箱格式",
		}, nil
	}

	if req.Password != req.ConfirmPassword {
		return ginx.Result{
			Code: errs.UserInvalidInput,
			Msg:  "两次输入的密码不相等",
		}, nil
	}

	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		return ginx.Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统错误",
		}, err
	}
	if !isPassword {
		return ginx.Result{
			Code: errs.UserInvalidInput,
			Msg:  "密码必须包含字母、数字、特殊字符",
		}, nil
	}

	err = h.svc.Signup(ctx.Request.Context(), req.Email, req.Password, req.Nickname)
	switch err {
	case nil:
		return ginx.Result{
			Msg: "注册成功",
		}, nil
	case service.ErrDuplicateEmail:
		return ginx.Result{
			Code: errs.UserDuplicateEmail,
			Msg:  "邮箱冲突",
		}, nil
	default:
		return ginx.Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统错误",
		}, err
	}
}

func (h *UserHandler) LoginJWT(ctx *gin.Context, req LoginJWTReq) (ginx.Result, error) {
	fmt.Printf("LoginJWT: 收到登录请求，邮箱: %s\n", req.Email)

	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		fmt.Printf("LoginJWT: 登录成功，用户ID: %d\n", u.Id)
		err = h.SetLoginToken(ctx, u.Id)
		if err != nil {
			fmt.Printf("LoginJWT: 设置token失败: %v\n", err)
			return ginx.Result{
				Code: 5,
				Msg:  "系统错误",
			}, err
		}
		fmt.Printf("LoginJWT: token设置成功\n")
		return ginx.Result{
			Msg: "OK",
		}, nil
	case service.ErrInvalidUserOrPassword:
		fmt.Printf("LoginJWT: 用户名或密码错误\n")
		return ginx.Result{Msg: "用户名或者密码错误"}, nil
	default:
		fmt.Printf("LoginJWT: 系统错误: %v\n", err)
		return ginx.Result{Msg: "系统错误"}, err
	}
}

//func (h *UserHandler) Logout(ctx *gin.Context) {
//	sess := sessions.Default(ctx)
//	sess.Options(sessions.Options{
//		MaxAge: -1,
//	})
//	sess.Save()
//}

func (h *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			// 十分钟
			MaxAge: 30,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	// 约定，前端在 Authorization 里面带上这个 refresh_token
	tokenStr := h.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RCJWTKey, nil
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if token == nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.CheckSession(ctx, rc.Ssid)
	if err != nil {
		// token 无效或者 redis 有问题
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = h.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "OK",
	})
}

func (h *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := h.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{Code: 5, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{Msg: "退出登录成功"})
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	u, err := h.svc.FindById(ctx.Request.Context(), uc.Uid)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{Code: 5, Msg: "系统错误"})
		return
	}

	resp := ProfileResp{
		Id:       u.Id,
		Email:    u.Email,
		Nickname: u.Nickname,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Phone:    u.Phone,
	}
	ctx.JSON(http.StatusOK, ginx.Result{Msg: "success", Data: resp})
}

func (h *UserHandler) UpdateProfile(ctx *gin.Context, req UpdateProfileReq) (ginx.Result, error) {
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	u, err := h.svc.FindById(ctx.Request.Context(), uc.Uid)
	if err != nil {
		return ginx.Result{Code: 5, Msg: "系统错误"}, err
	}

	u.Nickname = req.Nickname
	u.AboutMe = req.AboutMe
	u.Phone = req.Phone
	if req.Birthday != 0 {
		u.Birthday = time.UnixMilli(req.Birthday)
	}

	err = h.svc.Update(ctx.Request.Context(), u)
	if err != nil {
		return ginx.Result{Code: 5, Msg: "更新失败"}, err
	}

	return ginx.Result{Msg: "更新成功"}, nil
}
