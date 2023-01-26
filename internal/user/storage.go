package user

import (
	"context"
)

type Storage interface {
	GetOrCreate(ctx context.Context, tgId int) (User, error)
	Create(ctx context.Context, user User) (string, error)
	FindOne(ctx context.Context, tgId int) (User, error)
	FindAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, user User) error
}
