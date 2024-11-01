/*
* Copyright © 2024 Minand Manomohanan <minand.nell.mohan@gmail.com>
 */
package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

var (
	urlCache Cache
)

func SetUpServer(port, origin string) {
	fmt.Println("Setting up server...")
	app := fiber.New()
	initCache()
	setUpRoutes(app, origin)

	ctx, cancel := context.WithCancel(context.Background())
	go startServer(app, &port)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-ctx.Done():
		fmt.Println("Shutting down server...")
		app.Shutdown()
	case <-sigterm:
		fmt.Println("Shutting down server...")
		app.Shutdown()
	}
	cancel()
}

func setUpRoutes(app *fiber.App, origin string) {
	fmt.Println("Setting up routes...")

	// All get requests will be routed through this handler
	app.Get("/*", func(c *fiber.Ctx) error {
		return HandlerGet(c, &origin)
	})
}

func startServer(app *fiber.App, port *string) {
	fmt.Println("Listening on port " + *port)
	listenPort := fmt.Sprintf(":%s", *port)
	err := app.Listen(listenPort)
	if err != nil {
		fmt.Println("Error starting cache server!")
	}
}

func initCache() {
	urlCache.UrlMap = make(map[string]*Response)
}
