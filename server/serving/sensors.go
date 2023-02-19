package serving

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Sensors(c *gin.Context) {
	registredSensors := sensorRegister.GetAllSensors()

	c.JSON(http.StatusOK, gin.H{
		"sensor_list": registredSensors,
	})
}
