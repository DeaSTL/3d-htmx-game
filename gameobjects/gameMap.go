package gameobjects

import (
  "math/rand"
	"log"
)

type GameMap struct {
	Walls      []Wall
	TranslateX int
	TranslateY int
	Rotation   int
	Players    map[string]*Player
}

func NewGameMap(options GameMap) GameMap {
	newGameMap := options

	newGameMap.Walls = []Wall{}

  for i := 0; i < 200; i++ {
    randX := (rand.Int() % 4000) - 2000
    randY := (rand.Int() % 4000) - 2000
    randAngleIndex := rand.Int() % 8
    randAngle := 45 * randAngleIndex;
    newGameMap.Walls = append(newGameMap.Walls,NewWall(Wall{
      Height: 200,
      Rotation: randAngle,
      X: randX + 5000,
      Y: randY + 5000,
    }))

  }

	newGameMap.Players = map[string]*Player{}

	log.Printf("%+v", newGameMap.Walls)

	return newGameMap
}

func (game *GameMap) AddPlayer(player Player) {
	game.Players[player.ID] = &player
}

func (game *GameMap) LookupPlayer(playerID string) *Player {
	return game.Players[playerID]
}
