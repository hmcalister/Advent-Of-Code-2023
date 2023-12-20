package lib

type CommunicationModuleType interface {
	GetModuleBase() CommunicationModuleBase
	ReceivePulse(string, PulseTypeEnum) PulseTypeEnum
}

type CommunicationModuleBase struct {
	// The ID of this module
	ModuleID string

	// The ID of connected modules
	OutputModules []string
}

func (moduleBase CommunicationModuleBase) GetModuleBase() CommunicationModuleBase {
	return moduleBase
}

type BroadcastModule struct {
	CommunicationModuleBase
}

func (broadcast *BroadcastModule) ReceivePulse(senderModuleID string, incomingPulse PulseTypeEnum) PulseTypeEnum {
	return incomingPulse
}

type FlipFlopModule struct {
	CommunicationModuleBase

	// State of the flipflop module
	//
	// True means ON
	// False means OFF
	State bool
}

func (flipFlop *FlipFlopModule) ReceivePulse(senderModuleID string, incomingPulse PulseTypeEnum) PulseTypeEnum {
	if incomingPulse == LOW_PULSE {
		flipFlop.State = !flipFlop.State

		switch flipFlop.State {
		case false:
			// Module is now off, send low pulse
			return LOW_PULSE
		case true:
			return HIGH_PULSE
		}
	}

	return NO_PULSE
}

type ConjunctionModule struct {
	CommunicationModuleBase

	// The memory of the module
	// The most recent incoming pulse from EACH neighbor is remembered
	// Initialized to LOW_PULSE
	ReceivedPulseMemory map[string]PulseTypeEnum
}

func (conjunction *ConjunctionModule) ReceivePulse(senderModuleID string, incomingPulse PulseTypeEnum) PulseTypeEnum {
	conjunction.ReceivedPulseMemory[senderModuleID] = incomingPulse

	// Check if ALL memories are HIGH, and if so then return LOW
	for _, memory := range conjunction.ReceivedPulseMemory {
		// If ANY pulse is low, we return HIGH
		if memory == LOW_PULSE {
			return HIGH_PULSE
		}
	}
	return LOW_PULSE
}
