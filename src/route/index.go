package route

import (
	"github.com/gin-gonic/gin"
	"github.com/yzmw1213/demo-api/conf"
	"github.com/yzmw1213/demo-api/handle"
	"github.com/yzmw1213/demo-api/middleware"
)

func IndexRoute(router *gin.Engine) {
	indexHandler := handle.NewIndexHandler()

	authorityClient := []string{conf.CustomUserClaimClient}
	router.GET("/", middleware.AuthAPI(indexHandler.IndexHandler, authorityClient))
}
