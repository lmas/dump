package server

import "github.com/gin-gonic/gin"

var (
	// If VERBOSE is set to true it will start to print extra log messages.
	DEBUG bool
)

func StartWeb(addr string, requesthandler func(c *gin.Context)) {
	if DEBUG == false {
		gin.SetMode(gin.ReleaseMode)
	}

	web := gin.Default()

	web.Static("/static", "./static")
	web.StaticFile("/game", "./templates/game.html")
	web.GET("/client", requesthandler)

	web.Run(addr)
}
