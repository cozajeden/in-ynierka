package monitor

import (
	"os"
	"time"
)

type DSKind uint8
type replacement int

const (
	Mic = DSKind(iota)
	IMU
)

var FilePath string
var monitoringRefreshPeriod = time.Second
var monitoringActive bool

func RunMonitoring() {
	monitoringActive = true
	go func() {
		var err error
		ticker := time.NewTicker(monitoringRefreshPeriod)
		for {
			select {
			case <-ticker.C:
				err = printStats()
			case toKeep := <-storekeeper:
				err = keep(toKeep)
			}
			if err != nil {
				panic(err)
			}
		}
	}()
}

func printStats() error {
	output, err := os.Create(FilePath)
	if err != nil {
		return err
	}
	defer output.Close()
	createStats(output)
	if err != nil {
		return err
	}
	return nil
}

func sendPayload(p *payload) {
	if monitoringActive {
		storekeeper <- p
	}
}

func ReceiveBytes(quantity int) {
	sendPayload(&payload{
		target: bytesReceivedTarget,
		value:  quantity,
	})
}

func SendBytes(quantity int) {
	sendPayload(&payload{
		target: bytesSentTarget,
		value:  quantity,
	})
}

func ExpectDataStruct(quantity int, kind DSKind) {
	sendPayload(&payload{
		target: dataStructuresExpectedTarget,
		value:  &dataStructurePayload{kind: kind, quantity: quantity},
	})
}

func ExpectBytesWritten(quantity int) {
	sendPayload(&payload{
		target: expectedBytesWrittenTarget,
		value:  quantity,
	})
}

func CreateDataStruct(quantity int, kind DSKind) {
	sendPayload(&payload{
		target: dataStructuresCreatedTarget,
		value:  &dataStructurePayload{kind: kind, quantity: quantity},
	})
}

func SendDataStruct(quantity int, kind DSKind) {
	sendPayload(&payload{
		target: dataStructuresSentTarget,
		value:  &dataStructurePayload{kind: kind, quantity: quantity},
	})
}

func ReceiveDataStruct(quantity int, kind DSKind) {
	sendPayload(&payload{
		target: dataStructuresReceivedTarget,
		value:  &dataStructurePayload{kind: kind, quantity: quantity},
	})
}

func WriteBytes(quantity int) {
	sendPayload(&payload{
		target: bytesWrittenTarget,
		value:  quantity,
	})
}

func ProcessQueueLen(length int) {
	sendPayload(&payload{
		target: processQueueLengthTarget,
		value:  replacement(length),
	})
}

func InterpreterQueueLength(length int) {
	sendPayload(&payload{
		target: interpreterQueueLengthTarget,
		value:  replacement(length),
	})
}

func RegisterReceiveTime(duration time.Duration) {
	sendPayload(&payload{
		target: receiveTimeTarget,
		value:  duration,
	})
}

func RegisterParseTime(duration time.Duration) {
	sendPayload(&payload{
		target: parseTimeTarget,
		value:  duration,
	})
}

func RegisterInterpretingTime(duration time.Duration) {
	sendPayload(&payload{
		target: interpretingTimeTarget,
		value:  duration,
	})
}

func RegisterFileSavingTime(duration time.Duration) {
	sendPayload(&payload{
		target: fileSavingTimeTarget,
		value:  duration,
	})
}

func ByteBufferLen(len int) {
	sendPayload(&payload{
		target: bufferedBytesTarget,
		value:  replacement(len),
	})
}
