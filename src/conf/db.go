package conf

import "github.com/yzmw1213/demo-api/util"

type DBConfig struct {
	Host     string
	User     string
	Password string
}

type EnvConfig struct {
	DB DBConfig
}

func GetDBHost() string {
	if util.IsDevelopment() {
		return "tcp(db:3306)"
	}
	// TODO production host
	return ""
}
