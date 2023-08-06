package main

import (
	"flag"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/yzmw1213/demo-api/route"
)

func main() {
	addr := flag.String("addr", ":9090", "アプリケーションのアドレス")

	route.IndexRoute()
	log.Infof("Web server start port:[%s]", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Errorf("Error http.ListenAndServe %v", err)
	}
}