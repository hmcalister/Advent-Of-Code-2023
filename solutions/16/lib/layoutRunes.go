package lib

//go:generate stringer -type=LayoutRuneEnum
type LayoutRuneEnum rune

const (
	EMPTY_RUNE           LayoutRuneEnum = '.'
	FORWARD_SLASH_MIRROR LayoutRuneEnum = '/'
	BACK_SLASH_MIRROR    LayoutRuneEnum = '\\'
	VERTICAL_SPLITTER    LayoutRuneEnum = '|'
	HORIZONTAL_SPLITTER  LayoutRuneEnum = '-'

	priv_NON_ENERGIZED_RUNE LayoutRuneEnum = '.'
	priv_ENERGIZED_RUNE     LayoutRuneEnum = '#'
)
