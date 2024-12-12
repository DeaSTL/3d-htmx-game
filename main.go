package main

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/a-h/templ"
	"github.com/deastl/htmx-doom/gameobjects"
	"github.com/deastl/htmx-doom/network"
	"github.com/deastl/htmx-doom/views"
	"github.com/deastl/hxsocketsfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)


func main(){
  app := fiber.New()
  gameMap := gameobjects.NewGameMap(gameobjects.GameMap{
  })
  app.Static("/public","./public/")
  sServer := hx.NewServer(app)
  sServer.Mount("/ws")
  
  sServer.Listen("main",func(c *hx.Client, msg []byte){
    log.Printf("Client %s",c.ID)
  })

  // sServer.OnClientDisconnect = func(c *hx.Client){
  //   c.ID = "NULL";
  // }

  network.RegisterPlayerMessageHandlers(&sServer,&gameMap)
  go func(){
    for {
      for _, player := range gameMap.Players {
        time.Sleep(time.Millisecond * 30)
        player.Update() 

        buffer := new(bytes.Buffer)
        err := views.SceneTransform(player).
        Render(context.Background(),buffer)
        if err != nil {
          log.Printf("Error rendering plane")
        }
        client := sServer.GetClient(player.ID)
        // if client.ID == "NULL" {
        //   return;
        // }
        client.WriteMessage(buffer.Bytes())
      }
    }
  }()
  // go func(){
  //   count := int(0)
  //   for{
  //     count++
  //     clients := sServer.GetAllClients()
  //     for _, client := range clients {
  //       buffer := new(bytes.Buffer)
  //       err := views.Plane(gameMap.Walls[0],
  //       ).Render(context.Background(),buffer)
  //       if err != nil {
  //         log.Printf("Error rendering plane")
  //       }
  //       client.WriteMessage(buffer.Bytes())
  //       time.Sleep(time.Millisecond * 30) 
  //     }
  //   }
  // }()



  app.Get("/", func (c *fiber.Ctx) error {
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
