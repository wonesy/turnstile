package turnstile

import (
	"context"
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

	ts, nc = NewTurnstile(nil, 100)
	assert.Nil(t, ts)
	assert.Nil(t, nc)
}

func TestGo_OK(t *testing.T) {
	ctx := context.Background()
	ts, _ := NewTurnstile(ctx, 5)
	var mutex = &sync.RWMutex{}

	results := make(map[int]bool)

	for i := 0; i < 100; i++ {
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
	assert.Equal(t, 100, len(results))
}
