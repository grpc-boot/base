package components

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

type Event struct {
	name string
	data any
}

func (e *Event) Name() string {
	return e.name
}

func (e *Event) Data() any {
	return e.data
}
