package receive

import (
	"bytes"
	"server/sensor"
	"server/socket/codes"
	"testing"
)

func TestReceiver_Receive(t *testing.T) {
	t.Run("it returns message with sensor id", func(t *testing.T) {
		var sensorId byte = 10
		queue := make(chan *ReceiverOutput)
		inputBytes := make([]byte, 2000)
		inputBytes[0] = codes.SendData
		inputBytes[1] = sensor.Imu
		inputBytes[headerSize-1] = byte(50)
		expectedMessageSize := 50*sensor.AvailableReadingTypes[sensor.Imu].Bytes*sensor.AvailableReadingTypes[sensor.Imu].Count + headerSize
		reader := bytes.NewReader(inputBytes)
		receiver := NewReceiver(reader, queue, sensorId)

		go receiver.Receive()

		actualResult := <-queue
		if len(actualResult.Message) != expectedMessageSize {
			t.Errorf("Got message length %v, expected %v", len(actualResult.Message), expectedMessageSize)
		}
		if sensorId != actualResult.Id {
			t.Errorf("Expected different id: %v vs %v", actualResult.Id, sensorId)
		}
	})
}

func BenchmarkReceiver_Receive(b *testing.B) {
	var sensorId byte = 10
	queue := make(chan *ReceiverOutput, b.N)
	inputBytes := make([]byte, 2000)
	inputBytes[0] = codes.SendData
	inputBytes[1] = sensor.Imu
	inputBytes[headerSize-1] = byte(255)
	reader := bytes.NewReader(inputBytes)
	receiver := NewReceiver(reader, queue, sensorId)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		receiver.Receive()
		reader.Reset(inputBytes)
	}
}
