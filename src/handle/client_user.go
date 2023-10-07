package handle

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/yzmw1213/demo-api/service"
	"github.com/yzmw1213/demo-api/util"
)

type ClientUserHandler struct {
	clientUserService *service.ClientUserService
}

func NewClientUserHandler() *ClientUserHandler {
	return &ClientUserHandler{
		service.NewClientUserService(),
	}
}

func (h *ClientUserHandler) GetHandle(c *gin.Context) {
	log.Infof("GetHandle star")
	name, _ := c.Get("user_name")
	userID, _ := c.Get("user_id")
	log.Infof("user_id: %d,name: %s", userID.(int64), name.(string))

	var in service.InputGetClientUser
	if err := c.ShouldBind(in); err != nil {
		util.BadRequestJson(*c, err)
	}
	out := h.clientUserService.GetClientUser(in.GetParam())
	c.JSON(
		out.GetCode(),
		out.GetResult(),
	)
}

//func (h *UserHandler) SaveHandle(c *gin.Context) {
//	log.Infof("SaveHandle start")
//	name, _ := c.Get("user_name")
//	userID, _ := c.Get("user_id")
//	log.Infof("user_id: %d,name: %s", userID.(int64), name.(string))
//
//	var in service.InputGetUser
//	if err := c.ShouldBind(in); err != nil {
//		util.BadRequestJson(*c, err)
//	}
//	out := h.userService.GetUser(&in)
//	c.JSON(
//		out.GetCode(),
//		out.GetResult(),
//	)
//}
