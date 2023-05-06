package worker

import (
	"github.com/google/uuid"
	"log"
	"sync"
)

type Worker[T any] struct {
	id      uuid.UUID
	c       chan struct{}
	fn      func() T
	mu      sync.RWMutex
	running bool
	data    T
}

func New[T any](fn func() T) *Worker[T] {
	return &Worker[T]{
		id:      uuid.New(),
		c:       nil,
		fn:      fn,
		mu:      sync.RWMutex{},
		running: false,
	}
}

func (wkr *Worker[T]) bootstrap() {
	wkr.mu.Lock()
	if !wkr.running {
		log.Printf("[DEBUG] <WID: %s> Bootstrapping\n", wkr.id.String())
		wkr.c = make(chan struct{})
		wkr.running = true
		go func() {
			log.Printf("[DEBUG] <WID: %s> Executing worker\n", wkr.id.String())
			wkr.data = wkr.fn()
			wkr.mu.Lock()
			wkr.running = false
			log.Printf("[DEBUG] <WID: %s> Attempting to close worker\n", wkr.id.String())
			close(wkr.c)
			wkr.mu.Unlock()
		}()
	}
	wkr.mu.Unlock()
}

// Subscribe subscribes to the worker and executes the callback function on worker completion.
// Worker starts its execution upon the first subscription.
// If the worker has been executed already, it re-executes the worker.
// Therefore, workers are reusable.
func (wkr *Worker[T]) Subscribe(fn func(T)) (subscriber *Subscriber[T]) {
	defer func() {
		log.Printf("[DEBUG] <WID: %s, SID: %s> Subscribed to worker\n", wkr.id.String(), subscriber.subscriberId.String())
	}()

	wkr.mu.RLock()
	if !wkr.running {
		wkr.mu.RUnlock()
		wkr.bootstrap()
		return newSubscriber(wkr, fn)
	}

	defer wkr.mu.RUnlock()
	return newSubscriber(wkr, fn)
}

func (wkr *Worker[T]) Id() uuid.UUID {
	return wkr.id
}

type Subscriber[T any] struct {
	subscriberId uuid.UUID
	workerId     uuid.UUID
	c            chan struct{}
	stop         chan struct{}
}

func newSubscriber[T any](wkr *Worker[T], fn func(T)) *Subscriber[T] {
	s := &Subscriber[T]{
		subscriberId: uuid.New(),
		workerId:     wkr.id,
		c:            wkr.c,
		stop:         make(chan struct{}),
	}

	go func() {
		select {
		case <-wkr.c:
			log.Printf("[DEBUG] <WID: %s, SID: %s> Received %v\n", wkr.id.String(), s.subscriberId.String(), wkr.data)
			fn(wkr.data)
		case <-s.stop:
			log.Printf("[DEBUG] <WID: %s, SID: %s> Received stop signal\n", wkr.id.String(), s.subscriberId.String())
		}
	}()

	return s
}

func (s *Subscriber[T]) Unsubscribe() {
	close(s.stop)
}
