package socket

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"server/database"
	"server/monitor"
	"server/sensor"
	"server/socket/command"
	"server/socket/handshake"
	"server/socket/interpreting"
	"server/socket/receive"
	"time"
)

const interpreterThreads = 5

type listener struct {
	logger           *logrus.Logger
	sr               database.SensorRepository
	hh               *handshake.Handler
	port             string
	rr               database.ReadingsRepository
	interpreterQueue chan *receive.ReceiverOutput
	ce               *command.Exchange
}

type connectionHandler struct {
	conn                     net.Conn
	logger                   *logrus.Entry
	handshakeHandler         *handshake.Handler
	sensorRepository         database.SensorRepository
	readingsRepository       database.ReadingsRepository
	interpreterQueue         chan *receive.ReceiverOutput
	commandExchange          *command.Exchange
	commandRoutineKillSwitch chan bool
	connectionAbortSignal    chan bool
}

func (l *listener) Listen() {
	var err error
	defer func() {
		if err != nil {
			l.logger.WithFields(map[string]interface{}{
				"error": err,
			}).Error("Error occured while listening")
		}
	}()

	l.createInterpreterThreads()

	netListener, err := net.Listen("tcp", fmt.Sprintf(":%v", l.port))
	if err != nil {
		return
	}
	defer func() {
		l.logger.Debug("Closing connection now")
		netListener.Close()
	}()
	l.logger.WithFields(map[string]interface{}{
		"port": l.port,
	}).Debug("Created TCP server")

	for {
		conn, err := netListener.Accept()
		if err != nil {
			l.logger.WithFields(map[string]interface{}{
				"error": err,
			}).Error("Could not read incoming connection")
			continue
		}
		handler := l.connHandler(conn)
		go handler.handleConnection()
	}
}

func (l *listener) connHandler(conn net.Conn) *connectionHandler {
	logger := l.logger.WithFields(map[string]interface{}{
		"address": conn.RemoteAddr(),
	})
	return &connectionHandler{
		conn:                     conn,
		logger:                   logger,
		handshakeHandler:         l.hh,
		sensorRepository:         l.sr,
		readingsRepository:       l.rr,
		interpreterQueue:         l.interpreterQueue,
		commandExchange:          l.ce,
		commandRoutineKillSwitch: make(chan bool),
	}
}

func (l *listener) createInterpreterThreads() {
	for i := 0; i < interpreterThreads; i++ {
		go l.interpret()
	}
}

func (l *listener) interpret() {
	for {
		message := <-l.interpreterQueue
		start := time.Now()
		monitor.InterpreterQueueLength(len(l.interpreterQueue))
		monitor.SendBytes(len(message.Message))
		interpreter := interpreting.NewInterpreter(
			bufio.NewReader(bytes.NewBuffer(message.Message)), l.rr, message.Id)
		interpreted, err := interpreter.Interpret()
		if err != nil {
			l.logger.Debugf("Failed to interpret message: %v\n", err)
			continue
		}
		err = interpreted.Act()
		if err != nil {
			l.logger.Debugf("failed to act on interpreted message: %v\n", err)
		}
		monitor.RegisterInterpretingTime(time.Now().Sub(start))
	}
}

func (ch *connectionHandler) handleConnection() {
	var err error
	ch.logger.Debug("Opened new connection")
	defer func() {
		ch.logger.Debug("Closing connection")
		_ = ch.conn.Close()
	}()

	var id uint8
	var handshakeResponse []byte

	for {
		id, handshakeResponse, err = ch.handshakeHandler.Handle(ch.conn)
		if err == nil {
			break
		}
		ch.logger.Errorf("error occured while handling handshake: %v\n", err)
		if err == io.EOF {
			return
		}
		_, err = ch.conn.Write(handshakeResponse)
		if err != nil {
			ch.logger.Errorf(
				"error occured while trying to write to connection: %v\n", err)
			return
		}
	}
	sensorObject := ch.useSensor(id)
	go ch.startCommandRoutine(id)
	defer func() {
		ch.unuseSensor(sensorObject)
		ch.commandRoutineKillSwitch <- true
	}()
	_, err = ch.conn.Write(handshakeResponse)
	if err != nil {
		ch.logger.Errorf("error while responding to handshake: %v\n", err)
	}

	receiver := ch.buildReceiver(id)
	for {
		select {
		case <-ch.connectionAbortSignal:
			return
		default:
			err = receiver.Receive()
			if err != nil {
				ch.logger.Errorf("error occured while receiving data: %v\n", err)
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					return
				}
			}
		}
	}
}

func (ch *connectionHandler) addLoggerData(key string, value interface{}) {
	ch.logger = ch.logger.WithFields(map[string]interface{}{
		key: value,
	})
}

func (ch *connectionHandler) useSensor(id uint8) *sensor.Sensor {
	sensorObject, _ := ch.sensorRepository.GetById(id)
	_ = ch.sensorRepository.MarkSensorAsBeingUsed(id)
	ch.addLoggerData("id", id)
	ch.logger.Debug("Established connection with sensor")
	return sensorObject
}

func (ch *connectionHandler) unuseSensor(s *sensor.Sensor) {
	_ = ch.sensorRepository.MarkSensorAsUnused(s.ID)
	ch.logger.Debug("Connection with sensor has been broken")
}

func (ch *connectionHandler) buildReceiver(id uint8) receive.Receiver {
	return receive.NewReceiver(ch.conn, ch.interpreterQueue, id)
}

func (ch *connectionHandler) startCommandRoutine(id uint8) {
	sub := ch.commandExchange.GetSubscriberChannel()
	defer ch.commandExchange.Unsubscribe(sub)
	for {
		select {
		case <-ch.commandRoutineKillSwitch:
			return
		case comm := <-sub:
			if id != comm.ForId {
				continue
			}
			_, err := ch.conn.Write(append([]byte{comm.Command, comm.SensorType}, comm.Arguments...))
			if err != nil {
				ch.logger.Errorf("error while writing command: %v", err)
				ch.connectionAbortSignal <- true
			}
		}
	}
}
