package codes

// commands codes
const (
	RequireAuth byte = 1 + iota
	ImAuthorized
	SendData
	Settings
	Toggle
)

// response codes
const (
	Ok byte = 1 + iota
	IdDoesntExist
	IdIsAlreadyConnected
	NoIdAvailable
	HandshakeNotValid
	CommandNotRecognized
	UnexpectedPacketHeader
	UnexpectedPacketBody
	MeasureAlreadyStarted
	NoMeasureToStop
	UnknownCommand
)

// sensor types for commanding
const (
	Undefined = iota
	Imu
	Mic
)

// toggle arguments
const (
	TurnOn = iota + 1
	TurnOff
)

// Imu first arg
const (
	Imu250DPS = iota
	Imu500DPS
	Imu1000DPS
	Imu2000DPS
)

// Imu second arg
const (
	Imu2G = iota
	Imu4G
	Imu8G
	Imu16G
)

// MIC first arg
const (
	Mic22050Hz = iota
	Mic44100Hz
)
