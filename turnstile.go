package turnstile

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// Turnstile ...
type Turnstile struct {
	group     *errgroup.Group
	semaphore chan struct{}
}

// NewTurnstile ...
func NewTurnstile(ctx context.Context, size uint32) (*Turnstile, context.Context) {
	g, c := errgroup.WithContext(ctx)
	return &Turnstile{
		group:     g,
		semaphore: make(chan struct{}, size),
	}, c
}

// Go ...
func (t *Turnstile) Go(f func() error) {
	t.semaphore <- struct{}{}
	t.group.Go(func() error {
		defer func() {
			<-t.semaphore
		}()
		return f()
	})
}

// Wait ...
func (t *Turnstile) Wait() error {
	close(t.semaphore)
	return t.group.Wait()
}
