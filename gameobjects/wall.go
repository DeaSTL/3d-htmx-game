package gameobjects

import (
	"log"
	"math"

	"github.com/deastl/htmx-doom/utils"
)

var count = 0

type Wall struct {
  PhysicsObject
	Height     float64
	Brightness float64
	ID         int
}

func NewWall(position utils.Vector3, rotation utils.Vector3) Wall {
	newWall := Wall{}
  newWall.Position = position
  newWall.Rotation = rotation
	count++
  newWall.ID = count;
	if newWall.Height == 0 {
		newWall.Height = 255
    newWall.Position.Y = -255
	}
	if newWall.Brightness == 0 {
		newWall.Brightness = 80 * math.Abs(math.Cos(
			(180/math.Pi)*
				newWall.Rotation.Y)) + 20
	}

  newWall.Collider.Position = newWall.Position


  rotatedOffset := utils.Vector3{
    X: math.Cos((math.Pi / 180.00) * (newWall.Rotation.Y)) * 128,
    Z: math.Sin((math.Pi / 180.00) * (newWall.Rotation.Y)) * 128,
  }
  prevPosition := newWall.Collider.Position


  
  newWall.Collider.Position = newWall.Collider.Position.
      Add(utils.Vector3{X: 128, Z: 128}).Add(rotatedOffset)

  

  log.Printf("Rotated Offset %+v angle: %+v new Position %+v", rotatedOffset, newWall.Rotation.Y,newWall.Collider.Position)
  log.Printf("prev position %+v", prevPosition)

  // if newWall.Rotation == 0 {
  //   newWall.Brightness = 0
  // }
	return newWall
}
