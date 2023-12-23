package lib

import "fmt"

type Coordinate struct {
	X int
	Y int
}

func (coord Coordinate) String() string {
	return fmt.Sprintf("(%v, %v)", coord.X, coord.Y)
}
