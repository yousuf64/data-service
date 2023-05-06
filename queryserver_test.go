package main

import (
	"context"
	"data-service/requestctx"
	"github.com/google/uuid"
	"log"
	"sync"
	"testing"
	"time"
)

type dummyKey string

func (d dummyKey) Key() string {
	return string(d)
}

func TestOrchestrator(t *testing.T) {
	orchestrator := NewOrchestrator[string]()

	var wg sync.WaitGroup
	wg.Add(10)

	var key dummyKey = "hello"
	for i := 0; i < 10; i++ {
		if i%5 == 0 {
			key = "sunday"
			time.Sleep(time.Millisecond * 600)
		}
		go func(k dummyKey) {
			ctx := requestctx.New(context.Background(), uuid.New())

			orchestrator.Run(ctx, k, func() string {
				log.Printf("%s: [Key: %s] Querying...", requestctx.FromContext(ctx).String(), k)
				time.Sleep(time.Millisecond * 700)
				return "foo"
			})
			wg.Done()
		}(key)
		if i == 0 {
			time.Sleep(time.Millisecond * 400)
		}
	}

	wg.Wait()

	log.Printf("Short circuited %v", count.Load())
	time.Sleep(2 * time.Second)
}
