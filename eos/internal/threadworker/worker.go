package threadworker

import (
	"bytes"
	"context"
	"errors"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var ErrWorkerStopped = errors.New("worker is stopped")

// goroutineID returns the current goroutine's ID by parsing runtime.Stack.
// Used only for worker re-entrance detection in Submit. The 64-byte buffer
// is stack-allocated — the only line we care about is
// "goroutine NNN [running]:\n" which fits comfortably.
func goroutineID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	b := buf[:n]
	const prefix = "goroutine "
	if !bytes.HasPrefix(b, []byte(prefix)) {
		return 0
	}
	b = b[len(prefix):]
	space := bytes.IndexByte(b, ' ')
	if space < 0 {
		return 0
	}
	id, _ := strconv.ParseUint(string(b[:space]), 10, 64)
	return id
}

type TickFunc func()

type Option func(*Worker)

func WithTickInterval(d time.Duration) Option {
	return func(w *Worker) { w.tickInterval = d }
}

func WithWorkChannelSize(n int) Option {
	return func(w *Worker) { w.chanSize = n }
}

type workItem struct {
	fn   func()
	done chan struct{}
}

type Worker struct {
	tickFn       TickFunc
	tickInterval time.Duration
	chanSize     int
	workCh       chan workItem
	running      atomic.Bool
	cancel       context.CancelFunc
	ctx          context.Context
	wg           sync.WaitGroup
	// workerGID is the goroutine ID of the loop goroutine while it is
	// running. Submit consults it to detect re-entrant calls originating
	// from within tickFn or a previously-dispatched work item (typically
	// via a Cgo notification callback fired during EOS_Platform_Tick) and
	// executes fn inline instead of deadlocking on an empty workCh.
	workerGID atomic.Uint64
}

func New(tickFn TickFunc, opts ...Option) *Worker {
	w := &Worker{
		tickFn:       tickFn,
		tickInterval: 16 * time.Millisecond,
		chanSize:     256,
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (w *Worker) Start(ctx context.Context) {
	w.prepare(ctx)
	go w.loop()
}

// StartBlocking runs the worker loop on the calling goroutine.
// The caller must already hold runtime.LockOSThread (e.g. via init()).
// This is required on macOS where the EOS SDK's HTTP layer needs
// the main thread's run loop.
func (w *Worker) StartBlocking(ctx context.Context) {
	w.prepare(ctx)
	w.loop()
}

func (w *Worker) prepare(ctx context.Context) {
	w.workCh = make(chan workItem, w.chanSize)
	ctx, w.cancel = context.WithCancel(ctx)
	w.running.Store(true)
	w.wg.Add(1)
	w.ctx = ctx
}

func (w *Worker) loop() {
	defer w.wg.Done()
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer w.running.Store(false)

	w.workerGID.Store(goroutineID())
	defer w.workerGID.Store(0)

	ticker := time.NewTicker(w.tickInterval)
	defer ticker.Stop()

	for {
		select {
		case item := <-w.workCh:
			item.fn()
			close(item.done)
		case <-ticker.C:
			w.tickFn()
		case <-w.ctx.Done():
			w.drain()
			return
		}
	}
}

func (w *Worker) drain() {
	for {
		select {
		case item := <-w.workCh:
			item.fn()
			close(item.done)
		default:
			return
		}
	}
}

func (w *Worker) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
	w.wg.Wait()
}

func (w *Worker) Submit(fn func()) error {
	if !w.running.Load() {
		return ErrWorkerStopped
	}
	if w.isWorkerGoroutine() {
		fn()
		return nil
	}
	item := workItem{fn: fn, done: make(chan struct{})}
	select {
	case w.workCh <- item:
		<-item.done
		return nil
	default:
		return ErrWorkerStopped
	}
}

func (w *Worker) SubmitWithContext(ctx context.Context, fn func()) error {
	if !w.running.Load() {
		return ErrWorkerStopped
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if w.isWorkerGoroutine() {
		fn()
		return nil
	}

	item := workItem{fn: fn, done: make(chan struct{})}
	select {
	case w.workCh <- item:
	case <-ctx.Done():
		return ctx.Err()
	}

	select {
	case <-item.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// isWorkerGoroutine reports whether the caller is running on the worker's
// own loop goroutine — i.e. Submit was reached re-entrantly from within
// tickFn or a previously-dispatched work item (most commonly, a Cgo
// notification callback fired during EOS_Platform_Tick). In that case the
// caller already holds the worker's locked OS thread, so fn must run
// inline rather than enqueue — blocking on workCh would self-deadlock.
func (w *Worker) isWorkerGoroutine() bool {
	gid := w.workerGID.Load()
	return gid != 0 && gid == goroutineID()
}

func (w *Worker) IsRunning() bool {
	return w.running.Load()
}
