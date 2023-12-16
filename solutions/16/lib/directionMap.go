package lib

// Create a map detailing the resulting directions for a beam.
//
// The first map key, of type DirectionEnum, is the current direction of the light ray.
//
// The second may key, of type LayoutRuneEnum, is the current rune the light ray is encountering.
//
// The resulting output, of type []DirectionEnum, is a list of all the resulting light ray directions.
// Note this is an array as it is possible (see, splitters) for one light ray to become multiple.
func CreateDirectionMap() map[DirectionEnum]map[LayoutRuneEnum][]DirectionEnum {
	NorthwardsRunes := map[LayoutRuneEnum][]DirectionEnum{
		EMPTY_RUNE:           {DIRECTION_NORTH},
		FORWARD_SLASH_MIRROR: {DIRECTION_EAST},
		BACK_SLASH_MIRROR:    {DIRECTION_WEST},
		VERTICAL_SPLITTER:    {DIRECTION_NORTH},
		HORIZONTAL_SPLITTER:  {DIRECTION_WEST, DIRECTION_EAST},
	}
	EastwardsRunes := map[LayoutRuneEnum][]DirectionEnum{
		EMPTY_RUNE:           {DIRECTION_EAST},
		FORWARD_SLASH_MIRROR: {DIRECTION_NORTH},
		BACK_SLASH_MIRROR:    {DIRECTION_SOUTH},
		VERTICAL_SPLITTER:    {DIRECTION_NORTH, DIRECTION_SOUTH},
		HORIZONTAL_SPLITTER:  {DIRECTION_EAST},
	}
	SouthwardsRunes := map[LayoutRuneEnum][]DirectionEnum{
		EMPTY_RUNE:           {DIRECTION_SOUTH},
		FORWARD_SLASH_MIRROR: {DIRECTION_WEST},
		BACK_SLASH_MIRROR:    {DIRECTION_EAST},
		VERTICAL_SPLITTER:    {DIRECTION_SOUTH},
		HORIZONTAL_SPLITTER:  {DIRECTION_WEST, DIRECTION_EAST},
	}
	WestwardsRunes := map[LayoutRuneEnum][]DirectionEnum{
		EMPTY_RUNE:           {DIRECTION_WEST},
		FORWARD_SLASH_MIRROR: {DIRECTION_SOUTH},
		BACK_SLASH_MIRROR:    {DIRECTION_NORTH},
		VERTICAL_SPLITTER:    {DIRECTION_NORTH, DIRECTION_SOUTH},
		HORIZONTAL_SPLITTER:  {DIRECTION_WEST},
	}

	DirectionMap := map[DirectionEnum]map[LayoutRuneEnum][]DirectionEnum{
		DIRECTION_NORTH: NorthwardsRunes,
		DIRECTION_EAST:  EastwardsRunes,
		DIRECTION_SOUTH: SouthwardsRunes,
		DIRECTION_WEST:  WestwardsRunes,
	}
	return DirectionMap
}
