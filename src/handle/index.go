package handle

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/yzmw1213/demo-api/util"
)

type IndexHandler struct {
}

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

func (h *IndexHandler) IndexHandler(c *gin.Context) {
	log.Infof("IndexHandler start")
	name, _ := c.Get("user_name")
	log.Infof("name: %s", name.(string))
	res := &util.OutputBasic{
		Code:    http.StatusOK,
		Result:  "OK",
		Message: "OK",
	}
	c.JSON(
		res.GetCode(),
		res.GetResult(),
	)
}
