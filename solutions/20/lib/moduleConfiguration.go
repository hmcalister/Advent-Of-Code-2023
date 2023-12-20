package lib

import (
	"bufio"
	"strings"

	"github.com/rs/zerolog/log"
)

type ModuleConfigurationData struct {
	AllModules  map[string]CommunicationModuleType
	TotalPulses map[PulseTypeEnum]int
}

func ParseFileToModuleConfiguration(fileScanner *bufio.Scanner) *ModuleConfigurationData {
	moduleConfig := &ModuleConfigurationData{
		AllModules: map[string]CommunicationModuleType{},
		TotalPulses: map[PulseTypeEnum]int{
			LOW_PULSE:  0,
			HIGH_PULSE: 0,
		},
	}

	var line string
	for fileScanner.Scan() {
		line = fileScanner.Text()

		splitLine := strings.Split(line, "->")
		nameField := strings.TrimSpace(splitLine[0])
		moduleTypeIndicatorRune := nameField[0]
		neighborsField := strings.TrimSpace(splitLine[1])
		neighbors := strings.Split(neighborsField, ", ")

		var moduleID string
		if moduleTypeIndicatorRune == 'b' {
			// We have found the broadcaster module
			moduleID = nameField
		} else {
			// The first rune is not part of the name
			moduleID = nameField[1:]
		}

		log.Debug().
			Str("RawLine", line).
			Str("ModuleID", moduleID).
			Str("ModuleTypeIndicator", string(moduleTypeIndicatorRune)).
			Interface("Neighbors", neighbors).
			Send()

		moduleBase := CommunicationModuleBase{
			ModuleID:      moduleID,
			OutputModules: neighbors,
		}

		var newModule CommunicationModuleType
		switch moduleTypeIndicatorRune {
		case 'b':
			newModule = &BroadcastModule{
				CommunicationModuleBase: moduleBase,
			}
		case '%':
			newModule = &FlipFlopModule{
				CommunicationModuleBase: moduleBase,
				State:                   false,
			}
		case '&':
			newModule = &ConjunctionModule{
				CommunicationModuleBase: moduleBase,
				ReceivedPulseMemory:     make(map[string]PulseTypeEnum),
			}
		}

		moduleConfig.AllModules[moduleID] = newModule
	}

	// Now we are sure we know of all modules, we must clean up.
	//
	// The conjunction modules need to know about their inputs.
	// Walk over each module, and if the moduleOutput is a conjunction module,
	// add the current module to the memory map
	for currentModuleID, currentModule := range moduleConfig.AllModules {
		for _, outputModuleID := range currentModule.GetModuleBase().OutputModules {
			outputModule := moduleConfig.AllModules[outputModuleID]
			switch outputModule := outputModule.(type) {
			case *ConjunctionModule:
				outputModule.ReceivedPulseMemory[currentModuleID] = LOW_PULSE
			default:
				continue
			}
		}
	}

	return moduleConfig
}

func (moduleConfig *ModuleConfigurationData) PushButton() {
	pulsesQueue := make([]pulseEvent, 0)
	pulsesQueue = append(pulsesQueue, pulseEvent{
		PulseValue: LOW_PULSE,
		SenderID:   "button",
		ReceiverID: "broadcaster",
	})

	var nextPulseEvent pulseEvent
	for len(pulsesQueue) > 0 {
		nextPulseEvent, pulsesQueue = pulsesQueue[0], pulsesQueue[1:]
		moduleConfig.TotalPulses[nextPulseEvent.PulseValue] += 1

		log.Debug().Interface("CurrentPulse", nextPulseEvent).Msg("NextPulseEvent")

		targetModule, ok := moduleConfig.AllModules[nextPulseEvent.ReceiverID]
		if !ok {
			log.Trace().Msgf("failed to find any module with name %v, continuing", nextPulseEvent.ReceiverID)
			continue
		}
		moduleResponse := targetModule.ReceivePulse(nextPulseEvent.SenderID, nextPulseEvent.PulseValue)
		if moduleResponse == NO_PULSE {
			log.Trace().Msg("no pulse received from event")
			continue
		}

		for _, neighborID := range targetModule.GetModuleBase().OutputModules {
			newPulseEvent := pulseEvent{
				PulseValue: moduleResponse,
				SenderID:   targetModule.GetModuleBase().ModuleID,
				ReceiverID: neighborID,
			}

			pulsesQueue = append(pulsesQueue, newPulseEvent)

			log.Trace().
				Interface("GeneratedPulseEvent", newPulseEvent).
				Msg("PulseEventGenerated")
		}
	}
}

func (moduleConfig *ModuleConfigurationData) FindLowestButtonPushesToAchieve_RX_LOW() int {
	// The strategy here is to detect all of the cycles in the module config
	//
	// The puzzleInput is set up such that the broadcast module sends a signal to 4 other modules
	// then those four modules constitute the start of a separate cycle each. The module RX only
	// gets a LOW input when all of these cycles line up, so taking the LCM of the cycle lengths should be enough.

	var rxInput string
rxInputSearchLoop:
	for currentModuleID, currentModule := range moduleConfig.AllModules {
		for _, outputModuleID := range currentModule.GetModuleBase().OutputModules {
			if outputModuleID == "rx" {
				rxInput = currentModuleID
				break rxInputSearchLoop
			}
		}
	}

	cycleEnds := make([]string, 0)
	for currentModuleID, currentModule := range moduleConfig.AllModules {
		for _, outputModuleID := range currentModule.GetModuleBase().OutputModules {
			if outputModuleID == rxInput {
				cycleEnds = append(cycleEnds, currentModuleID)
				break
			}
		}
	}

	log.Debug().Interface("CycleEnds", cycleEnds).Send()
	cycleEndToLengthMap := make(map[string]int)
	for _, cycleEndID := range cycleEnds {
		cycleEndToLengthMap[cycleEndID] = 0
	}
	cycleLengthsDetected := 0

	var numButtonPushes int
	numButtonPushes = 0
	for cycleLengthsDetected < len(cycleEnds) {
		numButtonPushes += 1

		pulsesQueue := make([]pulseEvent, 0)
		pulsesQueue = append(pulsesQueue, pulseEvent{
			PulseValue: LOW_PULSE,
			SenderID:   "button",
			ReceiverID: "broadcaster",
		})

		var nextPulseEvent pulseEvent
		for len(pulsesQueue) > 0 {
			nextPulseEvent, pulsesQueue = pulsesQueue[0], pulsesQueue[1:]

			// If we have an emission from a cycle end
			if currentValue, ok := cycleEndToLengthMap[nextPulseEvent.SenderID]; ok {
				// If we have both not registered this cycle end AND the pulse is high, we are good to go!
				if currentValue == 0 && nextPulseEvent.PulseValue == HIGH_PULSE {
					cycleEndToLengthMap[nextPulseEvent.SenderID] = numButtonPushes
					log.Debug().
						Int("CycleLength", numButtonPushes).
						Str("CycleEndID", nextPulseEvent.SenderID).
						Interface("CycleLengthsMap", cycleEndToLengthMap).
						Send()
					cycleLengthsDetected += 1
				}
			}

			log.Trace().Interface("CurrentPulse", nextPulseEvent).Msg("NextPulseEvent")

			targetModule, ok := moduleConfig.AllModules[nextPulseEvent.ReceiverID]
			if !ok {
				log.Trace().Msgf("failed to find any module with name %v, continuing", nextPulseEvent.ReceiverID)
				continue
			}
			moduleResponse := targetModule.ReceivePulse(nextPulseEvent.SenderID, nextPulseEvent.PulseValue)
			if moduleResponse == NO_PULSE {
				log.Trace().Msg("no pulse received from event")
				continue
			}

			for _, neighborID := range targetModule.GetModuleBase().OutputModules {
				newPulseEvent := pulseEvent{
					PulseValue: moduleResponse,
					SenderID:   targetModule.GetModuleBase().ModuleID,
					ReceiverID: neighborID,
				}

				pulsesQueue = append(pulsesQueue, newPulseEvent)

				log.Trace().
					Interface("GeneratedPulseEvent", newPulseEvent).
					Msg("PulseEventGenerated")
			}
		}
	}

	// Now, cycleLengths has the lengths of each cycle, we just have to find the LCM of this

	cycleLengths := make([]int, 0)
	for _, cycleLength := range cycleEndToLengthMap {
		cycleLengths = append(cycleLengths, cycleLength)
	}

	return lcmOfSlice(cycleLengths)
}
