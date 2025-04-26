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
  closeToWall bool
	sync.Mutex
}

func NewPlayer(options Player) Player {
	newPlayer := options
	newPlayer.updateDirection()
	newPlayer.MovementSpeed = 25
	return newPlayer
}

func (p *Player) Jump() {
	// if(p.PlayerJumping){ return}
	// p.YVel = 300
	// p.PlayerJumping = true
	// p.YAcc = 30
}

func (p *Player) updateDirection() {
	p.Direction.Z = math.Sin(float64(p.Rotation.Y-90) * (math.Pi / 180))
	p.Direction.X = math.Cos(float64(p.Rotation.Y-90) * (math.Pi / 180))
}

func (p *Player) CalaculateCollision(gameMap *GameMap){
  p.closeToWall = false;
  // lowestDistanceX, lowestDistanceY := 0,0;
  for _, wall := range gameMap.Walls {
    XDist := math.Abs(wall.Position.X - p.Position.X)
    YDist := math.Abs(wall.Position.Z - p.Position.Z)
    if XDist < 100  && YDist < 100{
      p.Position = p.PreviousPosition
      p.closeToWall = true
    }
  }
}
func (p *Player) Update() {
  p.Stats = []Stat{}
	if p.ControlsState.Space {
		p.Jump()
	}

	if p.ControlsState.TurningRight {
		p.Rotation.Y += 4
    p.CameraState.Rotation.Y -= 4
		p.updateDirection()
	}
	if p.ControlsState.TurningLeft {
		p.Rotation.Y -= 4
    p.CameraState.Rotation.Y += 4
		p.updateDirection()
	}
  p.PreviousPosition = p.Position
	if p.ControlsState.MovingBackward {
		p.Position = p.Position.Sub(p.Direction.Scale(p.MovementSpeed))
    // p.CameraState.Position = p.CameraState.Position.Add(p.Direction.Scale(p.MovementSpeed))
	}
	if p.ControlsState.MovingForward {
		p.Position = p.Position.Add(p.Direction.Scale(p.MovementSpeed))
    // p.CameraState.Position = p.CameraState.Position.Sub(p.Direction.Scale(p.MovementSpeed))
	}

	//Sets head level
	p.Position.Y = 3000


  p.Stats = append(p.Stats,Stat{
    Key: "Position",
    Value: p.Position.Scale(1),
  })
  p.Stats = append(p.Stats,Stat{
    Key: "PrevPosition",
    Value: p.Position.Scale(1),
  })
  p.Stats = append(p.Stats,Stat{
    Key: "Velocity",
    Value: p.Velocity,
  })
  p.Stats = append(p.Stats,Stat{
    Key: "Acceleration",
    Value: p.Acceleration,
  })

  p.Stats = append(p.Stats,Stat{
    Key: "Close to wall",
    Value: p.closeToWall,
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
