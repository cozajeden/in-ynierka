package translate

import (
	"encoding/binary"
	"math"
	"reflect"
	"server/sensor"
	"testing"
)

func TestTranslator_Translate(t *testing.T) {
	t.Run("It translates readings of type int16", func(t *testing.T) {
		testReadingType := sensor.ReadingType{Format: sensor.I16, Bytes: 2, Count: 2, Labels: []byte{0, 1}}
		translator := Translator{testReadingType}
		rawMessage := make([]byte, 0)
		byteBuffer := make([]byte, testReadingType.Bytes)
		firstValue := int16(-100)
		secondValue := int16(100)
		binary.LittleEndian.PutUint16(byteBuffer, uint16(firstValue))
		rawMessage = append(rawMessage, byteBuffer...)
		binary.LittleEndian.PutUint16(byteBuffer, uint16(secondValue))
		rawMessage = append(rawMessage, byteBuffer...)
		expected := map[string]string{
			"first":  "-100",
			"second": "100",
		}

		got, err := translator.Translate(rawMessage)
		if err != nil {
			t.Errorf("Translate returned error: %v", err)
		}

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Expected %v, got %v", expected, got)
		}
	})

	t.Run("It translates readings of type float32", func(t *testing.T) {
		testReadingType := sensor.ReadingType{Format: sensor.F32, Bytes: 4, Count: 2, Labels: []byte{0, 1}}
		translator := Translator{testReadingType}
		firstValue := float32(3.23)
		secondValue := float32(-101.21)
		expected := map[string]string{
			"first":  "3.23000",
			"second": "-101.21000",
		}
		rawMessage := make([]byte, 0)
		firstValueUint := math.Float32bits(firstValue)
		secondValueUint := math.Float32bits(secondValue)
		byteBuffer := make([]byte, 4)
		binary.LittleEndian.PutUint32(byteBuffer, firstValueUint)
		rawMessage = append(rawMessage, byteBuffer...)
		binary.LittleEndian.PutUint32(byteBuffer, secondValueUint)
		rawMessage = append(rawMessage, byteBuffer...)

		got, err := translator.Translate(rawMessage)
		if err != nil {
			t.Errorf("Translate returned error: %v", err)
		}

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Expected %v, got %v", expected, got)
		}
	})

	t.Run("It returns error when data size is unexpected", func(t *testing.T) {
		testReadingType := sensor.ReadingType{sensor.I16, 8, 2, []byte{}}
		translator := Translator{testReadingType}
		rawMessage := []byte{0, 0, 0}

		_, err := translator.Translate(rawMessage)

		if err == nil {
			t.Error()
		}
	})

	t.Run("It returns error when format is unknown", func(t *testing.T) {
		testReadingType := sensor.ReadingType{Format: 0xFF, Bytes: 2, Count: 2, Labels: []byte{}}
		translator := Translator{testReadingType}
		rawMessage := []byte{0, 0, 0, 0}

		_, err := translator.Translate(rawMessage)

		if err == nil {
			t.Error()
		}
	})
}
