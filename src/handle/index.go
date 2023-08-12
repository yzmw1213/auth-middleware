package handle

import (
	"github.com/gin-gonic/gin"
	"net/http"

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
	res := &util.OutputBasic{
		Code: http.StatusOK,
		Result: "OK",
		Message: "OK" ,
	}
	c.JSON(
		res.GetCode(),
		res.GetResult(),
	)
}