package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequestJson(c gin.Context, err error) {
	switch err.(type) {
	case error:
		c.JSON(
			http.StatusBadRequest,
			map[string]interface{}{
				"code":    http.StatusBadRequest,
				"result":  "NG",
				"message": err.Error(),
			},
		)
	}
	c.JSON(
		http.StatusBadRequest,
		map[string]interface{}{
			"code":    http.StatusBadRequest,
			"result":  "NG",
			"message": err,
		},
	)
}
