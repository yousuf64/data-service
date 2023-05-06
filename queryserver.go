package main

import (
	"context"
	"data-service/query"
	"data-service/requestctx"
	"data-service/tmap"
	"fmt"
	"github.com/google/uuid"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Idle uint8 = iota
	Waiting
	Received
	Broadcast
)

type Worker[TResponse any] struct {
	workerId    uuid.UUID
	initiatorId uuid.UUID
	status      uint8
	statusMu    sync.RWMutex
	mu          sync.Mutex
	channels    map[uuid.UUID]chan TResponse
	resp        TResponse
}

func NewWorker[TResponse any](ctx context.Context) *Worker[TResponse] {
	wkr := &Worker[TResponse]{
		workerId:    uuid.New(),
		initiatorId: requestctx.FromContext(ctx),
		status:      Idle,
		statusMu:    sync.RWMutex{},
		mu:          sync.Mutex{},
		channels:    make(map[uuid.UUID]chan TResponse),
	}
	//log.Printf("Initiated worker [InitiatorId: %s, WorkerId: %s]", wkr.initiatorId.String(), wkr.workerId.String())
	return wkr
}

func (wkr *Worker[TResponse]) Start(fn func() TResponse) {
	wkr.statusMu.Lock()
	defer wkr.statusMu.Unlock()

	if wkr.status > Idle {
		//log.Printf("[WorkerId: %s] Worker already started", wkr.workerId.String())
		return
	}

	wkr.status = Waiting
	log.Printf("[InitiatorId: %s, WorkerId: %s] Starting worker", wkr.initiatorId.String(), wkr.workerId.String())

	go func() {

		wkr.resp = fn()

		wkr.status = Received
		log.Printf("[WorkerId: %s] Response received", wkr.workerId.String())

		wkr.mu.Lock()
		log.Printf("[WorkerId: %s] Found %d channels", wkr.workerId.String(), len(wkr.channels))
		count := 0
		for _, ch := range wkr.channels {
			ch <- wkr.resp
			count++
			//log.Printf("[WorkerId: %s] Broadcast to channel %s", wkr.workerId.String(), requestId.String())
		}
		wkr.status = Broadcast
		log.Printf("[WorkerId: %s] Broadcast to %d channels", wkr.workerId.String(), count)
		wkr.mu.Unlock()
	}()
}

var count atomic.Int32

func (wkr *Worker[TResponse]) Response(ctx context.Context) <-chan TResponse {
	ch := make(chan TResponse)
	if wkr.status >= Received {
		log.Printf("[WorkerId: %s] Has received a response already", wkr.workerId.String())
		go func() {
			count.Add(1)
			ch <- wkr.resp
		}()
	} else {
		wkr.mu.Lock()
		if wkr.status >= Received {
			go func() {
				count.Add(1)
				log.Printf("[WorkerId: %s] Has received a response already", wkr.workerId.String())
				ch <- wkr.resp
			}()
		} else {
			wkr.channels[requestctx.FromContext(ctx)] = ch
			log.Printf("[Id: %s] Subscribed to Worker %s", requestctx.FromContext(ctx).String(), wkr.workerId.String())
		}
		wkr.mu.Unlock()

	}
	return ch
}

type OrchestratorKey interface {
	Key() string
}

type GetMessagesKey struct {
	ChannelId int32
}

func (g GetMessagesKey) Key() string {
	return fmt.Sprintf("%d", g.ChannelId)
}

type Orchestrator[TResponse any] struct {
	workers  map[string]*Worker[TResponse]
	workers2 *tmap.Map[string, *Worker[TResponse]]
	mu       sync.Mutex
}

func NewOrchestrator[TResponse any]() *Orchestrator[TResponse] {
	return &Orchestrator[TResponse]{
		workers:  make(map[string]*Worker[TResponse]),
		workers2: tmap.New[string, *Worker[TResponse]](),
		mu:       sync.Mutex{},
	}
}

func (o *Orchestrator[TResponse]) Run(ctx context.Context, key OrchestratorKey, fn func() TResponse) TResponse {
	log.Printf("received key %v", key)

	wkrKey := key.Key()
	worker := NewWorker[TResponse](ctx)
	var exists bool
	if worker, exists = o.workers2.SetOnce(wkrKey, worker); !exists {
		worker.Start(fn)
	}
	//if wkr, ok := o.workers[wkrKey]; ok {
	//	worker = wkr
	//} else {
	//	o.mu.Lock()
	//	if wkr, ok := o.workers[wkrKey]; ok {
	//		worker = wkr
	//	} else {
	//		worker = NewWorker[TResponse](ctx)
	//		o.workers[wkrKey] = worker
	//		worker.Start(fn)
	//	}
	//	o.mu.Unlock()
	//}

	//if wkr, ok := o.workers2.Get(wkrKey); ok {
	//	worker = wkr
	//} else {
	//	var exists bool
	//	if worker, exists = o.workers2.SetOnce(wkrKey, worker); !exists {
	//		worker = NewWorker[TResponse](ctx)
	//		worker.Start(fn)
	//	}
	//}

	resp := <-worker.Response(ctx)
	//delete(o.workers, wkrKey)
	//o.workers2.Delete(wkrKey)
	return resp
}

type queryServer struct {
	query.UnimplementedQueryServer

	getMessagesOrch    *Orchestrator[[]*query.Message]
	getMessagesWorkers map[int32]*Worker[[]*query.Message]
}

func (svr queryServer) GetMessages(ctx context.Context, request *query.GetMessagesRequest) (*query.GetMessagesReply, error) {
	key := GetMessagesKey{ChannelId: request.ChannelId}

	resp := svr.getMessagesOrch.Run(ctx, key, func() []*query.Message {
		log.Printf("%s: [Key: %s] Querying...", requestctx.FromContext(ctx).String(), key.Key())
		time.Sleep(5 * time.Second)
		return []*query.Message{
			{
				FromId:  911,
				Message: "hello",
			},
			{
				FromId:  912,
				Message: "sup?",
			},
		}
	})

	return &query.GetMessagesReply{Messages: resp}, nil
}
