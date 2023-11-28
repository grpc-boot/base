package components

import (
	"testing"
)

func TestContext_Next(t *testing.T) {
	ctx := AcquireCtx([]Handler{
		func(ctx *Context) {
			t.Log("1 before")

			ctx.Next()

			t.Log("1 after")
		},
		func(ctx *Context) {
			t.Log("2 before")

			ctx.Next()

			t.Log("2 after")
		},
	})

	defer ctx.Close()

	ctx.Next()
}

func TestContext_Abort(t *testing.T) {
	ctx := AcquireCtx([]Handler{
		func(ctx *Context) {
			t.Log("1 before")

			ctx.Abort()

			t.Log("1 after")
		},
		func(ctx *Context) {
			t.Log("2 before")

			ctx.Next()

			t.Log("2 after")
		},
	})

	defer ctx.Close()

	ctx.Next()
}
