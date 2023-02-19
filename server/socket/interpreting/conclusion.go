package interpreting

import (
	"server/database"
)

type Conclusion interface {
	Act() error
}

type readConclusion struct {
	readings []*database.Reading
	repo     database.ReadingsRepository
	sensorId uint8
}

func (rc *readConclusion) Act() (err error) {
	err = rc.repo.SaveReadings(rc.sensorId, rc.readings)
	return
}
