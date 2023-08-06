package handle

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/yzmw1213/demo-api/util"
)

type IndexHandler struct {
}

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

func (h *IndexHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("IndexHandler start")
	res := &util.OutputBasic{
		Code: http.StatusOK,
		Result: "OK",
		Message: "OK" ,
	}

	output, err := json.Marshal(res.GetResult())

	if err != nil {
		http.Error(w, "JSONの生成に失敗しました。", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}