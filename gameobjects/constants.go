package gameobjects

type constants struct{}

var Constants = constants{}

func (c constants) DefaultWallWidth() int { return 255 }
