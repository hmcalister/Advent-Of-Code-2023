package part02

import (
	"bufio"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
)

type mappingData struct {
	rangeLength      int
	sourceRangeStart int
	destRangeStart   int
}

func parseLineToMappingData(line string) *mappingData {
	var err error
	parts := strings.Fields(line)
	destinationRangeStart, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Fatal().Msgf("error parsing destinationRange:%v", err)
	}

	sourceRangeStart, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatal().Msgf("error parsing sourceRangeStart:%v", err)
	}

	rangeLength, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Fatal().Msgf("error parsing rangeLen:%v", err)
	}

	md := &mappingData{
		rangeLength,
		sourceRangeStart,
		destinationRangeStart,
	}

	log.Trace().
		Str("InitString", line).
		Int("NumFields", len(parts)).
		Int("rangeLength", rangeLength).
		Int("sourceRangeStart", sourceRangeStart).
		Int("destinationRangeStart", destinationRangeStart).
		Send()
	return md
}

func (md *mappingData) IsInRange(testValue int) bool {
	return md.sourceRangeStart <= testValue && testValue < md.sourceRangeStart+md.rangeLength
}

func (md *mappingData) MapValue(initialValue int) int {
	return initialValue + (md.destRangeStart - md.sourceRangeStart)
}

type domainMapper struct {
	domainMappingID string
	MapDataArray    []*mappingData
}

func (dm *domainMapper) MapValue(val int) int {
	for mdID, md := range dm.MapDataArray {
		if md.IsInRange(val) {
			log.Trace().Int("ApplyingMapID", mdID).Send()
			return md.MapValue(val)
		}
	}
	return val
}

type endToEndMapper struct {
	domainMappers []*domainMapper
}

func (e2em *endToEndMapper) MapValue(val int) int {
	log.Trace().Int("InitValue", val).Send()
	for _, dm := range e2em.domainMappers {
		val = dm.MapValue(val)
		log.Trace().
			Str("ApplyingDomainID", dm.domainMappingID).
			Int("MappedValue", val).
			Send()
	}
	return val
}

func parseFileToEndToEndMapper(fileScanner *bufio.Scanner) *endToEndMapper {
	consecutiveMaps := make([]*domainMapper, 0)
	// Parse each mapping in turn
	for fileScanner.Scan() {
		// Get the initial line
		domainMappingID := fileScanner.Text()
		log.Debug().Str("MappingID", domainMappingID).Send()

		currentMapData := make([]*mappingData, 0)
		for fileScanner.Scan() {
			line := fileScanner.Text()
			if len(line) == 0 {
				break
			}
			log.Trace().Str("NextLine", line).Send()
			currentMapData = append(currentMapData, parseLineToMappingData(line))
		}
		consecutiveMaps = append(consecutiveMaps, &domainMapper{domainMappingID: domainMappingID, MapDataArray: currentMapData})
	}
	return &endToEndMapper{domainMappers: consecutiveMaps}

}

type seedRangeData struct {
	SeedRangeStart  int
	SeedRangeLength int
}

// Given a range of seed to check over, return the minimum final mapping for all seeds in that range
func checkSeedRangeWorker(WorkerID int, dataChan chan *seedRangeData, resultChan chan int, fullMapper *endToEndMapper, progressBar *progressbar.ProgressBar) {
	for currentRangeData := range dataChan {
		log.Debug().Msgf("WorkerID %v starting on seed range (%v, %v)", WorkerID, currentRangeData.SeedRangeStart, currentRangeData.SeedRangeLength)
		minSeedVal := math.MaxInt

		for seed := currentRangeData.SeedRangeStart; seed < currentRangeData.SeedRangeStart+currentRangeData.SeedRangeLength; seed += 1 {
			seedMappedValue := fullMapper.MapValue(seed)
			if seedMappedValue < minSeedVal {
				minSeedVal = seedMappedValue
			}
			progressBar.Add(1)
		}
		log.Debug().Msgf("WorkerID %v finished, found min value %v", WorkerID, minSeedVal)
		resultChan <- minSeedVal
	}
}

// Given some range of seeds to check over, break the range into chunks and give to workers to check
func checkSeedRange(seedRange *seedRangeData, fullMapper *endToEndMapper) int {
	WORKER_RANGE_LENGTH := 100000
	NUM_WORKERS := 12

	minMappedValue := math.MaxInt
	var accumulatorWaitGroup sync.WaitGroup
	var workerWaitGroup sync.WaitGroup
	dataChan := make(chan *seedRangeData)
	resultsChan := make(chan int)

	pbar := progressbar.Default(int64(seedRange.SeedRangeLength), fmt.Sprintf("Seed Range Progress(%v)", seedRange.SeedRangeStart))

	accumulatorWaitGroup.Add(1)
	go func() {
		for resultMappedValue := range resultsChan {
			if resultMappedValue < minMappedValue {
				log.Debug().Msgf("New lowest value found: %v", resultMappedValue)
				minMappedValue = resultMappedValue
			}
		}
		accumulatorWaitGroup.Done()
	}()

	for workerID := 0; workerID < NUM_WORKERS; workerID += 1 {
		workerID := workerID
		workerWaitGroup.Add(1)
		go func() {
			checkSeedRangeWorker(workerID, dataChan, resultsChan, fullMapper, pbar)
			workerWaitGroup.Done()
		}()
	}

	for targetRangeStart := 0; targetRangeStart < seedRange.SeedRangeLength; targetRangeStart += WORKER_RANGE_LENGTH {
		targetRangeLength := min(WORKER_RANGE_LENGTH, seedRange.SeedRangeStart+seedRange.SeedRangeLength-targetRangeStart)
		dataChan <- &seedRangeData{SeedRangeStart: targetRangeStart, SeedRangeLength: targetRangeLength}
	}
	close(dataChan)
	workerWaitGroup.Wait()

	close(resultsChan)
	accumulatorWaitGroup.Wait()
	return minMappedValue
}

func ProcessInput(fileScanner *bufio.Scanner) (int, error) {
	// Handle seeds
	fileScanner.Scan()
	seedLine := fileScanner.Text()
	seedLine = seedLine[7:]
	seedValuesStrs := strings.Fields(seedLine)
	fileScanner.Scan()

	fullMapper := parseFileToEndToEndMapper(fileScanner)

	minMappedValue := math.MaxInt
	for i := 0; i < len(seedValuesStrs); i += 2 {
		seedValueStartStr := seedValuesStrs[i]
		seedValueRangeStr := seedValuesStrs[i+1]

		seedValueStart, err := strconv.Atoi(seedValueStartStr)
		if err != nil {
			log.Fatal().Msgf("error parsing seed value start:%v", err)
		}
		seedValueRange, err := strconv.Atoi(seedValueStartStr)
		if err != nil {
			log.Fatal().Msgf("error parsing seed value start:%v", err)
		}
		log.Debug().
			Str("SeedValueStartStr", seedValueStartStr).
			Str("SeedValueRangeStr", seedValueRangeStr).
			Int("SeedValueStart", seedValueStart).
			Int("SeedValueRange", seedValueRange).
			Send()

		rangeValue := checkSeedRange(&seedRangeData{
			SeedRangeStart:  seedValueStart,
			SeedRangeLength: seedValueRange,
		}, fullMapper)
		if rangeValue < minMappedValue {
			minMappedValue = rangeValue
		}
	}

	return minMappedValue, nil
}
