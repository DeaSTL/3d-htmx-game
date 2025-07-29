package network

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/deastl/htmx-doom/gameobjects"
	"github.com/deastl/htmx-doom/utils"
	"github.com/deastl/htmx-doom/views"
	hx "github.com/deastl/hxsocketsfiber"
)

func StartRenderUpdateLoop(p *gameobjects.Player, gameMap *gameobjects.GameMap) {
	go func() {
		for {
			time.Sleep(time.Millisecond * 30)
			p.Update()

			buffer := new(bytes.Buffer)
			err := views.SceneTransform(p).
				Render(context.Background(), buffer)
			if err != nil {
				log.Printf("Error rendering plane")
			}
			// if client.ID == "NULL" {
			//   return;
			// }
			p.CalaculateCollision(gameMap)
			err = p.Socket.WriteMessage(buffer.Bytes())
			if err != nil {
				p.Socket.Close()
				break;
			}
			p.FrameCount++
			//Send stats every 3 frames
			if p.FrameCount%3 == 0 {
				buffer := new(bytes.Buffer)
				err := views.Stats(p.Stats).
					Render(context.Background(), buffer)
				if err != nil {
					log.Printf("Error rendering Stats")
				}
				err = p.Socket.WriteMessage(buffer.Bytes())

				if err != nil {
					//something bad happened
					p.Socket.Conn.Close()
					break;
				}
			}
		}
	}()
}
func RegisterPlayerMessageHandlers(s *hx.Server, game *gameobjects.GameMap) {
	s.Listen("player_init", func(client *hx.Client, msg []byte) {
		log.Printf("player_init %+v", client.ID)
		//might parse some meta data here
		if game.LookupPlayer(client.ID) != nil {
			return
		}
		buffer := new(bytes.Buffer)
		err := views.Scene(game).
			Render(context.Background(), buffer)
		if err != nil {
			log.Printf("Error rendering scene")
		}

		newPlayer := gameobjects.NewPlayer(&gameobjects.Player{
			ID:     client.ID,
			Socket: client,
		})
		newPlayer.PreviousPosition = utils.NewVector3(2024, 0, 2024)
		newPlayer.Position = utils.NewVector3(2024, 0, 2024)
		newPlayer.Rotation.Y = 128
		game.AddPlayer(newPlayer)
		log.Printf("Player Connected %v", client.ID)
		newPlayer.Socket = client

		err = newPlayer.Socket.WriteMessage(buffer.Bytes())

		if err != nil {
			newPlayer.Socket.Conn.Close()
		}

		time.Sleep(time.Second * 2)
		StartRenderUpdateLoop(newPlayer, game)
	})

	s.Listen("player_space_up", func(client *hx.Client, msg []byte) {
		game.LookupPlayer(client.ID).ControlsState.Space = false
	})
	s.Listen("player_space_down", func(client *hx.Client, msg []byte) {
		game.LookupPlayer(client.ID).ControlsState.Space = true
	})
	s.Listen("player_forward_up", func(client *hx.Client, msg []byte) {
		game.LookupPlayer(client.ID).ControlsState.MovingForward = false
	})
	s.Listen("player_forward_down", func(client *hx.Client, msg []byte) {
		game.LookupPlayer(client.ID).ControlsState.MovingForward = true
	})

	s.Listen("player_backward_up", func(client *hx.Client, msg []byte) {
		game.LookupPlayer(client.ID).ControlsState.MovingBackward = false
	})
	s.Listen("player_backward_down", func(client *hx.Client, msg []byte) {
		game.LookupPlayer(client.ID).ControlsState.MovingBackward = true
	})

	s.Listen("player_left_up", func(client *hx.Client, msg []byte) {
		game.LookupPlayer(client.ID).ControlsState.TurningLeft = false
	})
	s.Listen("player_left_down", func(client *hx.Client, msg []byte) {
		game.LookupPlayer(client.ID).ControlsState.TurningLeft = true
	})

	s.Listen("player_right_up", func(client *hx.Client, msg []byte) {
		game.LookupPlayer(client.ID).ControlsState.TurningRight = false
	})
	s.Listen("player_right_down", func(client *hx.Client, msg []byte) {
		game.LookupPlayer(client.ID).ControlsState.TurningRight = true
	})
}
