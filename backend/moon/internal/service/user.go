package service

import (
	"context"
	"errors"
	"moon/internal/domain"
	"moon/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = errors.New("邮箱冲突")
	ErrInvalidUserOrPassword = errors.New("用户不存在或者密码不对")
)

type UserService interface {
	Signup(ctx context.Context, email, password, nickname string) error
	Login(ctx context.Context, email, password string) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Signup(ctx context.Context, email, password, nickname string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := domain.User{
		Email:    email,
		Password: string(hashedPassword),
		Nickname: nickname,
	}

	err = s.repo.Create(ctx, user)
	if err == repository.ErrDuplicateUser {
		return ErrDuplicateEmail
	}
	return err
}

func (s *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := s.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}
