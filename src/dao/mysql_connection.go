package dao

import (
	"database/sql"
	"os"

	"github.com/yzmw1213/demo-api/conf"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

var db *sql.DB
var initedDB = false

func InitMysql() {
	log.Infof("InitMysql")
	var err error
	dataSouece := ""
	// TODO 環境変数を集中管理できる仕組みを作る
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dataSouece += user + ":" + password + "@" + conf.GetDBHost() + "/" + "user_db" + "?charset=utf8mb4&parseTime=true"

	db, err = sql.Open("mysql", dataSouece)
	if err != nil {
		log.Errorf("initMysql err %v", err)
		return
	}
	err = db.Ping()
	if err != nil {
		log.Errorf("PingDB err %v", err)
		return
	}

	initedDB = true
}
