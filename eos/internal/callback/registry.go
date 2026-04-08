package callback

import (
	"context"
	"runtime/cgo"
	"sync"
	"sync/atomic"
)

type OneShotResult struct {
	ResultCode int
	Data       any
}

type OneShotCallback struct {
	resultCh chan OneShotResult
	handle   cgo.Handle
	deleted  atomic.Bool
}

func NewOneShot() *OneShotCallback {
	cb := &OneShotCallback{
		resultCh: make(chan OneShotResult, 1),
	}
	cb.handle = cgo.NewHandle(cb)
	return cb
}

func (c *OneShotCallback) Handle() cgo.Handle {
	return c.handle
}

// HandleValue returns the handle as a uintptr suitable for passing as
// C void* ClientData. The caller is responsible for the unsafe.Pointer
// conversion at the Cgo boundary where it is permitted.
func (c *OneShotCallback) HandleValue() uintptr {
	return uintptr(c.handle)
}

func (c *OneShotCallback) Wait(ctx context.Context) (OneShotResult, error) {
	select {
	case result := <-c.resultCh:
		return result, nil
	case <-ctx.Done():
		return OneShotResult{}, ctx.Err()
	}
}

func (c *OneShotCallback) Complete(result OneShotResult) {
	c.resultCh <- result
}

// Delete frees the cgo.Handle. Safe to call multiple times.
func (c *OneShotCallback) Delete() {
	if c.deleted.CompareAndSwap(false, true) {
		c.handle.Delete()
	}
}

// CompleteByHandle completes the oneshot callback identified by handle
// and frees the handle. Called from C callback trampolines.
func CompleteByHandle(handle cgo.Handle, result OneShotResult) {
	cb := handle.Value().(*OneShotCallback)
	cb.Complete(result)
	cb.Delete()
}

type NotifyFunc func(data any)

type RemoveNotifyFunc func()

type NotificationRegistry struct {
	mu        sync.RWMutex
	callbacks map[uint64]NotifyFunc
}

func NewNotificationRegistry() *NotificationRegistry {
	return &NotificationRegistry{
		callbacks: make(map[uint64]NotifyFunc),
	}
}

func (r *NotificationRegistry) Register(id uint64, fn NotifyFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.callbacks[id] = fn
}

func (r *NotificationRegistry) Dispatch(id uint64, data any) {
	r.mu.RLock()
	fn, ok := r.callbacks[id]
	r.mu.RUnlock()

	if ok {
		fn(data)
	}
}

func (r *NotificationRegistry) Unregister(id uint64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.callbacks, id)
}

func (r *NotificationRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.callbacks)
}
