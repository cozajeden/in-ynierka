package translate

import (
	"errors"
	"server/sensor"
)

type Translator struct {
	sensor.ReadingType
}

func (t *Translator) Translate(reading []byte) (result map[uint8][]byte, err error) {
	if len(reading) != t.Bytes*t.Count {
		err = errors.New("unexpected reading length")
		return
	}

	if t.Format == sensor.I16 {
		result = t.translateInt16(reading)
		return
	}

	if t.Format == sensor.F32 {
		result = t.translateFloat32(reading)
		return
	}

	err = errors.New("reading type is not recognized")
	return
}

func (t *Translator) translateInt16(reading []byte) map[uint8][]byte {
	result := make(map[uint8][]byte, 0)
	for i := 0; i < t.Count; i++ {
		result[t.Labels[i]] = reading[t.Bytes*i : (i+1)*t.Bytes]
	}
	return result
}

func (t *Translator) translateFloat32(reading []byte) map[uint8][]byte {
	result := make(map[uint8][]byte, 0)
	for i := 0; i < t.Count; i++ {
		result[t.Labels[i]] = reading[t.Bytes*i : (i+1)*t.Bytes]
	}
	return result
}
