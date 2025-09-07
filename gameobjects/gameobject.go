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

// GetPenetrationDepth calculates how deep two colliders are intersecting
func (a *BoxCollider) GetPenetrationDepth(b *BoxCollider) utils.Vector3 {
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

	// Calculate overlap on each axis
	overlapX := math.Min(aMax.X, bMax.X) - math.Max(aMin.X, bMin.X)
	overlapY := math.Min(aMax.Y, bMax.Y) - math.Max(aMin.Y, bMin.Y)
	overlapZ := math.Min(aMax.Z, bMax.Z) - math.Max(aMin.Z, bMin.Z)

	return utils.Vector3{X: overlapX, Y: overlapY, Z: overlapZ}
}

// GetSeparationVector returns the minimum translation vector to separate colliders
func (a *BoxCollider) GetSeparationVector(b *BoxCollider) utils.Vector3 {
	penetration := a.GetPenetrationDepth(b)
	
	// Find the axis with minimum penetration
	minPenetration := math.Abs(penetration.X)
	separationAxis := utils.Vector3{X: 1, Y: 0, Z: 0}
	
	if math.Abs(penetration.Y) < minPenetration {
		minPenetration = math.Abs(penetration.Y)
		separationAxis = utils.Vector3{X: 0, Y: 1, Z: 0}
	}
	
	if math.Abs(penetration.Z) < minPenetration {
		minPenetration = math.Abs(penetration.Z)
		separationAxis = utils.Vector3{X: 0, Y: 0, Z: 1}
	}
	
	// Determine direction based on relative positions
	direction := a.Position.Sub(b.Position)
	
	// Apply the separation along the minimum axis
	if separationAxis.X != 0 {
		if direction.X < 0 {
			return utils.Vector3{X: -penetration.X, Y: 0, Z: 0}
		}
		return utils.Vector3{X: penetration.X, Y: 0, Z: 0}
	} else if separationAxis.Y != 0 {
		if direction.Y < 0 {
			return utils.Vector3{X: 0, Y: -penetration.Y, Z: 0}
		}
		return utils.Vector3{X: 0, Y: penetration.Y, Z: 0}
	} else {
		if direction.Z < 0 {
			return utils.Vector3{X: 0, Y: 0, Z: -penetration.Z}
		}
		return utils.Vector3{X: 0, Y: 0, Z: penetration.Z}
	}
}

// SweptAABB performs continuous collision detection
func (a *BoxCollider) SweptAABB(b *BoxCollider, velocity utils.Vector3) (bool, float64, utils.Vector3) {
	// Exit early if no movement
	if velocity.X == 0 && velocity.Y == 0 && velocity.Z == 0 {
		return a.IsColliding(b), 0, utils.Vector3{}
	}

	// Expand target by source dimensions
	expandedB := BoxCollider{
		Position: b.Position,
		Size:     b.Size.Add(a.Size),
	}

	// Ray from center of a to center + velocity
	rayOrigin := a.Position
	
	// Calculate entry and exit times for each axis
	var entryTime, exitTime utils.Vector3
	
	for axis := 0; axis < 3; axis++ {
		var vel, orig, bPos, bSize float64
		switch axis {
		case 0:
			vel = velocity.X
			orig = rayOrigin.X
			bPos = expandedB.Position.X
			bSize = expandedB.Size.X
		case 1:
			vel = velocity.Y
			orig = rayOrigin.Y
			bPos = expandedB.Position.Y
			bSize = expandedB.Size.Y
		case 2:
			vel = velocity.Z
			orig = rayOrigin.Z
			bPos = expandedB.Position.Z
			bSize = expandedB.Size.Z
		}
		
		if vel == 0 {
			if orig < bPos-bSize/2 || orig > bPos+bSize/2 {
				return false, 1.0, utils.Vector3{}
			}
			switch axis {
			case 0:
				entryTime.X = math.Inf(-1)
				exitTime.X = math.Inf(1)
			case 1:
				entryTime.Y = math.Inf(-1)
				exitTime.Y = math.Inf(1)
			case 2:
				entryTime.Z = math.Inf(-1)
				exitTime.Z = math.Inf(1)
			}
		} else {
			t1 := (bPos - bSize/2 - orig) / vel
			t2 := (bPos + bSize/2 - orig) / vel
			
			switch axis {
			case 0:
				entryTime.X = math.Min(t1, t2)
				exitTime.X = math.Max(t1, t2)
			case 1:
				entryTime.Y = math.Min(t1, t2)
				exitTime.Y = math.Max(t1, t2)
			case 2:
				entryTime.Z = math.Min(t1, t2)
				exitTime.Z = math.Max(t1, t2)
			}
		}
	}
	
	// Find the latest entry and earliest exit
	entry := math.Max(math.Max(entryTime.X, entryTime.Y), entryTime.Z)
	exit := math.Min(math.Min(exitTime.X, exitTime.Y), exitTime.Z)
	
	// No collision if entry > exit or entry > 1 or entry < 0
	if entry > exit || entry > 1.0 || entry < 0 {
		return false, 1.0, utils.Vector3{}
	}
	
	// Determine collision normal
	normal := utils.Vector3{}
	if entryTime.X > entryTime.Y && entryTime.X > entryTime.Z {
		if velocity.X < 0 {
			normal.X = 1
		} else {
			normal.X = -1
		}
	} else if entryTime.Y > entryTime.Z {
		if velocity.Y < 0 {
			normal.Y = 1
		} else {
			normal.Y = -1
		}
	} else {
		if velocity.Z < 0 {
			normal.Z = 1
		} else {
			normal.Z = -1
		}
	}
	
	return true, entry, normal
}



type PhysicsObject struct {
  GameObject
  PreviousPosition utils.Vector3
  RotationalVelocity utils.Vector3
  Velocity utils.Vector3
  Acceleration utils.Vector3
  Collider BoxCollider
}
