package lib

type TrailData struct {
	trailMap        map[Coordinate]SurfaceTypeEnum
	startCoordinate Coordinate
	endCoordiante   Coordinate
	mapWidth        int
	mapHeight       int
}
