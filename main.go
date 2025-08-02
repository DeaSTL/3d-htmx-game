package main

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/a-h/templ"
	"github.com/deastl/htmx-doom/assethandlers"
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
	prefabs := assethandlers.GeneratePrefabs()

	assethandlers.TransformPrefabs(&prefabs)
	assethandlers.CreateTestPrefabFiles(prefabs)

	app.Static("/public", "./public/")
	sServer := hx.NewServer(app)
	sServer.Mount("/ws")

	//player sync
	go func() {

		for {

			playerCount := 0
			for _,p := range gameMap.Players {

				if p.Exited {
					continue
				}

				playerCount++

				log.Printf("Player: %+v Position: %+v\n", p.ID, p.Position)
				
				buffer := new(bytes.Buffer)

				err := views.PlayerSync(gameMap,p.ID).
					Render(context.Background(), buffer)

				if err != nil {
					log.Printf("Error rendering player sync")
					continue
				}

				p.Lock()
				err = p.Socket.WriteMessage(buffer.Bytes())
				p.Unlock()

				if err != nil {
					continue
				}

			}

			log.Printf("Current Player Count: %+v",playerCount)
			time.Sleep(time.Millisecond * 250)
		}
	}()

	sServer.Listen("main", func(c *hx.Client, msg []byte) {
		log.Printf("Client %s", c.ID)
	})
	network.RegisterPlayerMessageHandlers(&sServer, &gameMap)

	// Root endpoint - serves a page with an iframe
	app.Get("/", func(c *fiber.Ctx) error {
		return Render(c, views.Index())
	})

	// Game endpoint - serves the actual game content
	app.Get("/game", func(c *fiber.Ctx) error {
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
