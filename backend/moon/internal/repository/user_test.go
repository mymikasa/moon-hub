package repository

import (
	"context"
	"database/sql"
	"moon/internal/domain"
	"moon/internal/repository/dao"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserDAO struct {
	mock.Mock
	users map[string]dao.User
	err   error
}

func (m *mockUserDAO) Insert(ctx context.Context, u dao.User) error {
	m.Called(ctx, u)
	if m.err != nil {
		return m.err
	}
	if m.users == nil {
		m.users = make(map[string]dao.User)
	}
	m.users[u.Email.String] = u
	return nil
}

func (m *mockUserDAO) FindByEmail(ctx context.Context, email string) (dao.User, error) {
	_ = m.Called(ctx, email)
	if m.err != nil {
		return dao.User{}, m.err
	}
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return dao.User{}, dao.ErrRecordNotFound
}

func (m *mockUserDAO) FindById(ctx context.Context, id int64) (dao.User, error) {
	_ = m.Called(ctx, id)
	if m.err != nil {
		return dao.User{}, m.err
	}
	for _, u := range m.users {
		if u.Id == id {
			return u, nil
		}
	}
	return dao.User{}, dao.ErrRecordNotFound
}

func (m *mockUserDAO) Update(ctx context.Context, u dao.User) error {
	_ = m.Called(ctx, u)
	return m.err
}

func TestGORMUserRepository_Create(t *testing.T) {
	tests := []struct {
		name      string
		user      domain.User
		mockSetup func(*mockUserDAO)
		wantErr   error
	}{
		{
			name: "创建成功",
			user: domain.User{
				Id:       1,
				Email:    "test@example.com",
				Password: "hashedPassword",
				Nickname: "testuser",
				Ctime:    time.Now(),
			},
			mockSetup: func(m *mockUserDAO) {
				m.On("Insert", mock.Anything, mock.AnythingOfType("dao.User")).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "邮箱冲突",
			user: domain.User{
				Id:       1,
				Email:    "existing@example.com",
				Password: "hashedPassword",
				Nickname: "testuser",
				Ctime:    time.Now(),
			},
			mockSetup: func(m *mockUserDAO) {
				m.err = dao.ErrDuplicateEmail
				m.On("Insert", mock.Anything, mock.AnythingOfType("dao.User")).Return(nil)
			},
			wantErr: ErrDuplicateUser,
		},
		{
			name: "数据库错误",
			user: domain.User{
				Id:       1,
				Email:    "test@example.com",
				Password: "hashedPassword",
				Nickname: "testuser",
				Ctime:    time.Now(),
			},
			mockSetup: func(m *mockUserDAO) {
				m.err = assert.AnError
				m.On("Insert", mock.Anything, mock.AnythingOfType("dao.User")).Return(nil)
			},
			wantErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDAO := new(mockUserDAO)
			tt.mockSetup(mockDAO)

			repo := NewGORMUserRepository(mockDAO)
			err := repo.Create(context.Background(), tt.user)

			assert.Equal(t, tt.wantErr, err)
			mockDAO.AssertExpectations(t)
		})
	}
}

func TestGORMUserRepository_FindByEmail(t *testing.T) {
	birthday := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	ctime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	daoUser := dao.User{
		Id:       1,
		Email:    sql.NullString{String: "test@example.com", Valid: true},
		Password: "hashedPassword",
		Nickname: "testuser",
		Birthday: birthday.UnixMilli(),
		AboutMe:  "about me",
		Phone:    sql.NullString{String: "1234567890", Valid: true},
		Ctime:    ctime.UnixMilli(),
	}

	tests := []struct {
		name      string
		email     string
		mockSetup func(*mockUserDAO)
		wantUser  domain.User
		wantErr   error
	}{
		{
			name:  "查找成功",
			email: "test@example.com",
			mockSetup: func(m *mockUserDAO) {
				m.users = make(map[string]dao.User)
				m.users["test@example.com"] = daoUser
				m.On("FindByEmail", mock.Anything, "test@example.com").Return(nil)
			},
			wantUser: domain.User{
				Id:       1,
				Email:    "test@example.com",
				Password: "hashedPassword",
				Nickname: "testuser",
				Birthday: time.UnixMilli(birthday.UnixMilli()),
				AboutMe:  "about me",
				Phone:    "1234567890",
				Ctime:    time.UnixMilli(daoUser.Ctime),
			},
			wantErr: nil,
		},
		{
			name:  "用户不存在",
			email: "nonexistent@example.com",
			mockSetup: func(m *mockUserDAO) {
				m.users = make(map[string]dao.User)
				m.err = dao.ErrRecordNotFound
				m.On("FindByEmail", mock.Anything, "nonexistent@example.com").Return(nil)
			},
			wantErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDAO := new(mockUserDAO)
			tt.mockSetup(mockDAO)

			repo := NewGORMUserRepository(mockDAO)
			user, err := repo.FindByEmail(context.Background(), tt.email)

			assert.Equal(t, tt.wantErr, err)
			if tt.wantErr == nil {
				assert.Equal(t, tt.wantUser, user)
			}
			mockDAO.AssertExpectations(t)
		})
	}
}

func TestDomainToDaoUser(t *testing.T) {
	birthday := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	ctime := time.Now()

	domainUser := domain.User{
		Id:       1,
		Email:    "test@example.com",
		Password: "password",
		Nickname: "testuser",
		Birthday: birthday,
		AboutMe:  "about me",
		Phone:    "1234567890",
		Ctime:    ctime,
	}

	daoUser := domainToDaoUser(domainUser)

	assert.Equal(t, int64(1), daoUser.Id)
	assert.True(t, daoUser.Email.Valid)
	assert.Equal(t, "test@example.com", daoUser.Email.String)
	assert.Equal(t, "password", daoUser.Password)
	assert.Equal(t, "testuser", daoUser.Nickname)
	assert.Equal(t, birthday.UnixMilli(), daoUser.Birthday)
	assert.Equal(t, "about me", daoUser.AboutMe)
	assert.True(t, daoUser.Phone.Valid)
	assert.Equal(t, "1234567890", daoUser.Phone.String)
	assert.Equal(t, ctime.UnixMilli(), daoUser.Ctime)
}

func TestDaoToDomainUser(t *testing.T) {
	birthday := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	ctime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	daoUser := dao.User{
		Id:       1,
		Email:    sql.NullString{String: "test@example.com", Valid: true},
		Password: "password",
		Nickname: "testuser",
		Birthday: birthday.UnixMilli(),
		AboutMe:  "about me",
		Phone:    sql.NullString{String: "1234567890", Valid: true},
		Ctime:    ctime.UnixMilli(),
	}

	domainUser := daoToDomainUser(daoUser)

	assert.Equal(t, int64(1), domainUser.Id)
	assert.Equal(t, "test@example.com", domainUser.Email)
	assert.Equal(t, "password", domainUser.Password)
	assert.Equal(t, "testuser", domainUser.Nickname)
	assert.Equal(t, birthday.UnixMilli(), domainUser.Birthday.UnixMilli())
	assert.Equal(t, "about me", domainUser.AboutMe)
	assert.Equal(t, "1234567890", domainUser.Phone)
	assert.Equal(t, ctime.UnixMilli(), domainUser.Ctime.UnixMilli())
}
