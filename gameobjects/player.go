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
	Sprint         bool
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
	Exited bool
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
	newPlayer.MovementSpeed = 10  // Reduced from 30 for better control
	newPlayer.Position.Y = floorLevel
  newPlayer.Collider = BoxCollider{
    Size: utils.Vector3{100,300,100},  // Smaller collision box for easier navigation
    Position: newPlayer.Position,
  }
  options.Unlock()
	return newPlayer
}

func (p *Player) Jump() {
	if(p.PlayerJumping){ return}
	p.PlayerJumping = true
	// Give an immediate upward velocity for more responsive jumping
	p.Velocity.Y = 30  // Balanced jump height
	p.Acceleration.Y = 15  // Moderate boost
}

func (p *Player) updateDirection() {
	p.Direction.Z = math.Sin(float64(p.Rotation.Y-90) * (math.Pi / 180))
	p.Direction.X = math.Cos(float64(p.Rotation.Y-90) * (math.Pi / 180))
}

func (p *Player) CalaculateCollision(gameMap *GameMap){
  p.closeToWall = false;
  
  // First pass: Continuous collision detection to prevent tunneling
  var earliestCollision struct {
    time float64
    normal utils.Vector3
    wall *Wall
    hasCollision bool
  }
  earliestCollision.time = 1.0
  
  // Check swept collision for each wall
  for i := range gameMap.Walls {
    wall := &gameMap.Walls[i]
    if collides, time, normal := p.Collider.SweptAABB(&wall.Collider, p.Velocity); collides {
      if time < earliestCollision.time {
        earliestCollision.time = time
        earliestCollision.normal = normal
        earliestCollision.wall = wall
        earliestCollision.hasCollision = true
      }
    }
  }
  
  // Apply continuous collision response
  if earliestCollision.hasCollision {
    // Move to point of collision (slightly before to avoid penetration)
    safeTime := math.Max(0, earliestCollision.time - 0.01)
    p.Position = p.PreviousPosition.Add(p.Velocity.Scale(safeTime))
    
    // Calculate remaining velocity after collision
    remainingTime := 1.0 - earliestCollision.time
    remainingVelocity := p.Velocity.Scale(remainingTime)
    
    // Wall sliding: remove velocity component along collision normal
    dotProduct := remainingVelocity.Dot(earliestCollision.normal)
    slideVelocity := remainingVelocity.Sub(earliestCollision.normal.Scale(dotProduct))
    
    // Apply sliding motion
    p.Position = p.Position.Add(slideVelocity)
    
    // Dampen velocity along collision normal (bounce with energy loss)
    normalVelocity := earliestCollision.normal.Scale(p.Velocity.Dot(earliestCollision.normal))
    tangentVelocity := p.Velocity.Sub(normalVelocity)
    
    // Apply restitution (0.2 = 20% bounce, 0.8 = 80% energy loss)
    p.Velocity = tangentVelocity.Sub(normalVelocity.Scale(0.2))
    
    p.closeToWall = true
  }
  
  // Second pass: Resolve any remaining penetrations (safety check)
  for i := range gameMap.Walls {
    wall := &gameMap.Walls[i]
    if wall.Collider.IsColliding(&p.Collider) {
      // Get minimum separation vector
      separation := p.Collider.GetSeparationVector(&wall.Collider)
      
      // Push player out of wall
      p.Position = p.Position.Add(separation)
      
      // Kill velocity component toward wall
      if separation.X != 0 {
        p.Velocity.X = 0
      }
      if separation.Y != 0 {
        p.Velocity.Y = 0
      }
      if separation.Z != 0 {
        p.Velocity.Z = 0
      }
      
      p.closeToWall = true
    }
  }
}
func (p *Player) Update() {
  p.Stats = []Stat{}
  
  // Store previous position for collision rollback
  p.PreviousPosition = p.Position
  
  // Handle jumping
	if p.ControlsState.Space {
		p.Jump()
	}

  // Handle rotation
	if p.ControlsState.TurningRight {
		p.RotationalVelocity.Y += 4
	}
	if p.ControlsState.TurningLeft {
		p.RotationalVelocity.Y -= 4
	}

  // Handle movement input (add to acceleration for smoother physics)
	if p.ControlsState.MovingBackward {
    p.Acceleration = p.Acceleration.Sub(p.Direction.Scale(p.MovementSpeed * 0.5))
	}
	if p.ControlsState.MovingForward {
    p.Acceleration = p.Acceleration.Add(p.Direction.Scale(p.MovementSpeed * 0.5))
	}

  // Air control adjustment and sprint
  if p.ControlsState.Sprint && !p.PlayerJumping {
    p.MovementSpeed = 15  // Sprint speed when on ground
  } else if p.PlayerJumping {
    p.MovementSpeed = 10  // Slightly reduced air control
  } else {
    p.MovementSpeed = 8  // Normal walking speed
  }

  // Apply acceleration to velocity
  p.Velocity = p.Velocity.Add(p.Acceleration)
  
  // Apply more realistic gravity (increases when falling)
  if p.PlayerJumping {
    if p.Velocity.Y > 0 {
      // Rising - moderate gravity for smooth arc
      p.Velocity.Y -= 4.5
    } else {
      // Falling - slightly stronger gravity for natural feel
      p.Velocity.Y -= 6.5
    }
  } else {
    // On ground - standard gravity to keep grounded
    p.Velocity.Y -= 25
  }
  
  // Calculate intended position based on velocity
  p.Position = p.Position.Add(p.Velocity)

  // Update rotation
  p.Rotation = p.Rotation.Add(p.RotationalVelocity)
  p.CameraState.Rotation = p.CameraState.Rotation.Sub(p.RotationalVelocity)
  
  // Update direction based on new rotation
  p.updateDirection()
  
  // Update collider position before collision check
  p.Collider.Position = p.Position
  
  // Collision detection and response will happen in CalaculateCollision
  // which modifies position and velocity appropriately
  
  // Apply friction AFTER collision to preserve sliding
  frictionFactor := 0.15
  if p.PlayerJumping {
    frictionFactor = 0.05  // Less friction in air
  }
  p.Velocity.X *= (1.0 - frictionFactor)
  p.Velocity.Z *= (1.0 - frictionFactor)
  
  // Rotational friction
  p.RotationalVelocity = p.RotationalVelocity.Scale(0.1)
  
  // Decay acceleration (faster decay for more responsive controls)
  p.Acceleration = p.Acceleration.Scale(0.25)

  // Floor collision with small bounce absorption
  if(p.Position.Y < floorLevel){
    p.Position.Y = floorLevel
    
    // Small bounce on landing if falling fast
    if p.Velocity.Y < -50 {
      p.Velocity.Y = -p.Velocity.Y * 0.1  // 10% bounce
    } else {
      p.Velocity.Y = 0
    }
    
    p.PlayerJumping = false
  }
  
  // Terminal velocity cap (prevent infinite falling speed)
  if p.Velocity.Y < -100 {
    p.Velocity.Y = -100
  }


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
