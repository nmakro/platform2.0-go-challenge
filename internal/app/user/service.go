package user

import (
	"context"

	"github.com/nmakro/platform2.0-go-challenge/internal/pkg/pagination"
)

type Service interface {
	GetAllUsers(ctx context.Context, pg pagination.PageInfo) ([]User, error)
	AddUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, user User) error
}

type UserService struct {
	repo Repository
}

func NewService(repo Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) AddUser(ctx context.Context, user User) error {
	return s.repo.Add(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, userID uint32) (User, error) {
	return s.repo.Get(ctx, userID)
}

func (s *UserService) GetAllUsers(ctx context.Context, pg pagination.PageInfo) ([]User, error) {
	return s.repo.List(ctx, pg)
}

func (s *UserService) DeleteUser(ctx context.Context, userID uint32) error {
	return s.repo.Delete(ctx, userID)
}
