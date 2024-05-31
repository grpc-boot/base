package controllers

import "github.com/gofiber/fiber/v3"

func LoadRouter(engine *fiber.App) {
	index := &indexController{}
	engine.Get("/", index.Index)
	engine.Get("/index", index.Index)
}
