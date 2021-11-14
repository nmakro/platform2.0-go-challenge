package user

import (
	"context"
)

type UserService struct {
	repo Repository
}

func NewService(repo Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) AddUser(ctx context.Context, user User) error {
	return s.repo.Add(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, userEmail string) (User, error) {
	return s.repo.GetByEmail(ctx, userEmail)
}

func (s *UserService) DeleteUser(ctx context.Context, userEmail string) error {
	return s.repo.Delete(ctx, userEmail)
}
