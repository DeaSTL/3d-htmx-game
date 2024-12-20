package main

import (
	"log"

	"github.com/a-h/templ"
	"github.com/deastl/htmx-doom/gameobjects"
	"github.com/deastl/htmx-doom/network"
	"github.com/deastl/htmx-doom/views"
	"github.com/deastl/hxsocketsfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func main() {
	app := fiber.New()
	gameMap := gameobjects.NewGameMap(gameobjects.GameMap{})
	app.Static("/public", "./public/")
	sServer := hx.NewServer(app)
	sServer.Mount("/ws")

	sServer.Listen("main", func(c *hx.Client, msg []byte) {
		log.Printf("Client %s", c.ID)
	})
	network.RegisterPlayerMessageHandlers(&sServer, &gameMap)
	app.Get("/", func(c *fiber.Ctx) error {
		return Render(c, views.Main())
	})

	log.Fatal(app.Listen(":3000"))

}
func Render(c *fiber.Ctx, component templ.Component, options ...func(*templ.ComponentHandler)) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}
