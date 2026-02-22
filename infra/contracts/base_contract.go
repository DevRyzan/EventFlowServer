package contracts

import (
	"context"
	"errors"
)

 
var ErrNotSupported = errors.New("operation not supported")
var ErrNotFound = errors.New("not found")

// base contract for crud operations
type BaseContract[T any] interface {
	Create(ctx context.Context, entity T) error
	GetByID(ctx context.Context, id string) (T, error)
	Update(ctx context.Context, entity T) error
	Delete(ctx context.Context, id string) error
}
