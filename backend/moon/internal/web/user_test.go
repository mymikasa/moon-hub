package web

import (
	"bytes"
	"context"
	"encoding/json"
	"moon/internal/domain"
	"moon/internal/service"
	"moon/pkg/ginx"
	"moon/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	ginx.L = logger.NewNopLogger()
}

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) Signup(ctx context.Context, email, password, nickname string) error {
	args := m.Called(ctx, email, password, nickname)
	return args.Error(0)
}

func (m *mockUserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	args := m.Called(ctx, email, password)
	if args.Get(0) == nil {
		return domain.User{}, args.Error(1)
	}
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *mockUserService) FindById(ctx context.Context, id int64) (domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return domain.User{}, args.Error(1)
	}
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *mockUserService) Update(ctx context.Context, u domain.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

type mockJWTHandler struct {
	mock.Mock
}

func (m *mockJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *mockJWTHandler) ExtractToken(ctx *gin.Context) string {
	args := m.Called(ctx)
	return args.String(0)
}

func (m *mockJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	args := m.Called(ctx, ssid)
	return args.Error(0)
}

func (m *mockJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	args := m.Called(ctx, uid, ssid)
	return args.Error(0)
}

func (m *mockJWTHandler) ClearToken(ctx *gin.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func setupTestRouter(handler *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler.RegisterRoutes(router)
	return router
}

func TestUserHandler_SignUp(t *testing.T) {
	tests := []struct {
		name      string
		reqBody   SignUpReq
		mockSetup func(*mockUserService)
		wantCode  int
		wantMsg   string
	}{
		{
			name: "成功注册",
			reqBody: SignUpReq{
				Email:           "test@example.com",
				Password:        "Password123!",
				ConfirmPassword: "Password123!",
				Nickname:        "testuser",
			},
			mockSetup: func(m *mockUserService) {
				m.On("Signup", mock.Anything, "test@example.com", "Password123!", "testuser").Return(nil)
			},
			wantCode: http.StatusOK,
			wantMsg:  "注册成功",
		},
		{
			name: "邮箱格式错误",
			reqBody: SignUpReq{
				Email:           "invalid-email",
				Password:        "Password123!",
				ConfirmPassword: "Password123!",
				Nickname:        "testuser",
			},
			mockSetup: func(m *mockUserService) {},
			wantCode:  http.StatusOK,
			wantMsg:   "非法邮箱格式",
		},
		{
			name: "密码不匹配",
			reqBody: SignUpReq{
				Email:           "test@example.com",
				Password:        "Password123!",
				ConfirmPassword: "DifferentPassword123!",
				Nickname:        "testuser",
			},
			mockSetup: func(m *mockUserService) {},
			wantCode:  http.StatusOK,
			wantMsg:   "两次输入的密码不相等",
		},
		{
			name: "密码格式错误",
			reqBody: SignUpReq{
				Email:           "test@example.com",
				Password:        "simple",
				ConfirmPassword: "simple",
				Nickname:        "testuser",
			},
			mockSetup: func(m *mockUserService) {},
			wantCode:  http.StatusOK,
			wantMsg:   "密码必须包含字母、数字、特殊字符",
		},
		{
			name: "邮箱已存在",
			reqBody: SignUpReq{
				Email:           "existing@example.com",
				Password:        "Password123!",
				ConfirmPassword: "Password123!",
				Nickname:        "testuser",
			},
			mockSetup: func(m *mockUserService) {
				m.On("Signup", mock.Anything, "existing@example.com", "Password123!", "testuser").Return(service.ErrDuplicateEmail)
			},
			wantCode: http.StatusOK,
			wantMsg:  "邮箱冲突",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockUserService)
			mockHdl := new(mockJWTHandler)
			tt.mockSetup(mockSvc)

			handler := NewUserHandler(mockSvc, mockHdl)
			router := setupTestRouter(handler)

			body, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest("POST", "/users/signup", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var resp ginx.Result
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantMsg, resp.Msg)
			mockSvc.AssertExpectations(t)
		})
	}
}

func TestUserHandler_LoginJWT(t *testing.T) {
	tests := []struct {
		name      string
		reqBody   LoginJWTReq
		mockSetup func(*mockUserService, *mockJWTHandler)
		wantCode  int
		wantMsg   string
	}{
		{
			name: "成功登录",
			reqBody: LoginJWTReq{
				Email:    "test@example.com",
				Password: "Password123!",
			},
			mockSetup: func(m *mockUserService, h *mockJWTHandler) {
				m.On("Login", mock.Anything, "test@example.com", "Password123!").Return(domain.User{Id: 1}, nil)
				h.On("SetLoginToken", mock.Anything, int64(1)).Return(nil)
			},
			wantCode: http.StatusOK,
			wantMsg:  "OK",
		},
		{
			name: "用户名或密码错误",
			reqBody: LoginJWTReq{
				Email:    "test@example.com",
				Password: "WrongPassword!",
			},
			mockSetup: func(m *mockUserService, h *mockJWTHandler) {
				m.On("Login", mock.Anything, "test@example.com", "WrongPassword!").Return(domain.User{}, service.ErrInvalidUserOrPassword)
			},
			wantCode: http.StatusOK,
			wantMsg:  "用户名或者密码错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockUserService)
			mockHdl := new(mockJWTHandler)
			tt.mockSetup(mockSvc, mockHdl)

			handler := NewUserHandler(mockSvc, mockHdl)
			router := setupTestRouter(handler)

			body, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var resp ginx.Result
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantMsg, resp.Msg)
			mockSvc.AssertExpectations(t)
			mockHdl.AssertExpectations(t)
		})
	}
}
