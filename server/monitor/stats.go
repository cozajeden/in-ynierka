package monitor

import (
	"html/template"
	"os"
	"time"
)

type monitoringData struct {
	TimeMeasured        string
	SaveQueueLen        int
	InterpreterQueueLen int
	MonitoringQueueLen  int

	BytesReceived              int
	BytesReceivedSinceLastTime int
	BytesSent                  int
	BytesWritten               int
	ExpectedBytesWritten       int

	DataStructuresCreated    int
	DataStructuresCreatedMic int
	DataStructuresCreatedImu int

	ExpectedDataStructures    int
	ExpectedDataStructuresMic int
	ExpectedDataStructuresImu int

	DataStructuresSent    int
	DataStructuresSentMic int
	DataStructuresSentImu int

	DataStructuresReceived    int
	DataStructuresReceivedMic int
	DataStructuresReceivedImu int

	BufferedBytes int

	ReceiveTime   time.Duration
	InterpretTime time.Duration
	SaveTime      time.Duration
	ParseTime     time.Duration
}

var tmpl = `
Monitoring status for {{.TimeMeasured}}
======================================
Parse queue len: 	{{.SaveQueueLen}}
Interpreter queue len: 	{{.InterpreterQueueLen}}
Monitoring queue len: 	{{.MonitoringQueueLen}}

Avg time spent on receiving: 	{{.ReceiveTime}}
Avg time spent on interpreting: {{.InterpretTime}}
Avg time spent on parsing:	{{.ParseTime}}
Avg time spent on saving:	{{.SaveTime}}

Buffered bytes: 	{{.BufferedBytes}}

Bytes received: 					{{.BytesReceived}} 
	Since last measure: 				{{.BytesReceivedSinceLastTime}}
Bytes received by interpreter: 				{{.BytesSent}}
Bytes written: 						{{.BytesWritten}}
	Expected bytes written, based on receiver data: {{.ExpectedBytesWritten}}

Data structures created by interpreter: {{.DataStructuresCreated}}
	Of which:
		MIC: {{.DataStructuresCreatedMic}}
		IMU: {{.DataStructuresCreatedImu}}
Data structures expected to be created based on receiver data: {{.ExpectedDataStructures}}
	Of which:
		MIC: {{.ExpectedDataStructuresMic}}
		IMU: {{.ExpectedDataStructuresImu}}
Data structures sent by interpreter: {{.DataStructuresSent}}
	Of which:
		MIC: {{.DataStructuresSentMic}}
		IMU: {{.DataStructuresSentImu}}
Data structures received by file saver: {{.DataStructuresReceived}}
	Of which:
		MIC: {{.DataStructuresReceivedMic}}
		IMU: {{.DataStructuresReceivedImu}}
`

var lastTimeBytesReceived = 0

func createStats(output *os.File) {
	data := monitoringData{
		TimeMeasured:               time.Now().Format(time.StampMilli),
		SaveQueueLen:               rawQuantities[processQueueLengthTarget],
		InterpreterQueueLen:        rawQuantities[interpreterQueueLengthTarget],
		MonitoringQueueLen:         len(storekeeper),
		BytesReceived:              rawQuantities[bytesReceivedTarget],
		BytesReceivedSinceLastTime: rawQuantities[bytesReceivedTarget] - lastTimeBytesReceived,
		BytesSent:                  rawQuantities[bytesSentTarget],
		BytesWritten:               rawQuantities[bytesWrittenTarget],
		ExpectedBytesWritten:       rawQuantities[expectedBytesWrittenTarget],

		DataStructuresCreated:     dataStructureQuantities[dataStructuresCreatedTarget][Mic] + dataStructureQuantities[dataStructuresCreatedTarget][IMU],
		DataStructuresCreatedMic:  dataStructureQuantities[dataStructuresCreatedTarget][Mic],
		DataStructuresCreatedImu:  dataStructureQuantities[dataStructuresCreatedTarget][IMU],
		ExpectedDataStructures:    dataStructureQuantities[dataStructuresExpectedTarget][Mic] + dataStructureQuantities[dataStructuresExpectedTarget][IMU],
		ExpectedDataStructuresMic: dataStructureQuantities[dataStructuresExpectedTarget][Mic],
		ExpectedDataStructuresImu: dataStructureQuantities[dataStructuresExpectedTarget][IMU],
		DataStructuresSent:        dataStructureQuantities[dataStructuresSentTarget][Mic] + dataStructureQuantities[dataStructuresSentTarget][IMU],
		DataStructuresSentMic:     dataStructureQuantities[dataStructuresSentTarget][Mic],
		DataStructuresSentImu:     dataStructureQuantities[dataStructuresSentTarget][IMU],
		DataStructuresReceived:    dataStructureQuantities[dataStructuresReceivedTarget][Mic] + dataStructureQuantities[dataStructuresReceivedTarget][IMU],
		DataStructuresReceivedMic: dataStructureQuantities[dataStructuresReceivedTarget][Mic],
		DataStructuresReceivedImu: dataStructureQuantities[dataStructuresReceivedTarget][IMU],
		ReceiveTime:               avgDuration(times[receiveTimeTarget]),
		InterpretTime:             avgDuration(times[interpretingTimeTarget]),
		SaveTime:                  avgDuration(times[fileSavingTimeTarget]),
		ParseTime:                 avgDuration(times[parseTimeTarget]),
		BufferedBytes:             rawQuantities[bufferedBytesTarget],
	}
	lastTimeBytesReceived = rawQuantities[bytesReceivedTarget]
	parsedTemplate, err := template.New("test").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	err = parsedTemplate.Execute(output, data)
	if err != nil {
		panic(err)
	}
	clearTimes()
}

func clearTimes() {
	for i := range times {
		times[i] = make([]time.Duration, 0, 1000)
	}
}

func avgDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	sum := time.Duration(0)
	for _, d := range durations {
		sum += d
	}
	return sum / time.Duration(len(durations))
}
