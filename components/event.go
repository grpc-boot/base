package components

import "github.com/grpc-boot/base/v3/kind"

type EventManager struct {
	events map[string]*Chain
}

func (em *EventManager) On(name string, handlers ...Handler) {
	if em.events == nil {
		em.events = map[string]*Chain{
			name: NewChain(handlers...),
		}
		return
	}

	chain, exists := em.events[name]
	if !exists {
		em.events[name] = NewChain(handlers...)
		return
	}

	chain.Use(handlers...)
}

func (em *EventManager) Has(name string) bool {
	_, exists := em.events[name]
	return exists
}

func (em *EventManager) Trigger(name string, data any) {
	chain, exists := em.events[name]
	if !exists {
		return
	}

	if len(chain.handlers) == 0 {
		return
	}

	ctx := AcquireCtx(nil)
	ctx.SetEvent(&Event{
		name: name,
		data: data,
	})

	chain.RunWithCtx(ctx)
}

func (em *EventManager) TriggerWithCtx(name string, ctx *Context) {
	chain, exists := em.events[name]
	if !exists {
		return
	}

	if len(chain.handlers) == 0 {
		return
	}

	if event := ctx.Event(); event == nil {
		ctx.SetEvent(&Event{
			name: name,
		})
	}

	chain.RunWithCtx(ctx)
}

type Event struct {
	name string
	data any
}

func NewEvent(name string, data any) *Event {
	return &Event{
		name: name,
		data: data,
	}
}

func (e *Event) Name() string {
	return e.name
}

func (e *Event) Data() any {
	return e.data
}

func (e *Event) ParamData() kind.JsonParam {
	if p, ok := e.data.(kind.JsonParam); ok {
		return p
	}

	p, _ := e.data.(map[string]any)

	return p
}
