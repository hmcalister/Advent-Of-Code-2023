package lib

import "math"

type stretchData struct {
	stretchStart int
	stretchLen   int
	interior     bool
}

type rowStretchCoverData []stretchData

func newRowStretchCover() rowStretchCoverData {
	return []stretchData{{
		stretchStart: 0,
		stretchLen:   math.MaxInt,
		interior:     false,
	}}
}

// Add a new stretch to the rowStretchCover. The stretch should already be fully formed,
// with length and interior calculated. This is intended for trenches in the left and right direction.
func (rowStretchCover rowStretchCoverData) addStretch(newStretch stretchData) {
	// We first find the stretch that contain this stretch
	//
	// We know there cannot be multiple containing stretches as this would imply overlapping trenches
	// and we never walk back over a trench.

}
