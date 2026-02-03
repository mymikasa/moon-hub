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
	FindById(ctx context.Context, id int64) (domain.User, error)
	Update(ctx context.Context, u domain.User) error
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

func (s *userService) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := s.repo.FindById(ctx, id)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, err
}

func (s *userService) Update(ctx context.Context, u domain.User) error {
	return s.repo.Update(ctx, u)
}
