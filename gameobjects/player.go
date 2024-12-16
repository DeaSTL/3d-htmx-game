package gameobjects

import (
	"math"
	"sync"

	"github.com/deastl/htmx-doom/utils"
	hx "github.com/deastl/hxsocketsfiber"
)

type PlayerControls struct {
	TurningLeft    bool
	TurningRight   bool
	MovingForward  bool
	MovingBackward bool
	Space          bool
}

type PlayerDirection struct {
	X float64
	Y float64
}
type Camera struct {
	GameObject
}
type Player struct {
	PhysicsObject
	FrameCount    int64
	Stats         []Stat
	Socket        *hx.Client
	ControlsState PlayerControls
	CameraState   Camera
	Direction     utils.Vector3
	PlayerJumping bool
	ID            string
	MovementSpeed float64
	sync.Mutex
}

func NewPlayer(options Player) Player {
	newPlayer := options
	newPlayer.updateDirection()
	newPlayer.MovementSpeed = 250
	return newPlayer
}

func (p *Player) Jump() {
	// if(p.PlayerJumping){ return}
	// p.YVel = 300
	// p.PlayerJumping = true
	// p.YAcc = 30
}

func (p *Player) updateDirection() {
	p.Direction.Z = math.Sin(float64(p.Rotation.Y+90) * (math.Pi / 180))
	p.Direction.X = math.Cos(float64(p.Rotation.Y+90) * (math.Pi / 180))
}

func (p *Player) Update() {
  p.Stats = []Stat{}
	if p.ControlsState.Space {
		p.Jump()
	}

	if p.ControlsState.TurningRight {
		p.Rotation.Y += 4
		p.updateDirection()
	}
	if p.ControlsState.TurningLeft {
		p.Rotation.Y -= 4
		p.updateDirection()
	}
  p.PreviousPosition = p.Position
	if p.ControlsState.MovingBackward {
		p.Position = p.Position.Sub(p.Direction.Scale(p.MovementSpeed))
	}
	if p.ControlsState.MovingForward {
		p.Position = p.Position.Add(p.Direction.Scale(p.MovementSpeed))
	}

	p.CameraState.Position = p.Position
	p.CameraState.Rotation = p.Rotation
	//Sets head level
	p.CameraState.Position.Y = 2048


  p.Stats = append(p.Stats,Stat{
    Key: "Position",
    Value: p.Position,
  })
  p.Stats = append(p.Stats,Stat{
    Key: "PrevPosition",
    Value: p.Position,
  })
  p.Stats = append(p.Stats,Stat{
    Key: "Velocity",
    Value: p.Velocity,
  })
  p.Stats = append(p.Stats,Stat{
    Key: "Acceleration",
    Value: p.Acceleration,
  })

	// p.YVel -= 10;
	//
	// p.YVel += p.YAcc
	//
	//
	// p.Y += p.YVel
	//
	// if p.Y < -4096 {
	//   p.Y = -4096
	// }

	// p.Y += p.YVel
	// if(p.YVel > 1) {
	//   p.YVel -= 1
	// }
	// if(p.Y < -4096) {
	//   p.Y = -4096
	//   p.YVel = 0;
	//   p.PlayerJumping = false
	//   log.Printf("Stopping movement")
	// }
	//
	// if(p.Y > 10){
	//   p.YVel -= 10
	//   if p.YVel < -1 {
	//     p.YVel = -200;
	//   }
	//   log.Printf("Applying gravity")
	// }

	// log.Printf("Position: %+v", p.Position)
}
