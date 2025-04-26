package utils;

type Dimension2 struct {
  Width float64
  Height float64
}


func (d Dimension2) ToVector3() Vector3 {
  return Vector3{
    X: d.Width,
    Y: d.Height,
  } 
}
