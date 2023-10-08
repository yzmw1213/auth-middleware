package route

import (
	"github.com/gin-gonic/gin"
	"github.com/yzmw1213/demo-api/conf"
	"github.com/yzmw1213/demo-api/handle"
	"github.com/yzmw1213/demo-api/middleware"
)

func IndexRoute(router *gin.Engine) {
	indexHandler := handle.NewIndexHandler()
	clientUserHandler := handle.NewClientUserHandler()

	authorityClient := []string{conf.CustomUserClaimClient}
	authorityAdmin := []string{conf.CustomUserClaimAdmin}
	router.GET("/", middleware.AuthAPI(indexHandler.IndexHandler, authorityClient))

	// client user
	router.GET("/client/user", middleware.AuthAPI(clientUserHandler.GetHandle, authorityAdmin))
	router.GET("/client/user/csv", middleware.AuthAPI(clientUserHandler.DownloadClientUserCsvHandle, authorityAdmin))
	router.POST("/client/user", middleware.AuthAPI(clientUserHandler.SaveHandle, authorityAdmin))
}
