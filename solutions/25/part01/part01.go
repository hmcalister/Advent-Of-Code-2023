package part01

import (
	"bufio"
	"hmcalister/aoc25/lib"
	"os"

	"github.com/dominikbraun/graph/draw"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	componentGraph := lib.ParseFileToComponentGraph(fileScanner)
	// newGraph := componentGraph.MinimumCut()

	file, _ := os.Create("./graphVis.gv")
	_ = draw.DOT(componentGraph.Graph, file)
	return 0, nil
}
