package gameobjects

import (
	"math"
	"sync"

	"github.com/deastl/htmx-doom/utils"
	hx "github.com/deastl/hxsocketsfiber"
)
const (
  floorLevel = 300
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
func getDirection(degrees float64) string {
	// Normalize the degrees to the range [0, 360)
	degrees = float64(int(degrees) % 360)
	if degrees < 0 {
		degrees += 360
	}

	if degrees >= 337.5 || degrees < 22.5 {
		return "N"
	} else if degrees >= 22.5 && degrees < 67.5 {
		return "NE"
	} else if degrees >= 67.5 && degrees < 112.5 {
		return "E"
	} else if degrees >= 112.5 && degrees < 157.5 {
		return "SE"
	} else if degrees >= 157.5 && degrees < 202.5 {
		return "S"
	} else if degrees >= 202.5 && degrees < 247.5 {
		return "SW"
	} else if degrees >= 247.5 && degrees < 292.5 {
		return "W"
	} else if degrees >= 292.5 && degrees < 337.5 {
		return "N"
	}
	return "Unknown"
}
func NewPlayer(options *Player) *Player {
  options.Lock()
	newPlayer := options
	newPlayer.updateDirection()
	newPlayer.MovementSpeed = 30
	newPlayer.Position.Y = floorLevel
  newPlayer.Collider = BoxCollider{
    Size: utils.Vector3{128,300,128},
    Position: newPlayer.Position,
  }
  options.Unlock()
	return newPlayer
}

func (p *Player) Jump() {
	if(p.PlayerJumping){ return}
	// p.Velocity.Y = 500
	p.PlayerJumping = true
	p.Acceleration.Y = 150
}

func (p *Player) updateDirection() {
	p.Direction.Z = math.Sin(float64(p.Rotation.Y-90) * (math.Pi / 180))
	p.Direction.X = math.Cos(float64(p.Rotation.Y-90) * (math.Pi / 180))
}

func (p *Player) CalaculateCollision(gameMap *GameMap){
  p.closeToWall = false;
  // lowestDistanceX, lowestDistanceY := 0,0;
  collisions := 0 
  for _, wall := range gameMap.Walls {
    if(wall.Collider.IsColliding(&p.Collider)){
      collisions++
      // diff := p.Position.Sub(p.PreviousPosition)
      // p.Velocity = p.Velocity.Add(diff.Scale(-2))
      p.Position = p.PreviousPosition
      p.closeToWall = true
    }
    // XDist := math.Abs(wall.Position.X - p.Position.X)
    // YDist := math.Abs(wall.Position.Z - p.Position.Z)
    // if XDist < 100  && YDist < 100{
    // }
  }
  // log.Printf("Collissions %+v", collisions)
}
func (p *Player) Update() {
  p.Stats = []Stat{}
	if p.ControlsState.Space {
		p.Jump()
	}

	if p.ControlsState.TurningRight {
		p.RotationalVelocity.Y += 4
    // p.CameraState.Rotation.Y -= 4
	}
	if p.ControlsState.TurningLeft {
		p.RotationalVelocity.Y -= 4
    // p.CameraState.Rotation.Y += 4
	}


  p.PreviousPosition = p.Position


	if p.ControlsState.MovingBackward {
    p.Velocity = p.Velocity.Sub(p.Direction.Scale(p.MovementSpeed))
	}

	if p.ControlsState.MovingForward {
    p.Velocity = p.Velocity.Add(p.Direction.Scale(p.MovementSpeed))
	}


  //little bit of a jump boost when in the air
  if(p.PlayerJumping){
    //apply gravity
    // p.Velocity.Y += 15;
    p.MovementSpeed = 40
  }else{
    p.MovementSpeed = 30
  }

  p.Velocity = p.Velocity.Add(
    p.Acceleration,
  ) 
  p.Position = p.Position.Add(
    p.Velocity,
  )

  p.Rotation = p.Rotation.Add(
    p.RotationalVelocity,
  )
  p.CameraState.Rotation = p.CameraState.Rotation.Sub(
    p.RotationalVelocity,
  )
  
  //friction bitch
  p.Velocity = p.Velocity.Scale(0.1)
  p.RotationalVelocity = p.RotationalVelocity.Scale(0.1);

  p.updateDirection()

  //keep player from yeeting through the floor
  if(p.Position.Y < floorLevel){
    p.Position.Y = floorLevel;
    p.PlayerJumping = false
  }




  p.Velocity.Y -= 25;
  p.Acceleration = p.Acceleration.Scale(0.65)


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
  p.Stats = append(p.Stats,Stat{
    Key: "Rotation",
    Value: p.Rotation,
  })

  p.Stats = append(p.Stats,Stat{
    Key: "Compass",
    Value: getDirection(p.Rotation.Y),
  })

  p.Collider.Position = p.Position


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
