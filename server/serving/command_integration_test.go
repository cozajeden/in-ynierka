//go:build integration
// +build integration

package serving

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"server/database"
	"server/socket/codes"
	"server/socket/command"
	"testing"
	"time"
)

type mockedReadingsRepository struct{}

func (m mockedReadingsRepository) SaveReadings(uint8, []*database.Reading) error {
	return nil
}

func (m mockedReadingsRepository) GetReadingsSince(time.Time, uint8) ([]*database.Reading, error) {
	return []*database.Reading{}, nil
}

func TestCommand(t *testing.T) {

	testCases := []struct {
		body            CommandBody
		expectedCommand command.Command
	}{
		{
			body:            CommandBody{Command: TurnOn, SettingsDetails: &Settings{SensorType: codes.Imu}},
			expectedCommand: command.Command{Command: codes.Toggle, SensorType: codes.Imu, Arguments: []byte{codes.TurnOn}},
		},
		{
			body:            CommandBody{Command: TurnOff, SettingsDetails: &Settings{SensorType: codes.Imu}},
			expectedCommand: command.Command{Command: codes.Toggle, SensorType: codes.Imu, Arguments: []byte{codes.TurnOff}},
		},
		{
			body:            CommandBody{Command: TurnOn, SettingsDetails: &Settings{SensorType: codes.Mic}},
			expectedCommand: command.Command{Command: codes.Toggle, SensorType: codes.Mic, Arguments: []byte{codes.TurnOn}},
		},
		{
			body:            CommandBody{Command: TurnOff, SettingsDetails: &Settings{SensorType: codes.Mic}},
			expectedCommand: command.Command{Command: codes.Toggle, SensorType: codes.Mic, Arguments: []byte{codes.TurnOff}},
		},
		{
			body:            CommandBody{Command: ImuSettings, SettingsDetails: &Settings{SensorType: codes.Imu, DPS: 250, G: 2}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Imu, Arguments: []byte{codes.Imu250DPS, codes.Imu2G}},
		},
		{
			body:            CommandBody{Command: ImuSettings, SettingsDetails: &Settings{SensorType: codes.Imu, DPS: 500, G: 2}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Imu, Arguments: []byte{codes.Imu500DPS, codes.Imu2G}},
		},
		{
			body:            CommandBody{Command: ImuSettings, SettingsDetails: &Settings{SensorType: codes.Imu, DPS: 1000, G: 2}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Imu, Arguments: []byte{codes.Imu1000DPS, codes.Imu2G}},
		},
		{
			body:            CommandBody{Command: ImuSettings, SettingsDetails: &Settings{SensorType: codes.Imu, DPS: 2000, G: 2}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Imu, Arguments: []byte{codes.Imu2000DPS, codes.Imu2G}},
		},
		{
			body:            CommandBody{Command: ImuSettings, SettingsDetails: &Settings{SensorType: codes.Imu, DPS: 250, G: 4}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Imu, Arguments: []byte{codes.Imu250DPS, codes.Imu4G}},
		},
		{
			body:            CommandBody{Command: ImuSettings, SettingsDetails: &Settings{SensorType: codes.Imu, DPS: 500, G: 4}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Imu, Arguments: []byte{codes.Imu500DPS, codes.Imu4G}},
		},
		{
			body:            CommandBody{Command: ImuSettings, SettingsDetails: &Settings{SensorType: codes.Imu, DPS: 1000, G: 4}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Imu, Arguments: []byte{codes.Imu1000DPS, codes.Imu4G}},
		},
		{
			body:            CommandBody{Command: ImuSettings, SettingsDetails: &Settings{SensorType: codes.Imu, DPS: 2000, G: 4}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Imu, Arguments: []byte{codes.Imu2000DPS, codes.Imu4G}},
		},
		{
			body:            CommandBody{Command: ImuSettings, SettingsDetails: &Settings{SensorType: codes.Imu, DPS: 250, G: 8}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Imu, Arguments: []byte{codes.Imu250DPS, codes.Imu8G}},
		},
		{
			body:            CommandBody{Command: ImuSettings, SettingsDetails: &Settings{SensorType: codes.Imu, DPS: 500, G: 8}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Imu, Arguments: []byte{codes.Imu500DPS, codes.Imu8G}},
		},
		{
			body:            CommandBody{Command: ImuSettings, SettingsDetails: &Settings{SensorType: codes.Imu, DPS: 1000, G: 8}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Imu, Arguments: []byte{codes.Imu1000DPS, codes.Imu8G}},
		},
		{
			body:            CommandBody{Command: ImuSettings, SettingsDetails: &Settings{SensorType: codes.Imu, DPS: 2000, G: 8}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Imu, Arguments: []byte{codes.Imu2000DPS, codes.Imu8G}},
		},
		{
			body:            CommandBody{Command: ImuSettings, SettingsDetails: &Settings{SensorType: codes.Imu, DPS: 250, G: 16}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Imu, Arguments: []byte{codes.Imu250DPS, codes.Imu16G}},
		},
		{
			body:            CommandBody{Command: ImuSettings, SettingsDetails: &Settings{SensorType: codes.Imu, DPS: 500, G: 16}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Imu, Arguments: []byte{codes.Imu500DPS, codes.Imu16G}},
		},
		{
			body:            CommandBody{Command: ImuSettings, SettingsDetails: &Settings{SensorType: codes.Imu, DPS: 1000, G: 16}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Imu, Arguments: []byte{codes.Imu1000DPS, codes.Imu16G}},
		},
		{
			body:            CommandBody{Command: ImuSettings, SettingsDetails: &Settings{SensorType: codes.Imu, DPS: 2000, G: 16}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Imu, Arguments: []byte{codes.Imu2000DPS, codes.Imu16G}},
		},
		{
			body:            CommandBody{Command: MicSettings, SettingsDetails: &Settings{SensorType: codes.Mic, Hz: 22050}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Mic, Arguments: []byte{codes.Mic22050Hz}},
		},
		{
			body:            CommandBody{Command: MicSettings, SettingsDetails: &Settings{SensorType: codes.Mic, Hz: 44100}},
			expectedCommand: command.Command{Command: codes.Settings, SensorType: codes.Mic, Arguments: []byte{codes.Mic44100Hz}},
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("it sends appropriate command to exchange on %v", tc.body.Command), func(t *testing.T) {
			sr := database.NewSensorRepository()
			testId, err := sr.ReserveNewId()
			if err != nil {
				t.Errorf("failed to reserve new id: %v", err)
			}
			tc.body.Id = testId
			tc.expectedCommand.ForId = testId
			rr := &mockedReadingsRepository{}
			l := logrus.New()
			l.Out = ioutil.Discard
			ce := command.NewExchange()
			sub := ce.GetSubscriberChannel()
			router := SetupRouter(l, sr, ce)
			body, err := json.Marshal(tc.body)
			if err != nil {
				return
			}

			response := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/command", bytes.NewReader(body))
			router.ServeHTTP(response, req)

			if response.Code != http.StatusOK {
				t.Errorf("something gone wrong, got response code %v: %v", response.Code, response.Body.String())
			}
			select {
			case c := <-sub:
				if !reflect.DeepEqual(tc.expectedCommand, *c) {
					t.Errorf("received command does not match, expected %v, got %v", tc.expectedCommand, *c)
				}
			default:
				t.Errorf("no message appeared in sub channel")
			}
		})
	}
}
