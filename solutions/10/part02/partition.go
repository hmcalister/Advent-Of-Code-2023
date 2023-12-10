package part02

import "github.com/rs/zerolog/log"

type PartitionData struct {
	PartitionArray        [][]int
	NumPartitions         int
	PartitionSizes        []int
	PartitionsGroundCount []int
}

func newPartitionData() *PartitionData {
	return &PartitionData{
		PartitionArray:        make([][]int, 0),
		NumPartitions:         0,
		PartitionSizes:        make([]int, 0),
		PartitionsGroundCount: make([]int, 0),
	}
}

func (partition *PartitionData) determinePartition(xCoord int, yCoord int) {
	partition.NumPartitions += 1
	newPartitionIndex := partition.NumPartitions
	partition.PartitionsGroundCount = append(partition.PartitionsGroundCount, 0)
	partition.PartitionSizes = append(partition.PartitionsGroundCount, 0)

	log.Debug().
		Int("NewPartitionID", newPartitionIndex).
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
			Int("PartitionIndex", newPartitionIndex).
			Int("Step", step).
			Int("QueueLen", len(nodesToExpand)).
			Interface("CurrentNode", currentNode).
			Interface("Queue", nodesToExpand).
			Send()

		if partition.PartitionArray[yCoord][xCoord] != 0 {
			continue
		}

		partition.PartitionArray[yCoord][xCoord] = newPartitionIndex
		partition.PartitionSizes[newPartitionIndex-1] += 1
		if currentNode.NodeRune == GROUND_RUNE {
			partition.PartitionsGroundCount[newPartitionIndex-1] += 1
		}

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

		log.Debug().
			Int("PartitionIndex", newPartitionIndex).
			Int("Step", step).
			Int("QueueLen", len(nodesToExpand)).
			Interface("Queue", nodesToExpand).
			Send()

	}
}
