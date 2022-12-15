package context

import (
	"math"
	"net/http"
	"sync/atomic"
)

type Context struct {
	index       *int64
	count       int64
	resp        http.ResponseWriter
	req         *http.Request
	host        string
	path        string
	serviceID   string
	HandleChain []Handler
}

func NewContext(serviceID string) *Context {
	index := int64(0)
	return &Context{
		index:       &index,
		count:       0,
		HandleChain: make([]Handler, 0),
		serviceID:   serviceID,
	}
}

func (ctx *Context) InjectHandler(h ...Handler) {
	ctx.HandleChain = append(ctx.HandleChain, h...)
	ctx.count = int64(len(ctx.HandleChain))
}

func (ctx *Context) Next() {
	for *ctx.index < ctx.count {
		ctx.HandleChain[*ctx.index].Do(ctx)
		atomic.AddInt64(ctx.index, 1)
	}
}

func (ctx *Context) Abort() {
	*ctx.index = math.MaxInt64
}

func (ctx *Context) GetService() string {
	return ctx.serviceID
}
