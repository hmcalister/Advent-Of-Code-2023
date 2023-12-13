package lib

import "github.com/rs/zerolog/log"

func findReflectionIndices(s string) []int {
	reversedString := reverseString(s)
	reflectionIndices := make([]int, 0)

	log.Trace().
		Str("TestString", s).
		Str("ReversedString", reversedString).
		Int("I_Limit", len(s)-1).
		Send()

	// We can check for a reflection by taking progressively smaller
	// substrings of the reversed string, and asking if testString ends with it
	for i := 1; i <= len(s)-1; i += 1 {
		strLen := min(i, len(s)-i)
		forwardStart := i - strLen
		forwardEnd := i
		reverseStart := len(s) - i - strLen
		reverseEnd := len(s) - i
		forwardPartial := s[forwardStart:forwardEnd]
		reversePartial := reversedString[reverseStart:reverseEnd]
		log.Trace().
			Int("TestingReflectionIndex", i).
			Str("ForwardPartial", forwardPartial).
			Str("ReversePartial", reversePartial).
			Send()
		if reversePartial == forwardPartial {
			log.Trace().Msg("Match Found")
			reflectionIndices = append(reflectionIndices, i)
		}
	}

	return reflectionIndices
}
