package dao

import (
	"context"
	"os"

	firebase "firebase.google.com/go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

var app *firebase.App
var initedFB = false

func initFirebase() {
	var err error

	opt := option.WithCredentialsFile(os.Getenv("FIREBASE_KEY_FILE_JSON"))
	// https://firebase.google.com/docs/admin/setup?hl=ja#initialize-sdk
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Errorf("firebase.NewApp err: %v\n", err)
		return
	}
	initedFB = true
	log.Infof("InitFirebase")
}

func Firebase() *firebase.App {
	if !initedFB {
		initFirebase()
	}
	return app
}
