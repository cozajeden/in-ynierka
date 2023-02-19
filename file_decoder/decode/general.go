package decode

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"time"
)

type reading struct {
	id    uint8
	name  string
	value string
	time  time.Time
}

var names = map[uint8]string{
	0: "volume",
	1: "x-g",
	2: "y-g",
	3: "z-g",
	4: "x-dps",
	5: "y-dps",
	6: "z-dps",
}

var unwantedMeasureError = errors.New("this measure is not requested")

type Decoder struct {
	reader    io.Reader
	writer    io.Writer
	ids       map[uint8]bool
	startTime time.Time
	duration  time.Duration
}

func (d *Decoder) ParseOnlyIds(ids []int) {
	for _, id := range ids {
		d.ids[uint8(id)] = true
	}
}

func (d *Decoder) ParseSince(startTime time.Time) {
	d.startTime = startTime
}

func (d *Decoder) ParseFor(duration time.Duration) {
	d.duration = duration
}

func NewDecoder(reader io.Reader, writer io.Writer) *Decoder {
	return &Decoder{
		reader:    reader,
		writer:    writer,
		ids:       make(map[uint8]bool, 0),
		startTime: time.Time{},
		duration:  time.Duration(0),
	}
}

func (d *Decoder) parseLine() (*reading, error) {
	isThatLineToBeIgnored := false
	id, name, err := parseIdAndName(d.reader)
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, fmt.Errorf("couldn't parse id and name: %v", err)
	}
	isThatLineToBeIgnored = isThatLineToBeIgnored || (len(d.ids) > 0 && d.ids[id] != true)
	value, err := parseValue(d.reader, name)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse value: %v", err)
	}
	timeParsed, err := parseTime(d.reader)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse time: %v", err)
	}
	isThatLineToBeIgnored = isThatLineToBeIgnored ||
		(timeParsed.Before(d.startTime) || (d.duration.Nanoseconds() > 0 && timeParsed.After(d.startTime.Add(d.duration))))
	if isThatLineToBeIgnored {
		return nil, unwantedMeasureError
	}
	return &reading{
		id:    id,
		name:  name,
		value: value,
		time:  timeParsed,
	}, nil
}

func parseTime(reader io.Reader) (time.Time, error) {
	bytes := make([]byte, 8)
	_, err := reader.Read(bytes)
	if err != nil {
		return time.UnixMicro(0), fmt.Errorf("couldn't read time: %v", err)
	}
	microseconds := int64(binary.LittleEndian.Uint64(bytes))
	timeParsed := time.UnixMicro(microseconds)
	_, offset := timeParsed.Zone()
	timeUtc := timeParsed.Add(time.Duration(offset) * time.Second).UTC()
	return timeUtc, nil
}

func parseValue(reader io.Reader, name string) (string, error) {
	var binaryVal []byte
	var result string
	if name == "volume" {
		binaryVal = make([]byte, 2)
	} else {
		binaryVal = make([]byte, 4)
	}
	_, err := reader.Read(binaryVal)
	if err != nil {
		return result, fmt.Errorf("couldn't read value: %v", err)
	}
	if name == "volume" {
		value := int16(binary.LittleEndian.Uint16(binaryVal))
		result = fmt.Sprintf("%d", value)
	} else {
		u32 := binary.LittleEndian.Uint32(binaryVal)
		value := math.Float32frombits(u32)
		result = fmt.Sprintf("%.5f", value)
	}
	return result, nil
}

func parseIdAndName(reader io.Reader) (u8 uint8, s string, err error) {
	idAndName := make([]byte, 2)
	_, err = reader.Read(idAndName)
	if err != nil {
		if err == io.EOF {
			return
		}
		err = fmt.Errorf("coldn't read id and name: %v", err)
		return
	}
	return idAndName[0], names[idAndName[1]], nil
}
