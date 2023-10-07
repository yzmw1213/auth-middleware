package handle

import (
	"errors"

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

func (h *ClientUserHandler) SaveHandle(c *gin.Context) {
	log.Infof("SaveHandle start")
	name, _ := c.Get("user_name")
	userID, _ := c.Get("user_id")
	log.Infof("user_id: %d,name: %s", userID.(int64), name.(string))

	in := &service.InputSaveClientUser{}
	if err := c.Bind(in); err != nil {
		util.BadRequestJson(*c, err)
		return
	}
	in.UpdateUserID = userID.(int64)
	// TODO バリデーションツールを入れる
	if in.Email == "" || in.Name == "" || in.Password == "" {
		log.Warn("invalid input")
		util.BadRequestJson(*c, errors.New("invalid input"))
		return
	}

	out := h.clientUserService.Save(in)
	c.JSON(
		out.GetCode(),
		out.GetResult(),
	)
}
