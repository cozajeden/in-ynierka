package receive

import (
	"bufio"
	"encoding/binary"
	"io"
	"server/monitor"
	"server/sensor"
	"time"
)

type ReceiverOutput struct {
	Message []byte
	Id      uint8
}

type Receiver interface {
	Receive() error
}

type receiver struct {
	conn             *bufio.Reader
	interpreterQueue chan *ReceiverOutput
	id               uint8
	headerBuffer     []byte
}

const headerSize = 16
const microphoneReadSize = 2
const readsQuantityByte = headerSize - 2
const typeByte = 1

func (rec *receiver) Receive() error {
	var err error
	rec.headerBuffer, err = rec.conn.Peek(headerSize)
	start := time.Now()
	if err != nil {
		return err
	}
	formatBytes := sensor.AvailableReadingTypes[rec.headerBuffer[typeByte]].Bytes
	formatCount := sensor.AvailableReadingTypes[rec.headerBuffer[typeByte]].Count
	readsQuantity := int(binary.LittleEndian.Uint16(rec.headerBuffer[readsQuantityByte:]))
	singleReadSize := formatBytes * formatCount
	messageSize := singleReadSize * readsQuantity
	packetBuffer := make([]byte, messageSize+headerSize)

	_, err = io.ReadFull(rec.conn, packetBuffer)
	if err != nil {
		return err
	}

	monitor.ReceiveBytes(len(packetBuffer))
	monitor.ExpectBytesWritten((8 + formatBytes + 1 + 1) * formatCount * readsQuantity)
	kind := monitor.IMU
	if singleReadSize == microphoneReadSize {
		kind = monitor.Mic
	}
	monitor.ExpectDataStruct(readsQuantity, kind)
	monitor.RegisterReceiveTime(time.Now().Sub(start))

	rec.interpreterQueue <- &ReceiverOutput{
		Id:      rec.id,
		Message: packetBuffer,
	}
	return nil
}

func NewReceiver(reader io.Reader, queue chan *ReceiverOutput, id uint8) Receiver {
	return &receiver{
		conn:             bufio.NewReader(reader),
		interpreterQueue: queue,
		id:               id,
		headerBuffer:     []byte{},
	}
}
