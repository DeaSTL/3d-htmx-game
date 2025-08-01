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

func RoughlyEqual(val uint32, expectedVal uint32, tolerance uint32) bool {
	//making sure that the lower bound doesn't overflow
	lowerBound := (uint32)(max((int)(expectedVal)-(int)(tolerance), 0))
	upperBound := expectedVal + tolerance
	// fmt.Printf("lower bound: %v upper bound: %v check value: %v value: %v \n", lowerBound, upperBound, expectedVal, val)
	return val >= lowerBound && val <= upperBound
}

// rgba value between 0-255
func RGBAToTile(r uint32, g uint32, b uint32, a uint32) uint32 {

	r = (r / 65353) * 255
	g = (g / 65353) * 255
	b = (b / 65353) * 255
	a = (a / 65353) * 255
	tolerance := (uint32)(10)
  if(r > 0 || b > 0 || g > 0){
  // fmt.Printf("%d , %d , %d \n", r,g,b)
}
	//roughly black
	if RoughlyEqual(r, 255, tolerance) && RoughlyEqual(g, 255, tolerance) && RoughlyEqual(b, 255, tolerance) {
		return WALLT_INDUSTRIAL
	}
	//roughly yellow
  if RoughlyEqual(r, 255, tolerance) && RoughlyEqual(g, 255, tolerance) && RoughlyEqual(b, 0, tolerance) {
		return WALLT_BRICK
	}
	//red
  if RoughlyEqual(r, 255, tolerance) && RoughlyEqual(g, 0, tolerance) && RoughlyEqual(b, 0, tolerance) {
		return WALLT_COPPER_INDUSTRIAL
	}
	//blue
  if RoughlyEqual(r, 0, tolerance) && RoughlyEqual(g, 0, tolerance) && RoughlyEqual(b, 255, tolerance) {
		return WALLT_HTMXCON
	}
	return WALLT_VOID
}

func (m *GameMap) LoadMap(filename string) error {
	f, err := os.Open(fmt.Sprintf("./maps/%s", filename))
	//this is so that we don't get weird glitchy shit on the edges
	imagePadding := 2

	if err != nil {
		return errors.Join(errors.New("Failed to open map file"), err)
	}
	image, _, err := image.Decode(f)

	if err != nil {
		return errors.Join(errors.New("Failed to read image file for map"), err)
	}
	// * 2 is to make sure that both sides get the same amount of padding
	imageWidth := (image.Bounds().Max.X - image.Bounds().Min.X)
	imageHeight := (image.Bounds().Max.Y - image.Bounds().Min.Y)
	mapWidth := imageWidth + (imagePadding * 2)
	mapHeight := imageHeight + (imagePadding * 2)

	fmt.Printf("image size %+v %+v", imageWidth, imageHeight)
	fmt.Printf("map size %+v %+v", mapWidth, mapHeight)

	newMap := make([][]int, mapWidth)
	for x := range mapWidth {
		newMap[x] = make([]int, mapHeight)
		for y := range mapHeight {
			//init everything else to 0
			newMap[x][y] = WALLT_VOID
			imageX := x - imagePadding
			imageY := y - imagePadding
			//skip if not within the image region
			if imageX < 0 || imageY < 0 {
				continue
			}
			if imageX > imageWidth-1 || imageY > imageHeight-1 {
				continue
			}
			r, g, b, a := image.At(imageX, imageY).RGBA()
			cellt := RGBAToTile(r, g, b, a)
			//offsets the projection by the image padding
			newMap[x][y] = (int)(cellt)
		}
	}

	m.Width = mapWidth
	m.Height = mapHeight
	m.Map = newMap

	return nil
}

func (m *GameMap) calculateNeighbors(xIndex int, zIndex int) ([]int, int) {
	var neighbors = make([]int, 0)
	//assume this wall type but anything else will override it
	wallType := WALLT_INDUSTRIAL
	//check sides and add neighbor directions to neighbor list
	if xIndex < m.Width-1 && xIndex != 0 {
		if m.Map[xIndex-1][zIndex] > 0 {
			neighbors = append(neighbors, WALLD_WEST)
			wallType = m.Map[xIndex-1][zIndex]
		}
		if m.Map[xIndex+1][zIndex] > 0 {
			neighbors = append(neighbors, WALLD_EAST)
			wallType = m.Map[xIndex+1][zIndex]
		}
	}

	if zIndex < m.Height-1 && zIndex != 0 {
		if m.Map[xIndex][zIndex-1] > 0 {
			neighbors = append(neighbors, WALLD_SOUTH)
			wallType = m.Map[xIndex][zIndex-1]
		}
		if m.Map[xIndex][zIndex+1] > 0 {
			neighbors = append(neighbors, WALLD_NORTH)
			wallType = m.Map[xIndex][zIndex+1]
		}
	}

	return neighbors, wallType
}

// returns the position and rotation of the wall
func neighborLookup(wallDirection int) (utils.Vector3, float64, error) {
	switch wallDirection {
	//up to down
	case WALLD_NORTH:
		return utils.Vector3{X: 0, Y: 0, Z: 0}, 0, nil
	case WALLD_SOUTH:
		return utils.Vector3{X: 1, Y: 0, Z: -1}, 180, nil
	//left to right
	case WALLD_WEST:
		return utils.Vector3{X: 0, Y: 0, Z: -1}, 270, nil
	case WALLD_EAST:
		return utils.Vector3{X: 1, Y: 0, Z: 0}, 90, nil
	default:
		return utils.Vector3{}, 0, errors.New("void cell")
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
			neighbors, wallType := newGameMap.calculateNeighbors(x, z)

			if newGameMap.Map[x][z] != 0 {
				continue
			}
			for _, n := range neighbors {
				offset, rotation, err := neighborLookup(n)

				if err != nil {
					continue
				}
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

				newGameMap.Walls = append(newGameMap.Walls, NewWall(utils.Vector3{
					X: float64(int(offset.X)*Constants.DefaultWallWidth()) + float64(xOffset),
					Z: float64(int(offset.Z)*Constants.DefaultWallWidth()) + float64(zOffset),
				},
					utils.Vector3{
						Y: rotation,
					},
					color,
					wallType,
				))
			}
		}
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
