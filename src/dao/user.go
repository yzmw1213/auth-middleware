package dao

import (
	"database/sql"
	"strings"

	"github.com/yzmw1213/demo-api/conf"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/yzmw1213/demo-api/entity"
)

type UserDao struct {
	db *sql.DB
}

func NewUserDao() *UserDao {
	return &UserDao{
		DB(),
	}
}

func (d *UserDao) Get(tx *sql.Tx, in *entity.User) (list []*entity.User, err error) {
	log.Infof("UserDao.Get %v", in)
	query := `
SELECT
	uu.user_id,
	uu.name,
	uu.email,
	uu.firebase_uid,
	uu.delete_flag
FROM
	user_db.users uu
WHERE
	uu.delete_flag = ?
`
	var params []interface{}
	params = append(params, in.DeleteFlag)
	var whereString []string
	if in.UserID > 0 {
		whereString = append(whereString, " uu.user_id = ? ")
		params = append(params, in.UserID)
	}
	if in.FirebaseUID != "" {
		whereString = append(whereString, " uu.firebase_uid = ? ")
		params = append(params, in.FirebaseUID)
	}
	log.Infof("query:%s params:%v", query, params)

	if len(whereString) > 0 {
		query += " AND " + strings.Join(whereString, " AND ")
	}
	var rows *sql.Rows

	if tx != nil {
		rows, err = tx.Query(query, params...)
	} else {
		rows, err = d.db.Query(query, params...)
	}
	if err != nil {
		if err != sql.ErrNoRows {
			log.Errorf("Error Query %v", err)
		}
		return
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Errorf("Error Rows Close:%v", err)
		}
	}()

	list = []*entity.User{}
	for rows.Next() {
		row := entity.User{}
		if e := rows.Scan(
			&row.UserID,
			&row.Name,
			&row.Email,
			&row.FirebaseUID,
			&row.DeleteFlag,
		); e != nil {
			list = []*entity.User{}
			err = e
			log.Errorf("Error Scan %v", err)
		}
		list = append(list, &row)
	}

	return
}

func (d *UserDao) GetEnable(tx *sql.Tx, in *entity.User) (list []*entity.User, err error) {
	log.Infof("UserDao.GetEnable %v", in)
	in.DeleteFlag = conf.DeleteFlagOFF
	return d.Get(tx, in)
}

func (d *UserDao) GetDisable(tx *sql.Tx, in *entity.User) (list []*entity.User, err error) {
	log.Infof("UserDao.GetDisable %v", in)
	in.DeleteFlag = conf.DeleteFlagON
	return d.Get(tx, in)
}
