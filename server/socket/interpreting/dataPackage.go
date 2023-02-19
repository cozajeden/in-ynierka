package interpreting

import "server/sensor"

type dataPackage struct {
	rt        sensor.ReadingType
	sec       uint32
	usec      uint32
	frequency float32
	quantity  uint16
}
