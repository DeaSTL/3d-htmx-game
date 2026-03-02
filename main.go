package main

import (
	"bytes"
	"html/template"
	"log"
	"time"

	"github.com/deastl/htmx-doom/assethandlers"
	"github.com/deastl/htmx-doom/gameobjects"
	"github.com/deastl/htmx-doom/network"
	"github.com/deastl/hxsocketsfiber"
	"github.com/gofiber/fiber/v2"
)

var templates *template.Template

func main() {
	// Load all templates
	var err error
	templates, err = template.New("").Funcs(template.FuncMap{
		"seq": func(start, end int) []int {
			var result []int
			for i := start; i < end; i++ {
				result = append(result, i)
			}
			return result
		},
		"mul": func(a, b interface{}) float64 {
			var aFloat, bFloat float64
			switch v := a.(type) {
			case int:
				aFloat = float64(v)
			case float64:
				aFloat = v
			case float32:
				aFloat = float64(v)
			}
			switch v := b.(type) {
			case int:
				bFloat = float64(v)
			case float64:
				bFloat = v
			case float32:
				bFloat = float64(v)
			}
			return aFloat * bFloat
		},
		"int": func(f float64) int { return int(f) },
		"sub": func(a, b interface{}) float64 {
			var aFloat, bFloat float64
			switch v := a.(type) {
			case int:
				aFloat = float64(v)
			case float64:
				aFloat = v
			case float32:
				aFloat = float64(v)
			}
			switch v := b.(type) {
			case int:
				bFloat = float64(v)
			case float64:
				bFloat = v
			case float32:
				bFloat = float64(v)
			}
			return aFloat - bFloat
		},
		"ne": func(a, b string) bool { return a != b },
		"not": func(b bool) bool { return !b },
		"dict": func(values ...interface{}) map[string]interface{} {
			if len(values)%2 != 0 {
				return nil
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil
				}
				dict[key] = values[i+1]
			}
			return dict
		},
	}).ParseGlob("views/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

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
			for _, p := range gameMap.Players {

				if p.Exited {
					continue
				}

				playerCount++


				buffer := new(bytes.Buffer)

				err := templates.ExecuteTemplate(buffer, "playerSync", map[string]interface{}{
					"Players":  gameMap.Players,
					"PlayerID": p.ID,
				})

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
			time.Sleep(time.Millisecond * 250)
		}
	}()
//status reporter
	go func() {
		maxPlayerCount := 0
		for {

			playerCount := 0

			for _, p := range gameMap.Players {

				if p.Exited {
					continue
				}

				log.Printf("Player: %+v Position: %+v\n", p.ID, p.Position)

				playerCount++
			}

			if playerCount > maxPlayerCount {
				maxPlayerCount = playerCount
			}

			log.Printf("Current player count: %+v max player count: %+v", playerCount, maxPlayerCount)

			time.Sleep(time.Millisecond * 750)
		}

	}()

	sServer.Listen("main", func(c *hx.Client, msg []byte) {
		log.Printf("Client %s", c.ID)
	})
	network.RegisterPlayerMessageHandlers(&sServer, &gameMap, templates)

	// Root endpoint - serves a page with an iframe
	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		var buf bytes.Buffer
		err := templates.ExecuteTemplate(&buf, "index", nil)
		if err != nil {
			return err
		}
		return c.Send(buf.Bytes())
	})

	// Game endpoint - serves the actual game content
	app.Get("/game", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		var buf bytes.Buffer
		err := templates.ExecuteTemplate(&buf, "main", nil)
		if err != nil {
			return err
		}
		return c.Send(buf.Bytes())
	})

	log.Fatal(app.Listen(":3000"))

}
