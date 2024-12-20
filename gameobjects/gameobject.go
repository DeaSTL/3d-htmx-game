package gameobjects


import (
	"github.com/deastl/htmx-doom/utils"
)


type GameObject struct {
  Position utils.Vector3;
  Rotation utils.Vector3;
}


type BoxCollider struct {
  //width depth and height
  Size utils.Vector3
  Position utils.Vector3
}


func (collider *BoxCollider) IsColliding(otherCollider *BoxCollider) bool {
    // Check for overlap in the x-axis
    if collider.Position.X+collider.Size.X < otherCollider.Position.X || 
       otherCollider.Position.X+otherCollider.Size.X < collider.Position.X {
        return false
    }

    // Check for overlap in the y-axis
    if collider.Position.Y+collider.Size.Y < otherCollider.Position.Y || 
       otherCollider.Position.Y+otherCollider.Size.Y < collider.Position.Y {
        return false
    }

    // Check for overlap in the z-axis
    if collider.Position.Z+collider.Size.Z < otherCollider.Position.Z || 
       otherCollider.Position.Z+otherCollider.Size.Z < collider.Position.Z {
        return false
    }

    // If none of the axes are disjoint, the boxes are colliding
    return true
}


type PhysicsObject struct {
  GameObject
  PreviousPosition utils.Vector3
  RotationalVelocity utils.Vector3
  Velocity utils.Vector3
  Acceleration utils.Vector3
  Collider BoxCollider
}
