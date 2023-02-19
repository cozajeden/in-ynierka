package serving

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Gui(c *gin.Context) {
	sensors := sensorRegister.GetAllActiveSensors()
	micOptions := map[int]string{
		22050: "22050 Hz",
		44100: "44100 Hz",
	}
	imuOptionsDPS := map[int]string{
		250:  "250 DPS",
		500:  "500 DPS",
		1000: "1000 DPS",
		2000: "2000 DPS",
	}
	imuOptionsG := map[int]string{
		2:  "2G",
		4:  "4G",
		8:  "8G",
		16: "16G",
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"sensors":       sensors,
		"imuOptionsDPS": imuOptionsDPS,
		"imuOptionsG":   imuOptionsG,
		"micOptions":    micOptions,
	})
}
