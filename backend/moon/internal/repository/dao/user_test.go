package dao

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	testCases := []struct {
		name string
		mock func(t *testing.T) *sql.DB
		ctx  context.Context
		user User

		wantErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(123, 1)
				mock.ExpectExec("INSERT INTO .*").WithArgs(
					sqlmock.AnyArg(), sqlmock.AnyArg(), "Tom", sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				).WillReturnResult(mockRes)
				return db
			},
			ctx: context.Background(),
			user: User{
				Nickname: "Tom",
			},
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("INSERT INTO .*").
					WillReturnError(&mysqlDriver.MySQLError{Number: 1062})
				return db
			},
			ctx: context.Background(),
			user: User{
				Nickname: "Tom",
			},
			wantErr: ErrDuplicateEmail,
		},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("INSERT INTO .*").
					WillReturnError(errors.New("数据库错误"))
				return db
			},
			ctx: context.Background(),
			user: User{
				Nickname: "Tom",
			},
			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			assert.NoError(t, err)
			dao := NewUserDAO(db)
			err = dao.Insert(tc.ctx, tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestGORMUserDAO_FindByEmail(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(t *testing.T) *sql.DB
		ctx      context.Context
		email    string
		wantErr  error
		wantUser User
	}{
		{
			name: "查找成功",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"id", "email", "password", "nickname", "birthday", "about_me", "phone", "ctime", "utime"}).
					AddRow(int64(1), "test@example.com", "hashedPassword", "testuser", int64(0), "", sql.NullString{String: "", Valid: false}, int64(1234567890000), int64(1234567890000))
				mock.ExpectQuery("SELECT \\* FROM .*").
					WithArgs("test@example.com", 1).
					WillReturnRows(rows)
				return db
			},
			ctx:     context.Background(),
			email:   "test@example.com",
			wantErr: nil,
			wantUser: User{
				Id:       1,
				Email:    sql.NullString{String: "test@example.com", Valid: true},
				Password: "hashedPassword",
				Nickname: "testuser",
				Birthday: 0,
				AboutMe:  "",
				Phone:    sql.NullString{String: "", Valid: false},
				Ctime:    1234567890000,
				Utime:    1234567890000,
			},
		},
		{
			name: "用户不存在",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery("SELECT \\* FROM .*").
					WithArgs("nonexistent@example.com", 1).
					WillReturnError(ErrRecordNotFound)
				return db
			},
			ctx:     context.Background(),
			email:   "nonexistent@example.com",
			wantErr: ErrRecordNotFound,
		},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectQuery("SELECT \\* FROM .*").
					WithArgs("test@example.com", 1).
					WillReturnError(errors.New("数据库错误"))
				return db
			},
			ctx:     context.Background(),
			email:   "test@example.com",
			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			assert.NoError(t, err)
			dao := NewUserDAO(db)
			user, err := dao.FindByEmail(tc.ctx, tc.email)
			assert.Equal(t, tc.wantErr, err)
			if tc.wantErr == nil {
				assert.Equal(t, tc.wantUser, user)
			}
		})
	}
}
