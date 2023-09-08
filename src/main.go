package main

import (
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/yzmw1213/demo-api/dao"
	"github.com/yzmw1213/demo-api/route"
)

func main() {
	dao.InitMysql()

	router := gin.Default()
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
		log.Infof("Defaulting to port %s", port)
	}
	route.IndexRoute(router)
	router.Run(":" + port)
}
