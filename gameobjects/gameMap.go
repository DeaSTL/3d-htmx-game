package gameobjects

import (
	"log"
)

const WIDTH = 16
const HEIGHT = 16

type GameMap struct {
	Walls      []Wall
	Map        [WIDTH][HEIGHT]int
	TranslateX int
	TranslateY int
	Rotation   int
	Width      int
	Height     int
	Players    map[string]*Player
}

func NewGameMap(options GameMap) GameMap {
	newGameMap := options

	newGameMap.Walls = []Wall{}

	newGameMap.Map = [WIDTH][HEIGHT]int{
		{000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, },
		{000, 255, 255, 255, 255, 255, 255, 000, 255, 255, 255, 255, 255, 255, 255, },
		{000, 255, 255, 255, 255, 255, 255, 000, 255, 255, 255, 255, 255, 255, 255, },
		{000, 255, 255, 255, 255, 255, 255, 000, 255, 255, 255, 255, 255, 255, 255, },
		{000, 255, 255, 255, 255, 255, 255, 000, 255, 255, 255, 255, 255, 255, 255, },
		{000, 255, 255, 255, 255, 255, 255, 000, 255, 255, 255, 255, 255, 255, 255, },
		{000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, },
		{000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, },
		{000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, },
		{000, 255, 255, 255, 255, 255, 255, 000, 255, 255, 255, 255, 255, 255, 255, },
		{000, 255, 255, 255, 255, 255, 255, 000, 255, 255, 255, 255, 255, 255, 255, },
		{000, 255, 255, 255, 255, 255, 255, 000, 255, 255, 255, 255, 255, 255, 255, },
		{000, 255, 255, 255, 255, 255, 255, 000, 255, 255, 255, 255, 255, 255, 255, },
		{000, 255, 255, 255, 255, 255, 255, 000, 255, 255, 255, 255, 255, 255, 255, },
		{000, 255, 255, 255, 255, 255, 255, 000, 255, 255, 255, 255, 255, 255, 255, },
		{000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, 000, },
	}

	newGameMap.genWalls()

	// for i := 0; i < 200; i++ {
	// randX := (rand.Int() % 40000) - 20000
	// randY := (rand.Int() % 40000) - 20000
	// randAngleIndex := rand.Int() % 8
	// randAngle := 45 * randAngleIndex
	// newGameMap.Walls = append(newGameMap.Walls, NewWall(Wall{
	// 	Height:   rand.Int() % 500,
	// 	Rotation: randAngle,
	// 	X:        randX + 5000,
	// 	Y:        randY + 5000,
	// }))

	// }

	newGameMap.Players = map[string]*Player{}

	log.Printf("%+v", newGameMap.Walls)

	return newGameMap
}

func (game *GameMap) AddPlayer(player Player) {
	game.Players[player.ID] = &player
}

func (game *GameMap) genWalls() {
	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			nNeighbor := 255
			sNeighbor := 255
			wNeighbor := 255
			eNeighbor := 255
			if y > 0 && y < HEIGHT {
				sNeighbor = game.Map[x][y-1]
			}
			if y < HEIGHT-1 {
				nNeighbor = game.Map[x][y+1]
			}
			if x > 0{
				wNeighbor = game.Map[x-1][y]
			}

      if x < WIDTH - 1 {
				eNeighbor = game.Map[x+1][y]
      }
			currentTile := game.Map[x][y]
			const ROTATION = 0
			const PLANE_WIDTH = 255
			if currentTile == 000 {
				//000|255
				if eNeighbor == 255 {
					game.Walls = append(game.Walls, NewWall(Wall{
						X:        (x) * PLANE_WIDTH,
						Z:        (y) * PLANE_WIDTH,
						Rotation: 270 + ROTATION,
					}))
				}
				//255
				//---
				//000
				if nNeighbor == 255 {
					game.Walls = append(game.Walls, NewWall(Wall{
						X:        (x) * PLANE_WIDTH,
						Z:        (y) * PLANE_WIDTH,
						Rotation: 180 + ROTATION,
					}))
				}
				//255|000
				if wNeighbor == 255 {
					game.Walls = append(game.Walls, NewWall(Wall{
						X:        (x) * PLANE_WIDTH,
						Z:        (y) * PLANE_WIDTH,
						Rotation: 90 + ROTATION,
					}))
				}
				//000
				//---
				//255
				if sNeighbor == 255 {
					game.Walls = append(game.Walls, NewWall(Wall{
						X:        (x) * PLANE_WIDTH,
						Z:        (y) * PLANE_WIDTH,
						Rotation: 0 + ROTATION,
					}))
				}
			}
		}
	}
}

func (game *GameMap) LookupPlayer(playerID string) *Player {
	return game.Players[playerID]
}
