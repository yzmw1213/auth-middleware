package dao

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/yzmw1213/demo-api/conf"
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

func (d *UserDao) GetCount(tx *sql.Tx, in *entity.User) (count int64, err error) {
	log.Infof("UserDao.GetCount %v", in)
	query := `
SELECT
	count(uu.user_id)
FROM
	user_db.users uu
WHERE
	uu.delete_flag = ?
`
	params, whereString := d.getWhere(in)

	if len(whereString) > 0 {
		query += " AND " + strings.Join(whereString, " AND ")
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

func (d *UserDao) Get(tx *sql.Tx, in *entity.User, limit, offset int64) (list []*entity.User, err error) {
	log.Infof("UserDao.Get %v", in)
	query := `
SELECT
	uu.user_id,
	uu.name,
	uu.email,
	uu.firebase_uid,
	uu.delete_flag,
	uu.created,
	uu.updated
FROM
	user_db.users uu
WHERE
	uu.delete_flag = ?
`
	params, whereString := d.getWhere(in)

	if len(whereString) > 0 {
		query += " AND " + strings.Join(whereString, " AND ")
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

	list = []*entity.User{}
	for rows.Next() {
		row := entity.User{}
		if e := rows.Scan(
			&row.UserID,
			&row.Name,
			&row.Email,
			&row.FirebaseUID,
			&row.DeleteFlag,
			&row.Created,
			&row.Updated,
		); e != nil {
			list = []*entity.User{}
			err = e
			log.Errorf("Error Scan %v", err)
		}
		list = append(list, &row)
	}

	return
}

func (d *UserDao) GetByEmail(tx *sql.Tx, in *entity.User) (userId int64, err error) {
	log.Infof("UserDao.GetByEmail %v", in)
	query := `
SELECT
	user_id
FROM
	user_db.users
WHERE
	email = ? 
`
	params := []interface{}{
		in.Email,
	}
	log.Infof("query:%s params:%v", query, params)
	if tx != nil {
		err = tx.QueryRow(query, params...).Scan(&userId)
	} else {
		err = d.db.QueryRow(query, params...).Scan(&userId)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Errorf("Error Query %v", err)
		}
		return
	}
	return
}

func (d *UserDao) getWhere(in *entity.User) (params []interface{}, whereString []string) {
	params = append(params, in.DeleteFlag)
	if in.UserID > 0 {
		whereString = append(whereString, " uu.user_id = ? ")
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
	return
}

func (d *UserDao) GetEnable(tx *sql.Tx, in *entity.User, limit, offset int64) (list []*entity.User, err error) {
	log.Infof("UserDao.GetEnable %v", in)
	in.DeleteFlag = conf.DeleteFlagOFF
	return d.Get(tx, in, limit, offset)
}

func (d *UserDao) GetDisable(tx *sql.Tx, in *entity.User, limit, offset int64) (list []*entity.User, err error) {
	log.Infof("UserDao.GetDisable %v", in)
	in.DeleteFlag = conf.DeleteFlagON
	return d.Get(tx, in, limit, offset)
}

func (d *UserDao) Save(tx *sql.Tx, in *entity.User) (insertedId int64, err error) {
	log.Infof("UserDao.Save %v", in)

	query := `
INSERT INTO user_db.users
(
	name,
	email,
	firebase_uid,
	update_user_id
) VALUES (?,?,?,?)
`
	params := []interface{}{
		in.Name,
		in.Email,
		in.FirebaseUID,
		in.UpdateUserID,
	}
	log.Infof("query: %v param: %v", query, params)

	var result sql.Result

	if tx != nil {
		result, err = tx.Exec(query, params...)
	} else {
		result, err = d.db.Exec(query, params...)
	}
	if err != nil {
		log.Errorf("error query %v", err)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Errorf("error LastInsertId %v", err)
		return
	}
	log.Infof("Inserted ID user %d", id)
	insertedId = id
	return
}
