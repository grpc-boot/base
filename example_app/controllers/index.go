package controllers

import (
	"app/lib/base"

	"github.com/gofiber/fiber/v3"
	"github.com/grpc-boot/base/v2/components"
)

type indexController struct {
	base.Controller
}

func (i *indexController) Index(ctx fiber.Ctx) error {
	_, err := i.Status(ctx, components.StatusOk("Hello World"))
	return err
}
