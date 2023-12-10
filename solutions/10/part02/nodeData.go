package part02

import "strings"

type NodeData struct {
	XCoordinate int
	YCoordinate int
	NodeRune    rune
}

func (node NodeData) nextNode(direction directionEnum) NodeData {
	var nextXCoord int
	var nextYCoord int
	switch direction {
	case DIRECTION_NORTH:
		nextXCoord, nextYCoord = node.XCoordinate, node.YCoordinate-1
	case DIRECTION_EAST:
		nextXCoord, nextYCoord = node.XCoordinate+1, node.YCoordinate
	case DIRECTION_SOUTH:
		nextXCoord, nextYCoord = node.XCoordinate, node.YCoordinate+1
	case DIRECTION_WEST:
		nextXCoord, nextYCoord = node.XCoordinate-1, node.YCoordinate
	}

	nextNode := NodeData{
		XCoordinate: nextXCoord,
		YCoordinate: nextYCoord,
		NodeRune:    PipeMaze[nextYCoord][nextXCoord],
	}

	return nextNode
}

func determineStartDirection(startNode NodeData) directionEnum {
	nextNode := startNode.nextNode(DIRECTION_NORTH)
	if strings.ContainsRune("|7F", nextNode.NodeRune) {
		return DIRECTION_NORTH
	}

	return DIRECTION_EAST
}

type LoopData struct {
	StartXCoordinate int
	StartYCoordinate int
	LoopNodes        []NodeData
	LoopDirection    []directionEnum
}
