package route

import (
	"github.com/gin-gonic/gin"
	"github.com/yzmw1213/demo-api/handle"
)



func IndexRoute(router *gin.Engine) {
	indexHandler := handle.NewIndexHandler()
	router.GET("/", indexHandler.IndexHandler)
}
