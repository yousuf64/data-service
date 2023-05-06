package worker

import (
	"log"
	"sync"
)

type Worker[T any] struct {
	c       chan struct{}
	fn      func() T
	mu      sync.RWMutex
	started uint8
	data    T
}

func New[T any](fn func() T) *Worker[T] {
	return &Worker[T]{
		c:       make(chan struct{}),
		fn:      fn,
		mu:      sync.RWMutex{},
		started: 0,
	}
}

func (wkr *Worker[T]) start() {
	wkr.mu.Lock()
	if wkr.started == 0 {
		wkr.started = 1
		go func() {
			wkr.data = wkr.fn()
			wkr.mu.Lock()
			wkr.started = 0
			wkr.mu.Unlock()
			close(wkr.c)
			wkr.c = make(chan struct{})
		}()
	}
	wkr.mu.Unlock()
}

type Subscriber[T any] struct {
	c    chan struct{}
	stop chan struct{}
}

func newSubscriber[T any](wkr *Worker[T], fn func(T)) *Subscriber[T] {
	sub := &Subscriber[T]{
		c:    wkr.c,
		stop: make(chan struct{}),
	}

	go func() {
		select {
		case <-wkr.c:
			log.Println("received", wkr.data)
			fn(wkr.data)
		case <-sub.stop:
		}
	}()

	return sub
}

func (s *Subscriber[T]) Unsubscribe() {
	close(s.stop)
}

func (wkr *Worker[T]) Subscribe(fn func(T)) (*Subscriber[T], bool) {
	wkr.mu.RLock()
	switch wkr.started {
	case 0:
		log.Println("starting")
		wkr.mu.RUnlock()
		wkr.start()
		log.Println("started")
		return newSubscriber(wkr, fn), true
	case 1:
		defer wkr.mu.RUnlock()
		return newSubscriber(wkr, fn), true
	default:
		defer wkr.mu.RUnlock()
		return nil, false
	}
}
