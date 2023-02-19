package database

import (
	"testing"
)

func TestSensorRepository_ReserveNewId(t *testing.T) {
	t.Run("it reserves unique ids", func(t *testing.T) {
		register := NewSensorRepository()
		alreadyReservedIds := make(map[uint8]interface{})
		var empty struct{}
		for i := uint8(0); i < ^uint8(0); i += 1 {
			id, err := register.ReserveNewId()
			if err != nil {
				t.Error("SensorRepository couldn't reserve id but it should")
				return
			}
			if _, exists := alreadyReservedIds[id]; exists {
				t.Error("SensorRepository reserved the same id twice:", id)
				return
			}
			alreadyReservedIds[id] = empty // ciekawe co by sie stalo jakby wypelnic to nilami
		}
	})

	t.Run("it doesn't reserve too many ids", func(t *testing.T) {
		register := NewSensorRepository()
		for i := uint8(0); i < ^uint8(0); i += 1 {
			_, _ = register.ReserveNewId()
		}
		_, err := register.ReserveNewId()
		if err == nil {
			t.Error("SensorRepository didn't return error while trying to reserve more id than it is capable of, but it should")
		}
	})
}

func TestSensorRegister_DoesIdExist(t *testing.T) {
	register := NewSensorRepository()
	id, _ := register.ReserveNewId()
	if !register.DoesIdExist(id) {
		t.Error("SensorRepository says id doesn't exist but it should")
		return
	}
	if register.DoesIdExist(id + 1) {
		t.Error("SensorRepository says id exist but sholudn't")
		return
	}
}

func TestSensorRegister_ReleaseId(t *testing.T) {
	t.Run("it releases id", func(t *testing.T) {
		register := NewSensorRepository()
		id, _ := register.ReserveNewId()
		err := register.ReleaseId(id)
		if err != nil {
			t.Error("SensorRepository throws error while releasing id but it shouldn't")
			return
		}
		err = register.ReleaseId(id)
		if err == nil {
			t.Error("SensorRepository doesn't throw error while releasing non existent id but it should")
			return
		}
	})

	t.Run("it does not release id in use", func(t *testing.T) {
		register := NewSensorRepository()
		id, _ := register.ReserveNewId()
		_ = register.MarkSensorAsBeingUsed(id)
		if err := register.ReleaseId(id); err == nil {
			t.Error("SensorRepository didn't throw error at try of releasing used id, but it should")
			return
		}
		_ = register.MarkSensorAsUnused(id)
		if err := register.ReleaseId(id); err != nil {
			t.Error("SensorRepository returned error while trying to release unused id, but it shouldn't")
		}
	})
}

func TestSensorRegister_MarkSensorAsBeingUsed(t *testing.T) {
	register := NewSensorRepository()
	id, _ := register.ReserveNewId()
	err := register.MarkSensorAsBeingUsed(id)
	if err != nil {
		t.Error("SensorRepository couldn't mark id as being used, but it should")
		return
	}
	err = register.MarkSensorAsBeingUsed(id)
	if err == nil {
		t.Error("SensorRepository didn't return error after trying to mark id as used twice, but it should")
		return
	}
	if !register.IsIdUsed(id) {
		t.Error("SensorRepository says id isn't used, but it should be")
		return
	}
	err = register.MarkSensorAsUnused(id)
	if err != nil {
		t.Error("SensorRepository couldn't mark sensor as unused, but it should")
		return
	}
	err = register.MarkSensorAsUnused(id)
	if err == nil {
		t.Error("SensorRepository didn't return error after trying to mark sensor as unused twice, but it should")
		return
	}
	if register.IsIdUsed(id) {
		t.Error("SensorRepository says id is used, but it shouldn't be")
		return
	}
}

func TestSensorRegister_GetById(t *testing.T) {
	register := NewSensorRepository()
	id, _ := register.ReserveNewId()
	sensor, err := register.GetById(id)
	if sensor == nil || err != nil {
		t.Error("SensorRepository couldn't get sensor by id, but it should")
		return
	}
	if sensor.ID != id {
		t.Error("ID on sensor object is not equal to requested ID")
		return
	}
	sensor, err = register.GetById(id + 1)
	if sensor != nil || err == nil {
		t.Error("SensorRepository returned sensor or did not return an error while getting non existent id, but it should")
		return
	}
}
