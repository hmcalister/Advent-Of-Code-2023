package lib

import "fmt"

type PropertySpaceRange struct {
	// inclusive start of the range
	RangeStart int

	// exclusive end of the range
	RangeEnd int
}

func (spaceRange PropertySpaceRange) size() int {
	return spaceRange.RangeEnd - spaceRange.RangeStart
}

func (spaceRange PropertySpaceRange) string() string {
	return fmt.Sprintf("(%v, %v)", spaceRange.RangeStart, spaceRange.RangeEnd)
}

type PartPropertySpaceRange struct {
	NextWorkflow        string
	ExtremelyCoolRating PropertySpaceRange
	MusicalRating       PropertySpaceRange
	AerodynamicRating   PropertySpaceRange
	ShinyRating         PropertySpaceRange
}

func InitialPartPropertySpaceRange() PartPropertySpaceRange {
	rangeStart := 1
	rangeEnd := 4001
	return PartPropertySpaceRange{
		NextWorkflow:        "in",
		ExtremelyCoolRating: PropertySpaceRange{rangeStart, rangeEnd},
		MusicalRating:       PropertySpaceRange{rangeStart, rangeEnd},
		AerodynamicRating:   PropertySpaceRange{rangeStart, rangeEnd},
		ShinyRating:         PropertySpaceRange{rangeStart, rangeEnd},
	}
}

func (partSpaceRange PartPropertySpaceRange) String() string {
	return fmt.Sprintf("%v, %v, %v, %v, %v",
		partSpaceRange.ExtremelyCoolRating.string(),
		partSpaceRange.MusicalRating.string(),
		partSpaceRange.AerodynamicRating.string(),
		partSpaceRange.ShinyRating.string(),
		partSpaceRange.NextWorkflow,
	)
}

func (partSpaceRange PartPropertySpaceRange) Size() int {
	return partSpaceRange.ExtremelyCoolRating.size() * partSpaceRange.MusicalRating.size() * partSpaceRange.AerodynamicRating.size() * partSpaceRange.ShinyRating.size()
}
