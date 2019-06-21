package turnstile

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTurnstile(t *testing.T) {
	ctx := context.Background()

	ts, nc := NewTurnstile(ctx, 0)
	assert.NotNil(t, ts)
	assert.NotNil(t, nc)
	assert.Equal(t, 0, cap(ts.semaphore))

	ts, nc = NewTurnstile(ctx, 100)
	assert.NotNil(t, ts)
	assert.NotNil(t, nc)
	assert.Equal(t, 100, cap(ts.semaphore))
}

func TestGo_OK(t *testing.T) {
	ctx := context.Background()
	ts, _ := NewTurnstile(ctx, 5)
	var mutex = &sync.RWMutex{}

	size := 10000
	results := make(map[int]bool)

	for i := 0; i < size; i++ {
		i := i
		ts.Go(func() error {
			mutex.Lock()
			results[i] = true
			mutex.Unlock()
			return nil
		})
	}

	err := ts.Wait()

	assert.Nil(t, err)
	assert.Equal(t, size, len(results))
}

func TestGo_OK_Cancel(t *testing.T) {
	ctx := context.Background()
	ts, _ := NewTurnstile(ctx, 5)

	size := 10000

	test := 0
	for ; test < size; test++ {
		test := test
		ts.Go(func() error {
			if test == 10 {
				return errors.New("badtest")
			}
			return nil
		})
	}

	err := ts.Wait()
	assert.Error(t, err)
	assert.Equal(t, size, test)
}
