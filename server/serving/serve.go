package serving

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"html/template"
	"server/database"
	"server/socket/command"
)

var sensorRegister database.SensorRepository
var exchange *command.Exchange

func SetupRouter(l *logrus.Logger, sr database.SensorRepository, ce *command.Exchange) *gin.Engine {
	sensorRegister = sr
	exchange = ce
	gin.DefaultErrorWriter = l.WriterLevel(logrus.ErrorLevel)
	gin.DefaultWriter = l.Writer()
	r := gin.Default()
	r.Static("/assets", "./gui/assets")
	html, err := template.New("index.html").Delims("[[", "]]").ParseFiles("gui/index.html")
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(html)
	_ = r.SetTrustedProxies(nil)
	r.GET("/", Gui)
	r.GET("/sensors", Sensors)
	r.POST("/command", Command)
	return r
}
