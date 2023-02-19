package interpreting

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"server/database"
	"server/monitor"
	"server/sensor"
	"server/socket/codes"
	"server/socket/translate"
	"time"
)

type Interpreter interface {
	Interpret() (Conclusion, error)
	Code() byte
}

type interpreter struct {
	input       *bufio.Reader
	code        byte
	readingRepo database.ReadingsRepository
	id          uint8
}

const microsecondsInSecond = 1_000_000
const nanosecondsInMicrosecond = 1000

func NewInterpreter(inputStream *bufio.Reader, rr database.ReadingsRepository, id uint8) Interpreter {
	return &interpreter{
		input:       inputStream,
		code:        codes.Ok,
		readingRepo: rr,
		id:          id,
	}
}

func (i *interpreter) Code() byte {
	code := i.code
	i.code = codes.Ok
	return code
}

func (i *interpreter) Interpret() (c Conclusion, err error) {
	_, err = i.readCommandCode()
	if err != nil {
		i.code = codes.CommandNotRecognized
		return
	}
	p, err := i.readDataTransferHeader()
	if err != nil {
		i.code = codes.UnexpectedPacketHeader
		return
	}
	readings, err := i.readDataTransferBody(p)
	if err != nil {
		i.code = codes.UnexpectedPacketBody
	}
	kind := monitor.IMU
	if len(p.rt.Labels) == 1 {
		kind = monitor.Mic
	}
	monitor.CreateDataStruct(len(readings), kind)
	c = &readConclusion{
		readings: readings,
		repo:     i.readingRepo,
		sensorId: i.id,
	}
	return
}

func (i *interpreter) readCommandCode() (code byte, err error) {
	code, err = i.input.ReadByte()
	if err != nil {
		return
	}
	if code != codes.SendData {
		err = fmt.Errorf("command code unknown: %v", code)
	}
	return
}

func (i *interpreter) readDataTransferHeader() (result *dataPackage, err error) {
	result = new(dataPackage)
	readingType, err := i.readReadingType()
	if err != nil {
		return
	}
	seconds, err := i.readUint32()
	if err != nil {
		return
	}
	microseconds, err := i.readUint32()
	if err != nil {
		return
	}
	frequency, err := i.readFloat32()
	if err != nil {
		return
	}
	readsQuantity, err := i.readUint16()
	if err != nil {
		return
	}
	result.rt = readingType
	result.sec = seconds
	result.usec = microseconds
	result.frequency = frequency
	result.quantity = readsQuantity
	return
}

func (i *interpreter) readReadingType() (rt sensor.ReadingType, err error) {
	typeId, err := i.input.ReadByte()
	if err != nil {
		return
	}
	rt, exists := sensor.AvailableReadingTypes[typeId]
	if !exists {
		err = errors.New("such reading type does not exist")
	}
	return
}

func (i *interpreter) readDataTransferBody(p *dataPackage) (readings []*database.Reading, err error) {
	translator := &translate.Translator{ReadingType: p.rt}
	readingBuffer := make([]byte, p.rt.Bytes*p.rt.Count)
	readings = make([]*database.Reading, 0, 256)
	microsecondPeriod := microsecondsInSecond / p.frequency
	for j := 0; j < int(p.quantity); j++ {
		_, err = io.ReadFull(i.input, readingBuffer)
		if err != nil {
			return
		}
		copiedBuffer := make([]byte, len(readingBuffer))
		copy(copiedBuffer, readingBuffer)

		translated, err := translator.Translate(copiedBuffer)
		if err != nil {
			return nil, err
		}

		microsecondsMeasured := p.usec + uint32(microsecondPeriod*float32(j))
		secondsMeasured := p.sec
		for microsecondsMeasured >= microsecondsInSecond {
			secondsMeasured += 1
			microsecondsMeasured -= microsecondsInSecond
		}

		timestamp := time.Unix(int64(secondsMeasured), int64(microsecondsMeasured)*nanosecondsInMicrosecond).UTC()

		readings = append(readings, &database.Reading{
			Value:     translated,
			Timestamp: timestamp,
		})
	}
	return readings, nil
}

func (i *interpreter) readUint32() (result uint32, err error) {
	buffer := make([]byte, 4)
	_, err = io.ReadFull(i.input, buffer)
	if err != nil {
		return
	}
	result = binary.LittleEndian.Uint32(buffer)
	return
}

func (i *interpreter) readFloat32() (result float32, err error) {
	buffer := make([]byte, 4)
	_, err = io.ReadFull(i.input, buffer)
	if err != nil {
		return
	}
	u32 := binary.LittleEndian.Uint32(buffer)
	result = math.Float32frombits(u32)
	return
}

func (i *interpreter) readUint16() (result uint16, err error) {
	buffer := make([]byte, 2)
	_, err = io.ReadFull(i.input, buffer)
	if err != nil {
		return
	}
	result = binary.LittleEndian.Uint16(buffer)
	return
}
