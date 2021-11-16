package user

import (
	"context"
)

// Repository interface for user.
type Repository interface {
	GetByEmail(ctx context.Context, userEmail string) (User, error)
	GetWithPassword(ctx context.Context, userEmail string) (AddUserCommand, error)
	Add(ctx context.Context, user AddUserCommand) error
	Delete(ctx context.Context, userEmail string) error
	List(ctx context.Context) ([]User, error)
}
