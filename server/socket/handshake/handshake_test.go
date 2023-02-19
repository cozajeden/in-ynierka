package handshake

import (
	"bytes"
	"errors"
	"reflect"
	"server/database"
	"server/socket/codes"
	"testing"
)

type sensorRegisterMock struct {
	database.SensorRepository
	isIdAvailable bool
	reservedId    uint8
	idUsed        bool
	idExists      bool
}

func (srm *sensorRegisterMock) ReserveNewId() (uint8, error) {
	if !srm.isIdAvailable {
		return 0, errors.New("no id available")
	}
	return srm.reservedId, nil
}

func (srm *sensorRegisterMock) IsIdUsed(uint8) bool {
	return srm.idUsed
}

func (srm *sensorRegisterMock) DoesIdExist(uint8) bool {
	return srm.idExists
}

type testCase struct {
	incoming             []byte
	expected             []byte
	isIdAvailable        bool
	idExists             bool
	idIsUsedAlready      bool
	shouldErrorBePresent bool
}

var testCases = []testCase{
	{[]byte{codes.RequireAuth, 0}, []byte{codes.Ok, 100}, true, false, false, false},
	{[]byte{codes.RequireAuth, 0xFF}, []byte{codes.Ok, 100}, true, false, false, false},
	{[]byte{codes.RequireAuth, 0xF0}, []byte{codes.Ok, 100}, true, false, false, false},
	{[]byte{codes.RequireAuth, 0}, []byte{codes.NoIdAvailable, 0}, false, false, false, true},
	{[]byte{codes.ImAuthorized, 100}, []byte{codes.Ok, 100}, false, true, false, false},
	{[]byte{codes.ImAuthorized, 100}, []byte{codes.IdIsAlreadyConnected, 100}, false, true, true, true},
	{[]byte{codes.ImAuthorized, 100}, []byte{codes.IdDoesntExist, 100}, false, false, false, true},
	{[]byte{0xF0, 100}, []byte{codes.HandshakeNotValid, 100}, false, false, false, true},
}

func TestAuthRequestHandshake(t *testing.T) {
	test := func(tc testCase) {
		reader := bytes.NewReader(tc.incoming)
		handshakeHandler := NewHandler(&sensorRegisterMock{
			isIdAvailable: tc.isIdAvailable,
			reservedId:    tc.expected[1],
			idExists:      tc.idExists,
			idUsed:        tc.idIsUsedAlready,
		})

		_, response, err := handshakeHandler.Handle(reader)

		if !reflect.DeepEqual(response, tc.expected) {
			t.Errorf("Handling %v returned %v, expected %v", tc.incoming, response, tc.expected)
		}
		if err == nil && tc.shouldErrorBePresent {
			t.Errorf("Handling %v returned no error, but it was expected", tc.incoming)
		}
		if err != nil && !tc.shouldErrorBePresent {
			t.Errorf("Handling %v returned an error, but it wasn't expected", tc.incoming)
		}
	}

	for _, tc := range testCases {
		test(tc)
	}
}
