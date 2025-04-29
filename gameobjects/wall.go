package gameobjects

import (
	"math"

	"github.com/deastl/htmx-doom/utils"
)

const PLANE_WIDTH = 255
//used for identifying walls
var count = 0

const (
	WALLT_INDUSTRIAL = 255
	WALLT_BRICK      = 254
	WALLT_COPPER_INDUSTRIAL = 253
	WALLT_VOID       = 0
)


type Wall struct {
	PhysicsObject
	Height     float64
	Width      float64
	Brightness float64
	Color      string
	WallType   int
	WallImage  string
	ID         int
}

func wallTypeUrl(wallType int) string {
  if(wallType == WALLT_INDUSTRIAL){return "/public/industrial-wall.jpg"}
  if(wallType == WALLT_COPPER_INDUSTRIAL){return "/public/copper-industrial-wall.jpg"}
  if(wallType == WALLT_BRICK){return "/public/brick.jpg"}

  //this will be helpful at some point
  return "IMAGE_NOT_FOUND"
}

func NewWall(position utils.Vector3, rotation utils.Vector3, color string, wallType int) Wall {
	newWall := Wall{}
	newWall.Position = position 
	newWall.Rotation = rotation
	newWall.Color = color
	newWall.WallType = wallType

  newWall.WallImage = wallTypeUrl(newWall.WallType)
	count++
	newWall.ID = count
	if newWall.Height == 0 {
		newWall.Height = 1024
		newWall.Width = 255
		newWall.Position.Y = -1024
	}
	if newWall.Brightness == 0 {
		newWall.Brightness = 80*math.Abs(
      math.Cos((math.Pi/180)*((newWall.Rotation.Y)/2))) + 40
	}

	newWall.Collider = BoxCollider{}
	//Square
	colliderPoints := []utils.Vector3{
		utils.Vector3{X: 1, Y: 1, Z: 1},
		utils.Vector3{X: -1, Y: 1, Z: 1},
		utils.Vector3{X: -1, Y: 1, Z: -1},
		utils.Vector3{X: 1, Y: 1, Z: -1},
	}

	theta := utils.Vector3{
		X: math.Sin((math.Pi / 180.00) * (newWall.Rotation.Y)),
		Z: math.Cos((math.Pi / 180.00) * (newWall.Rotation.Y)),
		Y: 1,
	}

	// log.Printf("Theta %+v", theta)

	for i, _ := range colliderPoints {
		// Shifts the points forward one and multiplies by theta
		colliderPoints[i] = colliderPoints[i].Mult(theta.Add(utils.Vector3{0, 0, 1}))
		colliderPoints[i] = colliderPoints[i].Scale(128)
		colliderPoints[i] = colliderPoints[i].Add(newWall.Position)
	}

	newWall.Collider = newWall.Collider.FromPoints(colliderPoints)

	// newWall.Collider.Size = utils.Vector3{}

	// newWall.Collider.Position = newWall.Collider.Position.
	//     Add(utils.Vector3{X: 0.5, Z: 0.5}.Scale(255)).
	//     Sub(rotatedOffset.Scale(128))

	// newWall.Collider.Size.X = math.Abs(cornerShift.X * 2)
	// newWall.Collider.Size.Z = math.Abs(cornerShift.Z * 2)
	// newWall.Collider.Size.Y = newWall.Height

	// log.Printf(
	// 	"Collider: %+v",
	// 	newWall.Collider,
	// )

	// if newWall.Rotation == 0 {
	//   newWall.Brightness = 0
	// }
	return newWall
}
