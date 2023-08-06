package route

import (
	"net/http"

	"github.com/yzmw1213/demo-api/handle"
)



func IndexRoute() {
	indexHandler := handle.NewIndexHandler()
	http.HandleFunc("/", indexHandler.IndexHandler)
}
