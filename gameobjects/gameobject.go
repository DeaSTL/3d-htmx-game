package gameobjects


import (
	"github.com/deastl/htmx-doom/utils"
)


type GameObject struct {
  Position utils.Vector3;
  Rotation utils.Vector3;
}


type PhysicsObject struct {
  GameObject
  PreviousPosition utils.Vector3
  RotationalVelocity utils.Vector3
  Velocity utils.Vector3
  Acceleration utils.Vector3
}
