package part01

import (
	"bufio"
	"hmcalister/aox22/lib"
)

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	lib.ParseFileToBrickPile(fileScanner)

	return 0, nil
}
