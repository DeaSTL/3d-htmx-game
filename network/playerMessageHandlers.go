package network

import (
	"bytes"
	"context"
	"log"

	"github.com/deastl/htmx-doom/gameobjects"
	"github.com/deastl/htmx-doom/views"
	hx "github.com/deastl/hxsocketsfiber"
)

func RegisterPlayerMessageHandlers(s *hx.Server, game *gameobjects.GameMap){
  s.Listen("player_init",func(client *hx.Client, msg []byte){
    //might parse some meta data here
    buffer := new(bytes.Buffer)
    err := views.Scene(game).
    Render(context.Background(),buffer)
    if err != nil {
      log.Printf("Error rendering plane")
    }
    game.AddPlayer(gameobjects.NewPlayer(gameobjects.Player{
      ID: client.ID,
      X: 0,
      Y: 14,
      Z: 0,
    }))
    log.Printf("Player Connected %v", client.ID)
    client.WriteMessage(buffer.Bytes()) 
  })
  s.Listen("player_forward_up",func(client *hx.Client, msg []byte){
    game.LookupPlayer(client.ID).ControlsState.MovingForward = false
  })
  s.Listen("player_forward_down",func(client *hx.Client, msg []byte){
    game.LookupPlayer(client.ID).ControlsState.MovingForward = true
  })

  s.Listen("player_backward_up",func(client *hx.Client, msg []byte){
    game.LookupPlayer(client.ID).ControlsState.MovingBackward = false
  })
  s.Listen("player_backward_down",func(client *hx.Client, msg []byte){
    game.LookupPlayer(client.ID).ControlsState.MovingBackward = true
  })

  s.Listen("player_left_up",func(client *hx.Client, msg []byte){
    game.LookupPlayer(client.ID).ControlsState.TurningLeft = false
  })
  s.Listen("player_left_down",func(client *hx.Client, msg []byte){
    game.LookupPlayer(client.ID).ControlsState.TurningLeft = true
  })

  s.Listen("player_right_up",func(client *hx.Client, msg []byte){
    game.LookupPlayer(client.ID).ControlsState.TurningRight = false
  })
  s.Listen("player_right_down",func(client *hx.Client, msg []byte){
    game.LookupPlayer(client.ID).ControlsState.TurningRight = true
  })
}
