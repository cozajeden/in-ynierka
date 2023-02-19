package serving

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"server/socket/codes"
	"server/socket/command"
)

type Settings struct {
	DPS        int `json:"dps"`
	G          int `json:"g"`
	Hz         int `json:"hz"`
	SensorType int `json:"sensor_type"`
}

type CommandBody struct {
	Id              uint8     `json:"id"`
	Command         string    `json:"command"`
	SettingsDetails *Settings `json:"settings"`
}

const (
	TurnOff     = "turn-off"
	TurnOn      = "turn-on"
	ImuSettings = "imu-settings"
	MicSettings = "mic-settings"
)

var commandMappings = map[string]command.Command{
	TurnOff: {
		Command:   codes.Toggle,
		Arguments: []byte{codes.TurnOff},
	},
	TurnOn: {
		Command:   codes.Toggle,
		Arguments: []byte{codes.TurnOn},
	},
	ImuSettings: {
		Command:    codes.Settings,
		SensorType: codes.Imu,
		Arguments:  make([]byte, 2),
	},
	MicSettings: {
		Command:    codes.Settings,
		SensorType: codes.Mic,
		Arguments:  make([]byte, 1),
	},
}

var gMapping = map[int]byte{
	2:  codes.Imu2G,
	4:  codes.Imu4G,
	8:  codes.Imu8G,
	16: codes.Imu16G,
}

var dpsMapping = map[int]byte{
	250:  codes.Imu250DPS,
	500:  codes.Imu500DPS,
	1000: codes.Imu1000DPS,
	2000: codes.Imu2000DPS,
}

var hzMapping = map[int]byte{
	22050: codes.Mic22050Hz,
	44100: codes.Mic44100Hz,
}

func Command(c *gin.Context) {
	var json CommandBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := dispatchCommand(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func dispatchCommand(cb *CommandBody) error {
	c, ok := commandMappings[cb.Command]
	if !ok {
		return fmt.Errorf("command '%v' not found", cb.Command)
	}
	if !sensorRegister.DoesIdExist(cb.Id) {
		return fmt.Errorf("id %v not found", cb.Id)
	}
	err := fillWithSettingsDetails(&c, cb)
	if err != nil {
		return err
	}
	c.ForId = cb.Id
	exchange.Dispatch(&c)
	return nil
}

func fillWithSettingsDetails(c *command.Command, cb *CommandBody) error {
	if c.Command == codes.Toggle {
		if cb.SettingsDetails.SensorType == codes.Undefined || cb.SettingsDetails.SensorType > codes.Mic {
			return fmt.Errorf("unknown sensor type code: %v", cb.SettingsDetails.SensorType)
		}
		c.SensorType = byte(cb.SettingsDetails.SensorType)
		return nil
	}

	if c.Command != codes.Settings { // every code requires arguments
		return fmt.Errorf("unknown command: %v", c.Command)
	}

	if cb.SettingsDetails.SensorType == codes.Mic {
		hz, ok := hzMapping[cb.SettingsDetails.Hz]
		if !ok {
			return fmt.Errorf("unknown Hz value: %v", cb.SettingsDetails.Hz)
		}
		c.Arguments[0] = hz
		return nil
	}

	// sensor must be imu
	dps, ok := dpsMapping[cb.SettingsDetails.DPS]
	if !ok {
		return fmt.Errorf("got unknown DPS value: %v", cb.SettingsDetails.DPS)
	}
	g, ok := gMapping[cb.SettingsDetails.G]
	if !ok {
		return fmt.Errorf("got unknown G value: %v", cb.SettingsDetails.G)
	}
	c.Arguments[0] = dps
	c.Arguments[1] = g
	return nil
}
