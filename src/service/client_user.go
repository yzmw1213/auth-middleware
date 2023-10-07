package service

import (
	"database/sql"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/yzmw1213/demo-api/dao"
	"github.com/yzmw1213/demo-api/entity"
	"github.com/yzmw1213/demo-api/util"
)

type ClientUserService struct {
	db            *sql.DB
	clientUserDao *dao.ClientUserDao
}

func NewClientUserService() *ClientUserService {
	return &ClientUserService{
		dao.DB(),
		dao.NewClientUserDao(),
	}
}

func (s *ClientUserService) GetDB() *sql.DB {
	return s.db
}

type InputSaveUser struct {
	UserID      int64  `json:"user_id" query:"user_id"`           // ユーザーID
	Name        string `json:"name" query:"name"`                 // 名前
	Email       string `json:"email" query:"email"`               // Eメール
	FirebaseUID string `json:"firebase_uid" query:"firebase_uid"` // firebase id
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

//func (s *ClientUserService) Save(in *InputSaveUser) util.OutputBasicInterface {
//	out := util.ExecTransactionService(s, func(tx *sql.Tx) util.OutputBasicServiceInterface {
//		return s.txSave(tx, in)
//	})
//	return out
//}
//
//func (s *ClientUserService) txSave(tx *sql.Tx, in *InputSaveUser) util.OutputBasicServiceInterface {
//	log.Infof("Save start %#v", in)
//
//	//	in.UserIDなし
//	if in.UserID == 0 {
//		//	firebase Userを作成
//		//	Claim clientを追加
//		//	user作成
//		//	client_user作成
//	} else {
//		//	Userを更新
//	}
//	//	ID指定し、1件しかないことを確認
//
//}
