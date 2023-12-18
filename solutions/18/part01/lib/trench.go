package lib

type ColorData struct {
	ColorString string
}

type TrenchData struct {
	// The direction this trench was dug in
	DigDirection DirectionEnum

	// The depth of the trench
	Depth int

	// The edges present in the trench (as keys) and the corresponding colors
	EdgeColors map[DirectionEnum]ColorData
}

func newTrench(digDirection DirectionEnum, depth int, edgeColor ColorData) TrenchData {
	t := TrenchData{
		DigDirection: digDirection,
		Depth:        depth,
		EdgeColors:   make(map[DirectionEnum]ColorData),
	}

	t.updateEdgeColors(digDirection, edgeColor)
	return t
}

func (trench TrenchData) updateEdgeColors(newTrenchDirection DirectionEnum, newColor ColorData) TrenchData {
	var edgeDirections []DirectionEnum
	switch newTrenchDirection {
	case DIRECTION_UP:
		fallthrough
	case DIRECTION_DOWN:
		edgeDirections = []DirectionEnum{DIRECTION_LEFT, DIRECTION_RIGHT}
	case DIRECTION_LEFT:
		fallthrough
	case DIRECTION_RIGHT:
		edgeDirections = []DirectionEnum{DIRECTION_UP, DIRECTION_DOWN}
	}

	for _, direction := range edgeDirections {
		trench.EdgeColors[direction] = newColor
	}

	return trench
}
