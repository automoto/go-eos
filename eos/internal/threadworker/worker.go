package threadworker

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var ErrWorkerStopped = errors.New("worker is stopped")

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
	wg           sync.WaitGroup
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
	w.workCh = make(chan workItem, w.chanSize)
	ctx, w.cancel = context.WithCancel(ctx)
	w.running.Store(true)
	w.wg.Add(1)

	go func() {
		defer w.wg.Done()
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		defer w.running.Store(false)

		ticker := time.NewTicker(w.tickInterval)
		defer ticker.Stop()

		for {
			select {
			case item := <-w.workCh:
				item.fn()
				close(item.done)
			case <-ticker.C:
				w.tickFn()
			case <-ctx.Done():
				w.drain()
				return
			}
		}
	}()
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

func (w *Worker) IsRunning() bool {
	return w.running.Load()
}
