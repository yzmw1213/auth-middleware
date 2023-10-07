package service

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/yzmw1213/demo-api/conf"

	log "github.com/sirupsen/logrus"

	"github.com/yzmw1213/demo-api/dao"
	"github.com/yzmw1213/demo-api/entity"
	"github.com/yzmw1213/demo-api/util"
)

type ClientUserService struct {
	db            *sql.DB
	clientUserDao *dao.ClientUserDao
	userDao       *dao.UserDao
	authDao       *dao.AuthDao
}

func NewClientUserService() *ClientUserService {
	return &ClientUserService{
		dao.DB(),
		dao.NewClientUserDao(),
		dao.NewUserDao(),
		dao.NewAuthDao(),
	}
}

func (s *ClientUserService) GetDB() *sql.DB {
	return s.db
}

type InputGetClientUser struct {
	ClientUserID int64  `json:"client_user_id" query:"client_user_id"` // 利用ユーザーID
	UserID       int64  `json:"user_id" query:"user_id"`               // ユーザーID
	Name         string `json:"name" query:"name"`                     // 名前
	Email        string `json:"email" query:"email"`                   // Eメール
	FirebaseUID  string `json:"firebase_uid" query:"firebase_uid"`     // firebase id
	Page         int64  `json:"page" query:"page"`                     // ページ番号
	Limit        int64  `json:"limit" query:"limit"`                   // リミット
	Offset       int64  `json:"-"`
}

func (in *InputGetClientUser) GetParam() *InputGetClientUser {
	if in.Limit <= 0 {
		in.Limit = 1000
	}
	if in.Page <= 1 {
		in.Page = 1
	}
	in.Offset = (in.Page - 1) * in.Limit
	return in
}

func (s *ClientUserService) GetClientUser(in *InputGetClientUser) util.OutputBasicInterface {
	log.Info("GetClientUser start %#v", in)
	count, err := s.clientUserDao.GetCount(nil, &entity.ClientUser{
		UserID:      in.UserID,
		Name:        in.Name,
		Email:       in.Email,
		FirebaseUID: in.FirebaseUID,
		DeleteFlag:  false,
	})
	if err != nil {
		log.Errorf("Error clientUserDao.GetCount %v", err)
		return &util.OutputBasic{
			Code:    http.StatusInternalServerError,
			Result:  "NG",
			Message: err,
		}
	}
	list, err := s.clientUserDao.GetEnable(nil, &entity.ClientUser{
		UserID:      in.UserID,
		Name:        in.Name,
		Email:       in.Email,
		FirebaseUID: in.FirebaseUID,
	}, in.Limit, in.Offset)

	if err != nil {
		return &util.OutputBasic{
			Code:    http.StatusInternalServerError,
			Result:  "Error userDao.GetEnable",
			Message: err,
		}
	}
	return util.NewOutputBasicListPaging(
		list,
		count,
		int64(len(list)),
		in.Page,
		in.Limit,
	)
}

type InputSaveClientUser struct {
	ClientUserID int64  `json:"client_user_id" form:"client_user_id"` // 利用ユーザーID
	UserID       int64  `json:"user_id" form:"user_id"`               // ユーザーID
	Name         string `json:"name" form:"name"`                     // 名前
	Email        string `json:"email" form:"email"`                   // Eメール
	Password     string `json:"password" form:"password"`             // パスワード
	UpdateUserID int64  `json:"-" form:"-"`
}

func (s *ClientUserService) Save(in *InputSaveClientUser) util.OutputBasicInterface {
	out := util.ExecTransactionService(s, func(tx *sql.Tx) util.OutputBasicServiceInterface {
		return s.txSave(tx, in)
	})
	return out
}

func (s *ClientUserService) txSave(tx *sql.Tx, in *InputSaveClientUser) util.OutputBasicServiceInterface {
	log.Infof("Save start %#v", in)
	// ユーザー作成済でないかチェックする
	userID, err := s.userDao.GetByEmail(tx, &entity.User{Email: in.Email})
	if err != nil && err != sql.ErrNoRows {
		return &util.OutputBasic{
			Code:    http.StatusInternalServerError,
			Result:  "NG",
			Message: err,
		}
	}
	if userID > 0 {
		in.UserID = userID
	}
	//	in.UserIDなし
	if in.UserID == 0 {
		//	firebase Userを作成
		defaultAuth := []string{conf.CustomUserClaimClient}
		firebaseID, err := s.authDao.CreateUser(in.Name, in.Password, in.Email, defaultAuth)
		if err != nil {
			log.Errorf("Error authDao.CreateUser %v", err)
			return &util.OutputBasic{
				Code:    http.StatusInternalServerError,
				Result:  "NG",
				Message: err,
			}
		}
		//	user作成
		userID, err := s.userDao.Save(tx, &entity.User{
			Name:         in.Name,
			Email:        in.Email,
			FirebaseUID:  firebaseID,
			UpdateUserID: in.UpdateUserID,
		})
		if err != nil {
			log.Errorf("Error userDao.Save %v", err)
			return &util.OutputBasic{
				Code:    http.StatusInternalServerError,
				Result:  "NG",
				Message: err,
			}
		}
		//	client_user作成
		id, err := s.clientUserDao.Save(tx, &entity.ClientUser{
			UserID:       userID,
			UpdateUserID: in.UpdateUserID,
		})
		in.ClientUserID = id
	} else {
		//	Userを更新
		_, err := s.clientUserDao.Save(tx, &entity.ClientUser{
			UserID:       in.UserID,
			UpdateUserID: in.UpdateUserID,
		})
		if err != nil {
			log.Errorf("Error clientUserDao.Save %v", err)
			return &util.OutputBasic{
				Code:    http.StatusInternalServerError,
				Result:  "Error clientUserDao.Save",
				Message: err,
			}
		}
	}
	//	ID指定し、1件しかないことを確認
	list, err := s.clientUserDao.Get(tx, &entity.ClientUser{
		ClientUserID: in.ClientUserID,
	}, 1, 0)
	if err != nil {
		log.Errorf("Error clientUserDao.Save %v", err)
		return &util.OutputBasic{
			Code:    http.StatusInternalServerError,
			Result:  "NG",
			Message: err,
		}
	}
	if len(list) != 1 {
		if len(list) != 1 {
			return &util.OutputBasic{
				Code:    http.StatusInternalServerError,
				Result:  "NG",
				Message: errors.New("no data"),
			}
		}
	}
	return util.NewOutputBasicObject(interface{}(list[0]))
}
