package gameobjects

import (
	"math"

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


func (bc *BoxCollider) FromPoints(points []utils.Vector3) BoxCollider{
	if len(points) == 0 {
		// Return an empty bounding box for no points.
		return BoxCollider{
			Position: utils.Vector3{X: 0, Y: 0, Z: 0},
			Size:     utils.Vector3{X: 0, Y: 0, Z: 0},
		}
	}

	// Initialize min and max to large/small values.
	min := utils.Vector3{X: math.Inf(1), Y: math.Inf(1), Z: math.Inf(1)}
	max := utils.Vector3{X: math.Inf(-1), Y: math.Inf(-1), Z: math.Inf(-1)}

	// Find the min and max for each axis.
	for _, p := range points {
		if p.X < min.X {
			min.X = p.X
		}
		if p.Y < min.Y {
			min.Y = p.Y
		}
		if p.Z < min.Z {
			min.Z = p.Z
		}
		if p.X > max.X {
			max.X = p.X
		}
		if p.Y > max.Y {
			max.Y = p.Y
		}
		if p.Z > max.Z {
			max.Z = p.Z
		}
	}

	// Calculate size and position.
	size := utils.Vector3{
		X: max.X - min.X,
		Y: max.Y - min.Y,
		Z: max.Z - min.Z,
	}
	position := utils.Vector3{
		X: (min.X + max.X) / 2,
		Y: (min.Y + max.Y) / 2,
		Z: (min.Z + max.Z) / 2,
	}

	return BoxCollider{
		Position: position,
		Size:     size,
	}
  
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
