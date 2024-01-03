package components

type Chain struct {
	handlers []Handler
}

func NewChain(handlers ...Handler) *Chain {
	c := &Chain{
		handlers: handlers,
	}

	return c
}

func (c *Chain) Use(handlers ...Handler) {
	if c.handlers == nil {
		c.handlers = handlers
		return
	}

	c.handlers = append(c.handlers, handlers...)
}

func (c *Chain) Run() {
	if len(c.handlers) == 0 {
		return
	}

	ctx := AcquireCtx(c.handlers)
	c.RunWithCtx(ctx)
}

func (c *Chain) RunWithCtx(ctx *Context) {
	if len(c.handlers) == 0 {
		return
	}

	defer ctx.Close()
	ctx.handlers = c.handlers
	ctx.Next()
}
