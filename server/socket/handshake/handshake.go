package handshake

import (
	"fmt"
	"io"
	"server/database"
	"server/socket/codes"
)

type Handler struct {
	sensorRegister database.SensorRepository
}

const (
	size = 2
)

func NewHandler(register database.SensorRepository) *Handler {
	return &Handler{
		register,
	}
}

func (h *Handler) Handle(reader io.Reader) (id uint8, response []byte, err error) {
	buffer := make([]byte, size)
	_, err = io.ReadFull(reader, buffer)
	if err != nil {
		return id, []byte{}, err
	}
	command, id := buffer[0], buffer[1]

	switch command {
	case codes.RequireAuth:
		id, err = h.sensorRegister.ReserveNewId()
		if err != nil {
			return id, []byte{codes.NoIdAvailable, id}, fmt.Errorf("handshake failed: %v", err)
		}
		return id, []byte{codes.Ok, id}, err
	case codes.ImAuthorized:
		if !h.sensorRegister.DoesIdExist(id) {
			return id, []byte{codes.IdDoesntExist, id}, fmt.Errorf("given id does not exist")
		}
		if h.sensorRegister.IsIdUsed(id) {
			return id, []byte{codes.IdIsAlreadyConnected, id}, fmt.Errorf("given id is already connected")
		}
		return id, []byte{codes.Ok, id}, err
	}
	return id, []byte{codes.HandshakeNotValid, id}, fmt.Errorf("handshake is not valid")
}
