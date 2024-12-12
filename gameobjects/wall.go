package gameobjects

import "math"

var count = 0

type Wall struct {
	X          int
	Y          int
	Rotation   int
	Height     int
	Brightness int
	ID         int
}

func NewWall(options Wall) Wall {
	newWall := options
	count++
  newWall.ID = count;
	if newWall.Height == 0 {
		newWall.Height = 255
	}
	if newWall.Brightness == 0 {
		newWall.Brightness = int(80 * math.Abs(math.Cos(
			(180/math.Pi)*
				float64(newWall.Rotation)))) + 20
	}
	return newWall
}