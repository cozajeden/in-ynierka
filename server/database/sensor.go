package database

import (
	"errors"
	"server/sensor"
)

type SensorRepository interface {
	ReserveNewId() (uint8, error)
	DoesIdExist(uint8) bool
	IsIdUsed(uint8) bool
	GetById(uint8) (*sensor.Sensor, error)
	MarkSensorAsBeingUsed(uint8) error
	ReleaseId(uint8) error
	MarkSensorAsUnused(id uint8) error
	GetAllSensors() []*sensor.Sensor
	GetAllActiveSensors() map[uint8]*sensor.Sensor
}

type mappedRepository struct {
	sensors        map[uint8]*sensor.Sensor
	usedSensorsIds map[uint8]interface{}
}

func (m *mappedRepository) GetAllActiveSensors() map[uint8]*sensor.Sensor {
	result := make(map[uint8]*sensor.Sensor, 0)
	for k := range m.usedSensorsIds {
		result[k] = m.sensors[k]
	}
	return result
}

func NewSensorRepository() SensorRepository {
	return &mappedRepository{
		map[uint8]*sensor.Sensor{},
		map[uint8]interface{}{},
	}
}

func (m *mappedRepository) GetAllSensors() []*sensor.Sensor {
	result := make([]*sensor.Sensor, 0, len(m.sensors))
	for _, s := range m.sensors {
		result = append(result, s)
	}
	return result
}

func (m *mappedRepository) ReserveNewId() (uint8, error) {
	for i := uint8(0); i < ^uint8(0); i += 1 {
		if !m.DoesIdExist(i) {
			m.sensors[i] = &sensor.Sensor{
				i,
			}
			return i, nil
		}
	}
	return 0, errors.New("no new IDs are available")
}

func (m *mappedRepository) DoesIdExist(id uint8) (exists bool) {
	_, exists = m.sensors[id]
	return
}

func (m *mappedRepository) IsIdUsed(id uint8) (exists bool) {
	_, exists = m.usedSensorsIds[id]
	return
}

func (m *mappedRepository) GetById(id uint8) (*sensor.Sensor, error) {
	s, ok := m.sensors[id]
	if !ok {
		return nil, errors.New("sensor with such id was not found")
	}
	return s, nil
}

func (m *mappedRepository) MarkSensorAsBeingUsed(id uint8) error {
	if m.IsIdUsed(id) {
		return errors.New("tried to use id already in use")
	}
	m.usedSensorsIds[id] = struct{}{}
	return nil
}

func (m *mappedRepository) MarkSensorAsUnused(id uint8) error {
	if !m.IsIdUsed(id) {
		return errors.New("tried to unuse id already unused")
	}
	delete(m.usedSensorsIds, id)
	return nil
}

func (m *mappedRepository) ReleaseId(id uint8) error {
	if !m.DoesIdExist(id) {
		return errors.New("tried to release non existent id")
	}
	if m.IsIdUsed(id) {
		return errors.New("tried to release used id")
	}
	delete(m.sensors, id)
	return nil
}
