package user

import (
	"context"

	"github.com/nmakro/platform2.0-go-challenge/internal/pkg/pagination"
)

// Repository interface for user.
type Repository interface {
	GetUser(ctx context.Context, userID uint32) (User, error)
	List(ctx context.Context, pageInfo pagination.PageInfo) ([]User, error)
	Add(ctx context.Context, user User) error
	Delete(ctx context.Context, user User) error
}
