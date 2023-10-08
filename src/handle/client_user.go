package handle

import (
	"encoding/csv"
	"errors"
	"net/http"
	"time"

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
		return
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

func (h *ClientUserHandler) DownloadClientUserCsvHandle(c *gin.Context) {
	log.Infof("DownloadClientUserCsvHandle star")
	var in service.InputGetClientUser
	if err := c.ShouldBind(in); err != nil {
		util.BadRequestJson(*c, err)
		return
	}
	csvList, err := h.clientUserService.GetClientListCSV(in.GetParam())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	now := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60))
	fileName := "client_user_list_" + now.Format("20060102030405") + ".csv"

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename='"+fileName+"'")
	w := csv.NewWriter(c.Writer)
	w.WriteAll(csvList)
}
