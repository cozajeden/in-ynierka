package monitor

import "time"

type target uint8

const (
	bytesReceivedTarget = target(iota)
	bytesSentTarget
	bytesWrittenTarget
	expectedBytesWrittenTarget
	bufferedBytesTarget

	processQueueLengthTarget
	interpreterQueueLengthTarget

	dataStructuresExpectedTarget
	dataStructuresCreatedTarget
	dataStructuresSentTarget
	dataStructuresReceivedTarget

	receiveTimeTarget
	interpretingTimeTarget
	fileSavingTimeTarget
	parseTimeTarget
)

type payload struct {
	target target
	value  interface{}
}

type dataStructurePayload struct {
	kind     DSKind
	quantity int
}

var rawQuantities = map[target]int{
	bytesReceivedTarget:          0,
	bytesSentTarget:              0,
	bytesWrittenTarget:           0,
	expectedBytesWrittenTarget:   0,
	processQueueLengthTarget:     0,
	interpreterQueueLengthTarget: 0,
	bufferedBytesTarget:          0,
}

var dataStructureQuantities = map[target]map[DSKind]int{
	dataStructuresExpectedTarget: make(map[DSKind]int),
	dataStructuresCreatedTarget:  make(map[DSKind]int),
	dataStructuresSentTarget:     make(map[DSKind]int),
	dataStructuresReceivedTarget: make(map[DSKind]int),
}

var times = map[target][]time.Duration{
	receiveTimeTarget:      make([]time.Duration, 0, 1000),
	interpretingTimeTarget: make([]time.Duration, 0, 1000),
	fileSavingTimeTarget:   make([]time.Duration, 0, 1000),
	parseTimeTarget:        make([]time.Duration, 0, 1000),
}
