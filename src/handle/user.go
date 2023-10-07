package handle

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/yzmw1213/demo-api/service"
	"github.com/yzmw1213/demo-api/util"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		service.NewUserService(),
	}
}

func (h *UserHandler) GetHandle(c *gin.Context) {
	log.Infof("GetHandle star")
	name, _ := c.Get("user_name")
	userID, _ := c.Get("user_id")
	log.Infof("user_id: %d,name: %s", userID.(int64), name.(string))

	var in service.InputGetUser
	if err := c.ShouldBind(in); err != nil {
		util.BadRequestJson(*c, err)
	}
	out := h.userService.GetUser(in.GetParam())
	c.JSON(
		out.GetCode(),
		out.GetResult(),
	)
}
