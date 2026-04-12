package threadworker

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func startWorker(t *testing.T, tickFn TickFunc, opts ...Option) (*Worker, context.CancelFunc) {
	t.Helper()
	w := New(tickFn, opts...)
	ctx, cancel := context.WithCancel(context.Background())
	w.Start(ctx)
	return w, func() { cancel(); w.Stop() }
}

func Test_worker_should_start_and_stop(t *testing.T) {
	w, cleanup := startWorker(t, func() {})
	assert.True(t, w.IsRunning())

	cleanup()
	assert.False(t, w.IsRunning())
}

func Test_worker_should_execute_submitted_work(t *testing.T) {
	w, cleanup := startWorker(t, func() {})
	defer cleanup()

	executed := false
	err := w.Submit(func() { executed = true })

	assert.NoError(t, err)
	assert.True(t, executed)
}

func Test_worker_should_execute_in_fifo_order(t *testing.T) {
	w, cleanup := startWorker(t, func() {}, WithTickInterval(100*time.Millisecond))
	defer cleanup()

	var order []int
	var mu sync.Mutex

	for i := range 5 {
		i := i
		err := w.Submit(func() {
			mu.Lock()
			order = append(order, i)
			mu.Unlock()
		})
		assert.NoError(t, err)
	}

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []int{0, 1, 2, 3, 4}, order)
}

func Test_worker_should_tick(t *testing.T) {
	var tickCount atomic.Int64
	_, cleanup := startWorker(t, func() { tickCount.Add(1) }, WithTickInterval(1*time.Millisecond))

	time.Sleep(50 * time.Millisecond)
	cleanup()

	assert.Greater(t, tickCount.Load(), int64(0))
}

func Test_worker_should_return_error_after_stop(t *testing.T) {
	w, cleanup := startWorker(t, func() {})
	cleanup()

	err := w.Submit(func() {})
	assert.Error(t, err)
}

func Test_worker_should_respect_context_cancellation_on_submit(t *testing.T) {
	w, cleanup := startWorker(t, func() {}, WithTickInterval(100*time.Millisecond))
	defer cleanup()

	submitCtx, submitCancel := context.WithCancel(context.Background())
	submitCancel()

	err := w.SubmitWithContext(submitCtx, func() {
		time.Sleep(1 * time.Second)
	})
	assert.Error(t, err)
}

func Test_worker_should_shutdown_on_context_cancel(t *testing.T) {
	w := New(func() {})
	ctx, cancel := context.WithCancel(context.Background())
	w.Start(ctx)

	cancel()
	time.Sleep(50 * time.Millisecond)

	assert.False(t, w.IsRunning())
}

func Test_worker_should_handle_concurrent_submits(t *testing.T) {
	w, cleanup := startWorker(t, func() {}, WithTickInterval(100*time.Millisecond))
	defer cleanup()

	var counter atomic.Int64
	var wg sync.WaitGroup

	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := w.Submit(func() { counter.Add(1) })
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
	assert.Equal(t, int64(100), counter.Load())
}

// Simulates a Cgo notification callback fired from inside EOS_Platform_Tick
// which then calls back into the SDK via Submit. Before the re-entrance
// guard, the nested Submit would block on workCh forever because the
// worker goroutine is sitting in tickFn waiting for its own done channel.
func Test_worker_submit_from_inside_tickFn_should_run_inline(t *testing.T) {
	var w *Worker
	var nestedRan atomic.Bool
	var nestedSubmitErr atomic.Value

	done := make(chan struct{})
	tickFn := func() {
		if nestedRan.Load() {
			return
		}
		err := w.Submit(func() { nestedRan.Store(true) })
		if err != nil {
			nestedSubmitErr.Store(err)
		}
		close(done)
	}

	w = New(tickFn, WithTickInterval(1*time.Millisecond))
	ctx, cancel := context.WithCancel(context.Background())
	w.Start(ctx)
	defer func() { cancel(); w.Stop() }()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("tickFn never completed — nested Submit deadlocked")
	}

	assert.True(t, nestedRan.Load(), "nested fn should have run")
	assert.Nil(t, nestedSubmitErr.Load(), "nested Submit should return nil")
}

// Nested Submits from inside a work item (not a tick) must also run
// inline — same deadlock shape, different entry point.
func Test_worker_submit_from_inside_work_item_should_run_inline(t *testing.T) {
	w, cleanup := startWorker(t, func() {}, WithTickInterval(100*time.Millisecond))
	defer cleanup()

	var nestedRan atomic.Bool
	done := make(chan error, 1)

	go func() {
		done <- w.Submit(func() {
			// Inside a work item — we're on the worker goroutine.
			nestedErr := w.Submit(func() { nestedRan.Store(true) })
			if nestedErr != nil {
				done <- nestedErr
			}
		})
	}()

	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(1 * time.Second):
		t.Fatal("nested Submit from work item deadlocked")
	}

	assert.True(t, nestedRan.Load())
}
