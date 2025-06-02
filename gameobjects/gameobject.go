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
  CorrectionNormal utils.Vector3
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

func (a *BoxCollider) IsColliding(b *BoxCollider) bool {
	aMin := utils.Vector3{
		X: a.Position.X - a.Size.X/2,
		Y: a.Position.Y - a.Size.Y/2,
		Z: a.Position.Z - a.Size.Z/2,
	}
	aMax := utils.Vector3{
		X: a.Position.X + a.Size.X/2,
		Y: a.Position.Y + a.Size.Y/2,
		Z: a.Position.Z + a.Size.Z/2,
	}

	bMin := utils.Vector3{
		X: b.Position.X - b.Size.X/2,
		Y: b.Position.Y - b.Size.Y/2,
		Z: b.Position.Z - b.Size.Z/2,
	}
	bMax := utils.Vector3{
		X: b.Position.X + b.Size.X/2,
		Y: b.Position.Y + b.Size.Y/2,
		Z: b.Position.Z + b.Size.Z/2,
	}

  // log.Printf("amin: %+v, amax: %+v bmin: %+v , bmax: %+v", aMin,aMax,bMin, bMax);

	return aMin.X <= bMax.X && aMax.X >= bMin.X &&
		aMin.Y <= bMax.Y && aMax.Y >= bMin.Y &&
		aMin.Z <= bMax.Z && aMax.Z >= bMin.Z
}



type PhysicsObject struct {
  GameObject
  PreviousPosition utils.Vector3
  RotationalVelocity utils.Vector3
  Velocity utils.Vector3
  Acceleration utils.Vector3
  Collider BoxCollider
}
