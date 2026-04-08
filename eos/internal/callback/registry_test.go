package callback

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_oneshot_should_return_result_on_complete(t *testing.T) {
	cb := NewOneShot()
	defer cb.Delete()

	go func() {
		cb.Complete(OneShotResult{ResultCode: 42, Data: "hello"})
	}()

	result, err := cb.Wait(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 42, result.ResultCode)
	assert.Equal(t, "hello", result.Data)
}

func Test_oneshot_should_return_error_on_context_cancel(t *testing.T) {
	cb := NewOneShot()
	defer cb.Delete()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := cb.Wait(ctx)
	assert.Error(t, err)
}

func Test_oneshot_delete_should_not_panic(t *testing.T) {
	cb := NewOneShot()
	assert.NotPanics(t, func() { cb.Delete() })
}

func Test_complete_by_handle_should_dispatch_to_oneshot(t *testing.T) {
	cb := NewOneShot()
	defer cb.Delete()

	go func() {
		CompleteByHandle(cb.Handle(), OneShotResult{ResultCode: 7})
	}()

	result, err := cb.Wait(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 7, result.ResultCode)
}

func Test_notification_should_call_registered_callback(t *testing.T) {
	reg := NewNotificationRegistry()
	var receivedData any

	reg.Register(100, func(data any) { receivedData = data })
	reg.Dispatch(100, "test-data")

	assert.Equal(t, "test-data", receivedData)
}

func Test_notification_should_not_call_after_unregister(t *testing.T) {
	reg := NewNotificationRegistry()
	called := false

	reg.Register(200, func(data any) { called = true })
	reg.Unregister(200)
	reg.Dispatch(200, nil)

	assert.False(t, called)
}

func Test_notification_dispatch_unknown_id_should_not_panic(t *testing.T) {
	reg := NewNotificationRegistry()
	assert.NotPanics(t, func() { reg.Dispatch(999, nil) })
}

func Test_notification_count_should_reflect_registrations(t *testing.T) {
	reg := NewNotificationRegistry()
	assert.Equal(t, 0, reg.Count())

	reg.Register(1, func(data any) {})
	reg.Register(2, func(data any) {})
	assert.Equal(t, 2, reg.Count())

	reg.Unregister(1)
	assert.Equal(t, 1, reg.Count())
}

func Test_notification_should_handle_concurrent_access(t *testing.T) {
	reg := NewNotificationRegistry()
	var wg sync.WaitGroup

	for i := range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := uint64(i)
			reg.Register(id, func(data any) {})
			reg.Dispatch(id, "data")
			reg.Unregister(id)
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for concurrent operations")
	}
}
