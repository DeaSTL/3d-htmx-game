package gameobjects

import "math"

type PlayerControls struct {
	TurningLeft    bool
	TurningRight   bool
	MovingForward  bool
	MovingBackward bool
}

type PlayerDirection struct {
	X float64
	Y float64
}
type Player struct {
	ControlsState     PlayerControls
	Direction         PlayerDirection
	Rotation          float64
	X                 float64
	Y                 float64
	Z                 float64
	TranslationMatrix [4][4]float64
	RotationMatrix    [4][4]float64
	ID                string
	MovementSpeed     float64
}

func NewPlayer(options Player) Player {
	newPlayer := options
	newPlayer.updateDirection()
	newPlayer.MovementSpeed = 250
	return newPlayer
}

func (p *Player) updateDirection() {
	p.Direction.X = math.Sin(float64(p.Rotation+90) * (math.Pi / 180))
	p.Direction.Y = math.Cos(float64(p.Rotation+90) * (math.Pi / 180))
}

func (p *Player) Update() {
	if p.ControlsState.TurningRight {
		p.Rotation += 3
		p.updateDirection()
	}
	if p.ControlsState.TurningLeft {
		p.Rotation -= 3
		p.updateDirection()
	}

	if p.ControlsState.MovingBackward {
		p.Z -= p.MovementSpeed * p.Direction.X
		p.X -= p.MovementSpeed * p.Direction.Y
	}
	if p.ControlsState.MovingForward {
		p.Z += p.MovementSpeed * p.Direction.X
		p.X += p.MovementSpeed * p.Direction.Y
	}
	// p.TranslationMatrix = [4][4]float64{
	// 	{1, 0, 0, float64(p.X)},
	// 	{0, 1, 0, float64(p.Y)},
	// 	{0, 0, 1, float64(p.Z)},
	// 	{0, 0, 0, 1},
	// }
	//
	// cosTheta := math.Cos(float64(p.Rotation) * (math.Pi / 180))
	// sinTheta := math.Sin(float64(p.Rotation) * (math.Pi / 180))
	//
	// p.RotationMatrix = [4][4]float64{
	// 	{cosTheta, 0, sinTheta, 0},
	// 	{0, 1, 0, 0},
	// 	{-sinTheta, 0, cosTheta, 0},
	// 	{0, 0, 0, 1},
	// }
	// log.Printf("%+v", *p)
}
