package service

import (
	"context"
	"moon/internal/domain"
	"moon/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	users     map[string]domain.User
	emailUsed bool
	createErr error
}

func (m *mockUserRepository) Create(ctx context.Context, u domain.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	if m.emailUsed {
		return repository.ErrDuplicateUser
	}
	m.users[u.Email] = u
	return nil
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return domain.User{}, repository.ErrUserNotFound
}

func (m *mockUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	for _, u := range m.users {
		if u.Id == id {
			return u, nil
		}
	}
	return domain.User{}, repository.ErrUserNotFound
}

func (m *mockUserRepository) Update(ctx context.Context, u domain.User) error {
	m.users[u.Email] = u
	return nil
}

func TestUserService_Signup(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		password  string
		nickname  string
		mockSetup func(*mockUserRepository)
		wantErr   error
	}{
		{
			name:     "成功注册",
			email:    "test@example.com",
			password: "Password123!",
			nickname: "testuser",
			mockSetup: func(m *mockUserRepository) {
				m.users = make(map[string]domain.User)
			},
			wantErr: nil,
		},
		{
			name:     "邮箱已存在",
			email:    "existing@example.com",
			password: "Password123!",
			nickname: "testuser",
			mockSetup: func(m *mockUserRepository) {
				m.users = make(map[string]domain.User)
				m.emailUsed = true
			},
			wantErr: ErrDuplicateEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepository{}
			tt.mockSetup(mockRepo)

			svc := &userService{repo: mockRepo}
			err := svc.Signup(context.Background(), tt.email, tt.password, tt.nickname)

			assert.Equal(t, tt.wantErr, err)

			if tt.wantErr == nil {
				u, found := mockRepo.users[tt.email]
				assert.True(t, found, "用户应该被创建")
				assert.Equal(t, tt.email, u.Email)
				assert.Equal(t, tt.nickname, u.Nickname)
				assert.NotEmpty(t, u.Password, "密码应该被加密")
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	password := "Password123!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tests := []struct {
		name      string
		email     string
		password  string
		mockSetup func(*mockUserRepository)
		wantErr   error
	}{
		{
			name:     "成功登录",
			email:    "test@example.com",
			password: password,
			mockSetup: func(m *mockUserRepository) {
				m.users = make(map[string]domain.User)
				m.users["test@example.com"] = domain.User{
					Id:       1,
					Email:    "test@example.com",
					Password: string(hashedPassword),
					Nickname: "testuser",
				}
			},
			wantErr: nil,
		},
		{
			name:     "用户不存在",
			email:    "nonexistent@example.com",
			password: password,
			mockSetup: func(m *mockUserRepository) {
				m.users = make(map[string]domain.User)
			},
			wantErr: ErrInvalidUserOrPassword,
		},
		{
			name:     "密码错误",
			email:    "test@example.com",
			password: "WrongPassword123!",
			mockSetup: func(m *mockUserRepository) {
				m.users = make(map[string]domain.User)
				m.users["test@example.com"] = domain.User{
					Id:       1,
					Email:    "test@example.com",
					Password: string(hashedPassword),
					Nickname: "testuser",
				}
			},
			wantErr: ErrInvalidUserOrPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepository{}
			tt.mockSetup(mockRepo)

			svc := &userService{repo: mockRepo}
			user, err := svc.Login(context.Background(), tt.email, tt.password)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.email, user.Email)
			}
		})
	}
}
