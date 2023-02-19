package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"server/database"
	"server/monitor"
	"server/serving"
	"server/socket"
	"server/socket/command"
	"server/socket/handshake"
)

const (
	ListenPortEnvKey     = "LISTEN_PORT"
	ServePortEnvKey      = "SERVE_PORT"
	MonitoringEnvKey     = "ENABLE_MONITORING"
	MonitoringPathEnvKey = "MONITORING_PATH"
)

func main() {
	prod := os.Getenv("DONT_LOAD_DOTENV")
	if prod != "1" {
		if err := godotenv.Load(".env"); err != nil {
			panic("Could not load .env file. Does it exist?")
		}
	}
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05.999",
	})
	sensorRegister := database.NewSensorRepository()
	handshakeHandler := handshake.NewHandler(
		sensorRegister,
	)
	listenPort := os.Getenv(ListenPortEnvKey)
	servePort := os.Getenv(ServePortEnvKey)
	readingRepository := database.NewReadingsRepository()
	defer readingRepository.Close()
	commandExchange := command.NewExchange()

	listener := socket.NewListener(logger, sensorRegister, handshakeHandler, listenPort, readingRepository, commandExchange)
	go listener.Listen()

	if os.Getenv(MonitoringEnvKey) == "1" {
		monitor.FilePath = os.Getenv(MonitoringPathEnvKey)
		monitor.RunMonitoring()
	}

	router := serving.SetupRouter(logger, sensorRegister, commandExchange)
	err := router.Run(fmt.Sprintf(":%v", servePort))
	if err != nil {
		logger.Fatalf("Couldn't run server: %v", err)
	}
}
