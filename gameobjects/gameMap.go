package gameobjects

import (
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/deastl/htmx-doom/utils"
)

type GameMap struct {
	Walls      []Wall
	Map        [][]int
	TranslateX int
	TranslateY int
	Rotation   int
	Width      int
	Height     int
	Players    map[string]*Player
}

const (
	WALLD_NONE  = 0
	WALLD_EAST  = 1
	WALLD_WEST  = 2
	WALLD_SOUTH = 3
	WALLD_NORTH = 4
)

const (
	CELLT_WALL = 255
	CELLT_COIN = 254
)



func RoughlyEqual(val uint32,expectedVal uint32,tolerance uint32) bool{
  //making sure that the lower bound doesn't overflow
  lowerBound := (uint32)(max((int)(expectedVal) - (int)(tolerance),0))
  upperBound := expectedVal + tolerance
  // fmt.Printf("lower bound: %v upper bound: %v check value: %v value: %v \n", lowerBound, upperBound, expectedVal, val)
  return val > lowerBound && val < upperBound
}


//rgba value between 0-255
func RGBAToTile(r uint32, g uint32, b uint32, a uint32) uint32 {

  r = (r/65353) * 255
  g = (g/65353) * 255
  b = (b/65353) * 255
  a = (a/65353) * 255
  tolerance := (uint32)(5)
  //roughly black
  if(RoughlyEqual(r,255,tolerance) && RoughlyEqual(g,255,tolerance) && RoughlyEqual(b,255,tolerance)){
    fmt.Print("Found wall cell")
    return CELLT_WALL
  }else if(RoughlyEqual(r,255,tolerance) && RoughlyEqual(g,255,tolerance) && RoughlyEqual(b,0,tolerance)){
    return CELLT_COIN
  }
  return 0;
}

func (m *GameMap) LoadMap(filename string) error {
	f, err := os.Open(fmt.Sprintf("./maps/%s", filename))

	if err != nil {
		return errors.Join(errors.New("Failed to open map file"), err)
	}
	image, _, err := image.Decode(f)

	if err != nil {
		return errors.Join(errors.New("Failed to read image file for map"), err)
	}
	imageWidth := image.Bounds().Max.X - image.Bounds().Min.X
	imageHeight := image.Bounds().Max.Y - image.Bounds().Min.Y

	newMap := make([][]int, imageHeight)
	for x := image.Bounds().Min.X; x < image.Bounds().Max.X; x++ {
		newMap[x] = make([]int, imageWidth)
		for y := image.Bounds().Min.Y; y < image.Bounds().Max.Y; y++ {
			r, g, b, a := image.At(x, y).RGBA()
      cellt := RGBAToTile(r,g,b,a);
			newMap[x][y] = (int)(cellt);
		}
	}

	m.Width = imageWidth
	m.Height = imageHeight
	m.Map = newMap

	return nil
}

func (m *GameMap) calculateNeighbors(xIndex int, zIndex int) []int {
	var neighbors = make([]int, 0)
	//check sides
	if xIndex < m.Width-1 && xIndex != 0 {
		if m.Map[xIndex-1][zIndex] == CELLT_WALL {
			neighbors = append(neighbors, WALLD_WEST)
		}
		if m.Map[xIndex+1][zIndex] == CELLT_WALL {
			neighbors = append(neighbors, WALLD_EAST)
		}
	}

	if zIndex < m.Height-1 && zIndex != 0 {
		if m.Map[xIndex][zIndex-1] == CELLT_WALL {
			neighbors = append(neighbors, WALLD_SOUTH)
		}
		if m.Map[xIndex][zIndex+1] == CELLT_WALL {
			neighbors = append(neighbors, WALLD_NORTH)
		}
	}

	return neighbors
}

// returns the position and rotation of the wall
func neighborLookup(wallDirection int) (utils.Vector3, float64) {
	switch wallDirection {
	//up to down
	case WALLD_NORTH:
		return utils.Vector3{X: 0, Y: 0, Z: 0}, 0
	case WALLD_SOUTH:
		return utils.Vector3{X: 0, Y: 0, Z: -1}, 0
	//left to right
	case WALLD_WEST:
		return utils.Vector3{X: 0, Y: 0, Z: -1}, -90
	case WALLD_EAST:
		return utils.Vector3{X: 1, Y: 0, Z: -1}, -90
	default:
		return utils.Vector3{}, 0
	}
}

func NewGameMap(options GameMap) GameMap {
	newGameMap := options

	newGameMap.Walls = []Wall{}

	err := newGameMap.LoadMap("smile.png")

	if err != nil {
    //big oof
		log.Panicf("Error loading map :%+v", err)
	}

	for x, _ := range newGameMap.Map {
		for z, _ := range newGameMap.Map[x] {
			xOffset := x * Constants.DefaultWallWidth()
			zOffset := z * Constants.DefaultWallWidth()
			neighbors := newGameMap.calculateNeighbors(x, z)
			if newGameMap.Map[x][z] != 0 {
				continue
			}
			for _, n := range neighbors {
				offset, rotation := neighborLookup(n)
				var color = ""

				switch n {
				case WALLD_EAST:
					color = "red"
					break
				case WALLD_NORTH:
					color = "yellow"
					break
				case WALLD_WEST:
					color = "blue"
					break
				case WALLD_SOUTH:
					color = "green"
					break
				default:
					color = "orange"
					break
				}

				fmt.Printf("color %s : %d\n", color, n)
				newGameMap.Walls = append(newGameMap.Walls, NewWall(utils.Vector3{
					X: float64(int(offset.X)*Constants.DefaultWallWidth()) + float64(xOffset),
					Z: float64(int(offset.Z)*Constants.DefaultWallWidth()) + float64(zOffset),
				},
					utils.Vector3{
						Y: rotation,
					}, color))
			}
		}

		fmt.Print("\n")
	}
	// newGameMap.calculateNeighbors()
	// newGameMap.genWalls()

	newGameMap.Players = map[string]*Player{}

	// log.Printf("%+v", newGameMap.Walls)

	return newGameMap
}

func (game *GameMap) AddPlayer(player *Player) {
	game.Players[player.ID] = player
}
func (game *GameMap) LookupPlayer(playerID string) *Player {
	return game.Players[playerID]
}
