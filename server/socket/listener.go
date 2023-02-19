package socket

import (
	"github.com/sirupsen/logrus"
	"server/database"
	"server/socket/command"
	"server/socket/handshake"
	"server/socket/receive"
)

type Listener interface {
	Listen()
}

func NewListener(logger *logrus.Logger, register database.SensorRepository, handler *handshake.Handler, port string, repository database.ReadingsRepository, exchange *command.Exchange) Listener {
	interpreterQueue := make(chan *receive.ReceiverOutput, 100)
	return &listener{logger: logger, rr: repository, port: port, sr: register, hh: handler, interpreterQueue: interpreterQueue, ce: exchange}
}
