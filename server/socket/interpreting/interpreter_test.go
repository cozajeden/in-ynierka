package interpreting

import (
	"bufio"
	"bytes"
	"server/database"
	"server/socket/codes"
	"testing"
)

func TestInterpreter_Interpret(t *testing.T) {
	t.Run("it returns error on unknown command", func(t *testing.T) {
		reader := bufio.NewReader(bytes.NewReader([]byte{255}))
		rr := database.NewReadingsRepository()
		var id uint8 = 100
		interpreter := NewInterpreter(reader, rr, id)

		_, err := interpreter.Interpret()

		if err == nil {
			t.Errorf("did not return error on unknown command")
		}
	})

	t.Run("it returns error on unknown reading type", func(t *testing.T) {
		reader := bufio.NewReader(bytes.NewReader([]byte{codes.SendData, 255}))
		rr := database.NewReadingsRepository()
		var id uint8 = 100
		interpreter := NewInterpreter(reader, rr, id)

		_, err := interpreter.Interpret()

		if err == nil {
			t.Errorf("did not return error on unknown reading type")
		}
	})
}

func TestInterpreter_Code(t *testing.T) {
	t.Run("it defaults to ok code", func(t *testing.T) {
		reader := bufio.NewReader(bytes.NewReader([]byte{codes.SendData, 255}))
		rr := database.NewReadingsRepository()
		var id uint8 = 100
		interpreter := NewInterpreter(reader, rr, id)

		if interpreter.Code() != codes.Ok {
			t.Errorf("default interpreter code is not ok")
		}
	})

	t.Run("it resets to default code after use", func(t *testing.T) {
		reader := bufio.NewReader(bytes.NewReader([]byte{255}))
		rr := database.NewReadingsRepository()
		var id uint8 = 100
		interpreter := NewInterpreter(reader, rr, id)

		_, _ = interpreter.Interpret()

		if interpreter.Code() == codes.Ok {
			t.Errorf("interpreter does not return appropriate error code on unexpected command")
		}
		if interpreter.Code() != codes.Ok {
			t.Errorf("interpreter does not come back to defualt code")
		}
	})
}
