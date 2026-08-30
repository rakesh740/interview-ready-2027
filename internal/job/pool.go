package job

import (
	"context"
	"sync"
)

// Pool is a generic fixed-size worker pool: `workers` goroutines pull items
// off an internal buffered queue and hand each to handler.
type Pool[T any] struct {
	handler func(context.Context, T)
	queue   chan T
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

// NewPool starts the pool immediately; workers block on the queue until
// Shutdown is called.
func NewPool[T any](workers, queueSize int, handler func(context.Context, T)) *Pool[T] {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool[T]{
		handler: handler,
		queue:   make(chan T, queueSize),
		cancel:  cancel,
	}
	p.wg.Add(workers)
	for i := 0; i < workers; i++ {
		go p.worker(ctx)
	}
	return p
}

func (p *Pool[T]) worker(ctx context.Context) {
	defer p.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case item, ok := <-p.queue:
			if !ok {
				return
			}
			p.handler(ctx, item)
		}
	}
}

// Submit enqueues an item, blocking until there's room or ctx is done.
func (p *Pool[T]) Submit(ctx context.Context, item T) error {
	select {
	case p.queue <- item:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Shutdown stops accepting work and waits for in-flight items to finish,
// up to ctx's deadline. Workers still mid-handler when ctx is done are left
// to finish on their own; Shutdown just stops waiting for them.
func (p *Pool[T]) Shutdown(ctx context.Context) error {
	close(p.queue)
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		p.cancel()
		return nil
	case <-ctx.Done():
		p.cancel()
		return ctx.Err()
	}
}
