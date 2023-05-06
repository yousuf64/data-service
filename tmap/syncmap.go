package tmap

import (
	"log"
	"sync"
)

type Map[K comparable, V any] struct {
	smap sync.Map
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		smap: sync.Map{},
	}
}

func (sm *Map[K, V]) Set(key K, value V) {
	sm.smap.Store(key, value)
}

func (sm *Map[K, V]) SetOnce(key K, value V) (actual V, exists bool) {
	v, loaded := sm.smap.LoadOrStore(key, value)
	return v.(V), !loaded
}

func (sm *Map[K, V]) Get(key K) (value V, ok bool) {
	v, ok := sm.smap.Load(key)
	if !ok {
		return value, ok
	}
	return v.(V), ok
}

func (sm *Map[K, V]) Delete(key K) bool {
	log.Println("Trying to delete")
	sm.smap.Delete(key)
	return true
}

func (sm *Map[K, V]) Range(fn func(key K, value V) bool) {
	sm.smap.Range(func(key, value any) bool {
		return fn(key.(K), value.(V))
	})
}
