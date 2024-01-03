package components

import (
	"testing"
	"time"
)

func TestEventManager_Trigger(t *testing.T) {
	var em EventManager
	ok := em.Has("login")
	if ok {
		t.Fatalf("want false, got %v", ok)
	}

	em.On("login", func(ctx *Context) {
		event := ctx.Event()
		t.Logf("login at %s with data: %v", time.Now().Format(time.DateTime), event.Data())
	})

	em.On("login", func(ctx *Context) {
		t.Logf("next with login")

		ctx.Next()

		t.Logf("next done with login")
	})

	em.On("login", func(ctx *Context) {
		t.Logf("abort with login")

		ctx.Abort()

		t.Logf("abort done with login")
	})

	em.On("login", func(ctx *Context) {
		t.Logf("after abort with login")

		ctx.Next()

		t.Logf("after done with login")
	})

	em.On("login out", func(ctx *Context) {
		event := ctx.Event()
		t.Logf("login out at %s with data: %v", time.Now().Format(time.DateTime), event.Data())
	})

	for i := 0; i < 8; i++ {
		em.Trigger("login", time.Now().Format(time.DateTime))
		em.Trigger("login out", time.Now().Format(time.DateTime))
	}
}
