package utils

import (
	"errors"
	"math"
)

type Vector3 struct {
	X, Y, Z float64
}


// NewVector3 creates a new 3D vector
func NewVector3(x, y, z float64) Vector3 {
	return Vector3{X: x, Y: y, Z: z}
}

func (v1 Vector3) Add(v2 Vector3) Vector3 {
	return Vector3{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
		Z: v1.Z + v2.Z,
	}
}

func (v1 Vector3) Zero() Vector3{
  return Vector3{0,0,0}
}
func (v1 Vector3) Sub(v2 Vector3) Vector3 {
	return Vector3{
		X: v1.X - v2.X,
		Y: v1.Y - v2.Y,
		Z: v1.Z - v2.Z,
	}
}

func (v1 Vector3) Mult(v2 Vector3) Vector3 {
	return Vector3{
		X: v1.X * v2.X,
		Y: v1.Y * v2.Y,
		Z: v1.Z * v2.Z,
	}
}

func (v1 Vector3) Dot(v2 Vector3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func (v1 Vector3) Cross(v2 Vector3) Vector3 {
	return Vector3{
		X: v1.Y*v2.Z - v1.Z*v2.Y,
		Y: v1.Z*v2.X - v1.X*v2.Z,
		Z: v1.X*v2.Y - v1.Y*v2.X,
	}
}

func (v Vector3) Norm() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector3) Normalize() (Vector3, error) {
	norm := v.Norm()
  //shit's fucked
	if norm == 0 {
    return Vector3{}, errors.New("cannot normalize a zero vector, rip!")
	}
	return Vector3{
		X: v.X / norm,
		Y: v.Y / norm,
		Z: v.Z / norm,
	}, nil
}

func (v Vector3) Scale(scalar float64) Vector3 {
	return Vector3{
		X: v.X * scalar,
		Y: v.Y * scalar,
		Z: v.Z * scalar,
	}
}

