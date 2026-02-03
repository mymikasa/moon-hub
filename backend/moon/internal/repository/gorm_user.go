package repository

import (
	"context"
	"database/sql"
	"moon/internal/domain"
	"moon/internal/repository/dao"
	"time"
)

type GORMUserRepository struct {
	dao dao.UserDAO
}

func NewGORMUserRepository(dao dao.UserDAO) UserRepository {
	return &GORMUserRepository{dao: dao}
}

func (r *GORMUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, domainToDaoUser(u))
}

func (r *GORMUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	daoUser, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return daoToDomainUser(daoUser), nil
}

func (r *GORMUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	daoUser, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return daoToDomainUser(daoUser), nil
}

func (r *GORMUserRepository) Update(ctx context.Context, u domain.User) error {
	return r.dao.Update(ctx, domainToDaoUser(u))
}

func domainToDaoUser(u domain.User) dao.User {
	return dao.User{
		Id:       u.Id,
		Email:    sql.NullString{String: u.Email, Valid: u.Email != ""},
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Phone:    sql.NullString{String: u.Phone, Valid: u.Phone != ""},
		Ctime:    u.Ctime.UnixMilli(),
	}
}

func daoToDomainUser(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
		AboutMe:  u.AboutMe,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}
