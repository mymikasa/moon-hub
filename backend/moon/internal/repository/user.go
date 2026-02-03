package repository

import (
	"context"
	"moon/internal/domain"
	"moon/internal/repository/dao"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateEmail
	ErrUserNotFound  = dao.ErrRecordNotFound
)

//go:generate mockgen -source=./user.go -package=repomocks -destination=./mocks/user.mock.go UserRepository
type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}

type CachedUserRepository struct {
	dao  dao.UserDAO
	repo UserRepository
}

// Create implements [UserRepository].
func (c *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return c.repo.Create(ctx, u)
}

// FindByEmail implements [UserRepository].
func (c *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	return c.repo.FindByEmail(ctx, email)
}

func NewCachedUserRepository(dao dao.UserDAO, repo UserRepository) UserRepository {
	return &CachedUserRepository{
		dao:  dao,
		repo: repo,
	}
}
