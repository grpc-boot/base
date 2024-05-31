package main

import (
	"flag"
	"os"

	"app/controllers"
	"app/lib/conf"

	"github.com/gofiber/fiber/v3"
	"github.com/grpc-boot/base/v2/grace"
	"github.com/grpc-boot/base/v2/utils"
)

func init() {
	var confPath string
	flag.StringVar(&confPath, "c", "", "config file path")
	flag.Parse()
	if confPath == "" {
		flag.Usage()
		os.Exit(0)
	}

	if err := conf.LoadConfig(confPath); err != nil {
		utils.RedFatal("failed to load config: %v", err)
	}
}

func main() {
	engine := fiber.New()
	controllers.LoadRouter(engine)

	gf := grace.New(engine.Server(), nil)
	if err := gf.Serve(":8080", "8081"); err != nil {
		utils.RedFatal("failed to start the server: %v", err)
	}
}
