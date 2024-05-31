package base

import (
	"github.com/gofiber/fiber/v3"
	"github.com/grpc-boot/base/v2/components"
)

type Controller struct{}

func (c *Controller) Status(ctx fiber.Ctx, sts *components.Status) (int, error) {
	defer sts.Close()
	return c.Json(ctx, sts.JsonMarshal())
}

func (c *Controller) Json(ctx fiber.Ctx, data []byte) (int, error) {
	ctx.Response().Header.SetContentType(fiber.MIMEApplicationJSONCharsetUTF8)
	return ctx.Write(data)
}

func (c *Controller) File(ctx fiber.Ctx, data []byte) (int, error) {
	ctx.Response().Header.SetContentType(fiber.MIMEOctetStream)
	return ctx.Write(data)
}
