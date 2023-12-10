package part02

import "github.com/rs/zerolog/log"

type PartitionData struct {
	PartitionIndicesArray [][]int
}

func newPartitionData() *PartitionData {
	return &PartitionData{
		PartitionIndicesArray: make([][]int, 0),
	}
}

func (partition *PartitionData) determinePartition(partitionID int, xCoord int, yCoord int) {
	log.Debug().
		Int("NewPartitionID", partitionID).
		Int("NewPartitionStartX", xCoord).
		Int("NewPartitionStartY", yCoord).
		Send()

	nodesToExpand := []NodeData{{
		XCoordinate: xCoord,
		YCoordinate: yCoord,
		NodeRune:    PipeMaze[yCoord][xCoord],
	}}

	step := 0
	for {
		step += 1
		if len(nodesToExpand) == 0 {
			break
		}

		currentNode := nodesToExpand[0]
		nodesToExpand = nodesToExpand[1:]
		log.Debug().
			Int("partitionID", partitionID).
			Int("Step", step).
			Int("QueueLen", len(nodesToExpand)).
			Interface("CurrentNode", currentNode).
			Interface("Queue", nodesToExpand).
			Send()

		if partition.PartitionIndicesArray[yCoord][xCoord] != 0 {
			continue
		}

		partition.PartitionIndicesArray[yCoord][xCoord] = partitionID

		if yCoord > 0 {
			nodesToExpand = append(nodesToExpand, currentNode.nextNode(DIRECTION_NORTH))
		}
		if yCoord < len(PipeMaze)-1 {
			nodesToExpand = append(nodesToExpand, currentNode.nextNode(DIRECTION_SOUTH))
		}
		if xCoord > 0 {
			nodesToExpand = append(nodesToExpand, currentNode.nextNode(DIRECTION_WEST))
		}
		if xCoord < len(PipeMaze[yCoord])-1 {
			nodesToExpand = append(nodesToExpand, currentNode.nextNode(DIRECTION_EAST))
		}
	}
}
