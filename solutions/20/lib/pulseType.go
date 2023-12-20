package lib

//go:generate stringer -type PulseTypeEnum
type PulseTypeEnum int

const (
	NO_PULSE   PulseTypeEnum = iota
	LOW_PULSE  PulseTypeEnum = iota
	HIGH_PULSE PulseTypeEnum = iota
)

type pulseEvent struct {
	PulseValue PulseTypeEnum
	SenderID   string
	ReceiverID string
}
