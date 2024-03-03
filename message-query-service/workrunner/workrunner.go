package workrunner

import (
	"github.com/yousuf64/request-coalescing/message-query-service/worker"
	"log"
	"sync"
)

type WorkRunner[K comparable, V any] struct {
	m  map[K]*worker.Worker[V]
	mu sync.RWMutex
}

func New[K comparable, V any]() *WorkRunner[K, V] {
	return &WorkRunner[K, V]{
		m:  make(map[K]*worker.Worker[V]),
		mu: sync.RWMutex{},
	}
}

func (wr *WorkRunner[K, V]) Observe(key K, fn func() V) <-chan V {
	c := make(chan V)
	wr.mu.RLock()
	wkr, ok := wr.m[key]
	wr.mu.RUnlock()
	if !ok {
		wr.mu.Lock()
		wkr, ok = wr.m[key] // Re-check after placing a full lock.
		if !ok {
			wkr = worker.New[V](fn)
			wr.m[key] = wkr
		}
		wr.mu.Unlock()
	}
	log.Printf("[DEBUG] <WID: %s> Observing worker\n", wkr.Id().String())
	wkr.Subscribe(func(v V) {
		// Though workers are reusable, it's better to keep workers map light as possible.
		wr.mu.RLock()
		_, ok := wr.m[key]
		wr.mu.RUnlock()
		if ok {
			wr.mu.Lock()
			_, ok = wr.m[key] // Re-check after placing a full lock.
			if ok {
				delete(wr.m, key)
			}
			wr.mu.Unlock()
		}

		c <- v
	})
	return c
}
