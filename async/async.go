package async

import (
	"context"
	"time"
)

var DefaultTimeout = 5 * time.Second

type Async[T any] interface {
	Await() (T, error)
}

type async[T any] struct {
	value *T
	err   error
	ch    chan struct{}
	ctx   context.Context
}

func New[T any](f func() (T, error)) Async[T] {
	return NewWithTimeout(DefaultTimeout, f)
}

func NewWithTimeout[T any](timeout time.Duration, f func() (T, error)) Async[T] {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return NewWithContext(ctx, func() (T, error) {
		defer cancel()
		return f()
	})
}

func NewWithContext[T any](ctx context.Context, f func() (T, error)) Async[T] {

	task := &async[T]{
		value: new(T),
		err:   nil,
		ch:    make(chan struct{}),
		ctx:   ctx,
	}

	Do(func() {
		*task.value, task.err = f()
	}, func(err error) {

		defer close(task.ch)

		if err != nil {
			task.value, task.err = nil, err
		}
	})

	return task
}

func (a *async[T]) Await() (T, error) {
	select {
	case <-a.ctx.Done():
		return *new(T), a.ctx.Err()
	case <-a.ch:
		return *a.value, a.err
	}
}
