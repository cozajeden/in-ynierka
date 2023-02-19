package command

type SensorType byte

type Command struct {
	ForId      uint8
	Command    byte
	SensorType byte
	Arguments  []byte
}
