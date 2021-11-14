package user

import (
	"context"
)

// Repository interface for user.
type Repository interface {
	GetByEmail(ctx context.Context, userEmail string) (User, error)
	Add(ctx context.Context, user User) error
	Delete(ctx context.Context, userEmail string) error
}
