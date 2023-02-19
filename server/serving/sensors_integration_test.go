//go:build integration
// +build integration

package serving

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"server/database"
	"testing"
)

func Test_Sensors(t *testing.T) {
	t.Run("it gets list of sensors registered", func(t *testing.T) {
		sensorRegister := database.NewSensorRepository()
		_, _ = sensorRegister.ReserveNewId()
		_, _ = sensorRegister.ReserveNewId()
		_, _ = sensorRegister.ReserveNewId()
		readingsRepository := database.NewReadingsRepository()
		logger := logrus.New()
		logger.Out = ioutil.Discard
		router := SetupRouter(logger, sensorRegister, nil)

		response := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/sensors", nil)
		router.ServeHTTP(response, req)

		var responseMap map[string]interface{}
		_ = json.Unmarshal(response.Body.Bytes(), &responseMap)
		if list, exists := responseMap["sensor_list"]; !exists {
			t.Error("sensor list doesn't exist in response")
		} else {
			if reflect.ValueOf(list).Len() != 3 {
				t.Error("list length is not matching expected")
			}
		}
	})
}
