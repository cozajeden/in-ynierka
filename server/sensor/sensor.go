package sensor

type Sensor struct {
	ID uint8 `json:"id"`
}

type ReadingType struct {
	Format uint8
	Bytes  int
	Count  int
	Labels []uint8
}

const (
	Imu = iota + 1
	Microphone
)

const (
	I16 = iota
	F32
)

const (
	Volume = uint8(iota)
	Xg
	Yg
	Zg
	Xdps
	Ydps
	Zdps
)

var AvailableReadingTypes = map[uint8]ReadingType{
	Microphone: {I16, 2, 1, []byte{Volume}},
	Imu:        {F32, 4, 6, []byte{Xg, Yg, Zg, Xdps, Ydps, Zdps}},
}
