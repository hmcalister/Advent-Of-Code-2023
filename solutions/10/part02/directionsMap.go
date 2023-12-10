package part02

func createDirectionMap() map[directionEnum]map[rune]directionEnum {
	NorthwardsRunes := map[rune]directionEnum{
		'|': DIRECTION_NORTH,
		'7': DIRECTION_WEST,
		'F': DIRECTION_EAST,
	}
	EastwardsRunes := map[rune]directionEnum{
		'-': DIRECTION_EAST,
		'7': DIRECTION_SOUTH,
		'J': DIRECTION_NORTH,
	}
	SouthwardsRunes := map[rune]directionEnum{
		'|': DIRECTION_SOUTH,
		'L': DIRECTION_EAST,
		'J': DIRECTION_WEST,
	}
	WestwardsRunes := map[rune]directionEnum{
		'-': DIRECTION_WEST,
		'L': DIRECTION_NORTH,
		'F': DIRECTION_SOUTH,
	}

	DirectionMap := map[directionEnum]map[rune]directionEnum{
		DIRECTION_NORTH: NorthwardsRunes,
		DIRECTION_EAST:  EastwardsRunes,
		DIRECTION_SOUTH: SouthwardsRunes,
		DIRECTION_WEST:  WestwardsRunes,
	}
	return DirectionMap
}
