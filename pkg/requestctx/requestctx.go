package requestctx

import (
	"context"
	"github.com/google/uuid"
)

var ctxKey uint8 = 0

type RequestCtx struct {
	context.Context

	id uuid.UUID
}

func New(ctx context.Context, id uuid.UUID) *RequestCtx {
	return &RequestCtx{
		Context: ctx,
		id:      id,
	}
}

func (ctx *RequestCtx) Value(key any) any {
	if key == ctxKey {
		return ctx
	}
	return ctx.Context.Value(key)
}

func FromContext(ctx context.Context) uuid.UUID {
	if rCtx, ok := ctx.Value(ctxKey).(*RequestCtx); ok {
		return rCtx.id
	}
	return uuid.Nil
}
