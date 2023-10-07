package util

import (
	"database/sql"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type InterfaceTransaction interface {
	GetDB() *sql.DB
}

func ExecTransactionService(i InterfaceTransaction, txFunc func(*sql.Tx) OutputBasicServiceInterface) (out OutputBasicInterface) {
	tx, err := i.GetDB().Begin()
	if err != nil {
		log.Errorf("Error Begin Transaction: %v", err)
		out = &OutputBasic{
			Code:    http.StatusInternalServerError,
			Result:  "Error Start Transaction",
			Message: err,
		}
		return
	}
	defer func() {
		if p := recover(); p != nil {
			if err := tx.Rollback(); err != nil {
				log.Errorf("Error tx.Rollback at Panic recovery: %v", err)

			}
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				log.Errorf("Error tx.Rollback:%v", err)
			}
		} else {
			err = tx.Commit()
			if err != nil {
				log.Errorf("Error tx.Commit: %v", err)
				out = &OutputBasic{
					Code:    http.StatusInternalServerError,
					Result:  "NG",
					Message: err,
				}
			}
		}
	}()
	txOut := txFunc(tx)
	out = txOut
	err = txOut.GetError()
	return
}
