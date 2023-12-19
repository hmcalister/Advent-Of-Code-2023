package lib

//go:generate stringer --type PartPropertyEnum
type partPropertyEnum rune

const (
	ExtremelyCoolProperty partPropertyEnum = 'x'
	MusicalProperty       partPropertyEnum = 'm'
	AerodynamicProperty   partPropertyEnum = 'a'
	ShinyProperty         partPropertyEnum = 's'
)
