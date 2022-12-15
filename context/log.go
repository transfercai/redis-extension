package context

import (
	"log"
)

type Log struct {
	Level int
}

func (l *Log) Do(ctx *Context) {
	ctx.Next()
	log.Println(ctx.resp)
}
