package user

import (
	"context"
	"errors"
	"go-crud-api/m/pkg/hashutil"
	"go-crud-api/m/pkg/jwtutil"
)

type Service interface {
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, user User) (*User, error)
	Profile(ctx context.Context, username string) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (svc *service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := svc.repo.GetByUsername(ctx, username)
	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}

	if !hashutil.CheckPassword(user.Password, password) {
		return "", errors.New("invalid credentials")
	}

	token, err := jwtutil.CreateToken(user.Username, string(user.Role))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (svc *service) Register(ctx context.Context, user User) (*User, error) {
	existing, err := svc.repo.GetByUsername(ctx, user.Username)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("username already exists")
	}

	hashed, err := hashutil.HashPassword(user.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	user.Password = hashed

	newUser, err := svc.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (svc *service) Profile(ctx context.Context, username string) (*User, error) {
	user, err := svc.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}
