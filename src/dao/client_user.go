package dao

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/yzmw1213/demo-api/conf"
	"github.com/yzmw1213/demo-api/entity"
)

type ClientUserDao struct {
	db *sql.DB
}

func NewClientUserDao() *ClientUserDao {
	return &ClientUserDao{
		DB(),
	}
}

func (d *ClientUserDao) GetCount(tx *sql.Tx, in *entity.ClientUser) (count int64, err error) {
	log.Infof("ClientUserDao.GetCount %v", in)
	query := `
SELECT
	count(cu.client_user_id)
FROM
	client_db.client_users cu
INNER JOIN user_db.users uu ON uu.user_id = cu.user_id
WHERE
	cu.delete_flag = ?
`
	params := []interface{}{
		in.DeleteFlag,
	}
	var whereString []string

	if in.UserID > 0 {
		whereString = append(whereString, " cu.client_user_id = ? ")
		params = append(params, in.UserID)
	}
	if len(in.Name) > 0 {
		whereString = append(whereString, " uu.name LIKE ? ")
		params = append(params, "%"+in.Name+"%")
	}
	if in.FirebaseUID != "" {
		whereString = append(whereString, " uu.firebase_uid = ? ")
		params = append(params, in.FirebaseUID)
	}

	if len(whereString) > 0 {
		query += " WHERE " + strings.Join(whereString, " AND ")
	}
	log.Infof("query:%s params:%v", query, params)

	if tx != nil {
		err = tx.QueryRow(query, params...).Scan(&count)
	} else {
		err = d.db.QueryRow(query, params...).Scan(&count)
	}
	if err != nil {
		if err != sql.ErrNoRows {
			log.Errorf("Error Query %v", err)
			return
		}
	}
	return
}

func (d *ClientUserDao) Get(tx *sql.Tx, in *entity.ClientUser, limit, offset int64) (list []*entity.ClientUser, err error) {
	log.Infof("ClientUserDao.Get %v", in)
	query := `
SELECT
	cu.client_user_id,
	cu.user_id,
	uu.name,
	uu.email,
	uu.firebase_uid,
	cu.delete_flag,
	cu.create_user_id,
	cu.update_user_id,
	cu.created,
	cu.updated
FROM
	client_db.client_users cu
INNER JOIN user_db.users uu ON uu.user_id = cu.user_id
WHERE
	cu.delete_flag = ?
`
	params := []interface{}{
		in.DeleteFlag,
	}
	var whereString []string

	if in.UserID > 0 {
		whereString = append(whereString, " cu.client_user_id = ? ")
		params = append(params, in.UserID)
	}
	if len(in.Name) > 0 {
		whereString = append(whereString, " uu.name LIKE ? ")
		params = append(params, "%"+in.Name+"%")
	}
	if in.FirebaseUID != "" {
		whereString = append(whereString, " uu.firebase_uid = ? ")
		params = append(params, in.FirebaseUID)
	}

	if len(whereString) > 0 {
		query += " WHERE " + strings.Join(whereString, " AND ")
	}
	query += ` ORDER BY uu.user_id DESC `
	query += ` LIMIT ?,?`
	params = append(params, offset)
	params = append(params, limit)

	log.Infof("query:%s params:%v", query, params)
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

	list = []*entity.ClientUser{}
	for rows.Next() {
		row := entity.ClientUser{}
		if e := rows.Scan(
			&row.ClientUserID,
			&row.UserID,
			&row.Name,
			&row.Email,
			&row.FirebaseUID,
			&row.DeleteFlag,
			&row.CreateUserID,
			&row.UpdateUserID,
			&row.Created,
			&row.Updated,
		); e != nil {
			list = []*entity.ClientUser{}
			err = e
			log.Errorf("Error Scan %v", err)
		}
		list = append(list, &row)
	}

	return
}

func (d *ClientUserDao) GetEnable(tx *sql.Tx, in *entity.ClientUser, limit, offset int64) (list []*entity.ClientUser, err error) {
	log.Infof("ClientUserDao.GetEnable %v", in)
	in.DeleteFlag = conf.DeleteFlagOFF
	return d.Get(tx, in, limit, offset)
}

func (d *ClientUserDao) GetDisable(tx *sql.Tx, in *entity.ClientUser, limit, offset int64) (list []*entity.ClientUser, err error) {
	log.Infof("ClientUserDao.GetDisable %v", in)
	in.DeleteFlag = conf.DeleteFlagON
	return d.Get(tx, in, limit, offset)
}
