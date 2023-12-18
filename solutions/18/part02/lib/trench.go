package lib

type TrenchData struct {
	// The direction this trench was dug in
	DigDirection DirectionEnum
}

func newTrench(digDirection DirectionEnum) TrenchData {
	t := TrenchData{
		DigDirection: digDirection,
	}
	return t
}
