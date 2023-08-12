package main

import (
	"github.com/gin-gonic/gin"

	"github.com/yzmw1213/demo-api/route"
)

func main() {
	router := gin.Default()

	route.IndexRoute(router)
	router.Run(":9090")
}
