package micro

import (
	"context"

	"github.com/pkg/errors"
)

var ErrLocked = errors.New("locked")

type Locker interface {
	TryLock(ctx context.Context, key string) (error, bool)
	Lock(ctx context.Context, key string) error
	Unlock(ctx context.Context, key string) error
}
