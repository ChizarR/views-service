package intaraction

import (
	"context"
)

type Storage interface {
	GetOrCreate(ctx context.Context, date string) (Intaraction, error)
	Create(ctx context.Context, intr Intaraction) (string, error)
	FindOne(ctx context.Context, date string) (Intaraction, error)
	Update(ctx context.Context, intr Intaraction) error
	FindAll(ctx context.Context) ([]Intaraction, error)
}
