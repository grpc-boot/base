package connctx

import "sync"

var ctxPool = sync.Pool{
	New: func() any {
		return newCtx()
	},
}

func AcquireCtx() Context {
	c := ctxPool.Get().(Context)
	return c
}

func ReleaseCtx(c Context) {
	c.Close()
	ctxPool.Put(c)
}
