package context

type Handler interface {
	Do(ctx *Context)
}
