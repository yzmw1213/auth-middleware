package service

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/yzmw1213/demo-api/dao"
	"github.com/yzmw1213/demo-api/entity"
	"github.com/yzmw1213/demo-api/util"
)

type UserService struct {
	userDao *dao.UserDao
}

func NewUserService() *UserService {
	return &UserService{
		dao.NewUserDao(),
	}
}

type InputGetUser struct {
	UserID      int64  `json:"user_id" query:"user_id"`           // ユーザーID
	Name        string `json:"name" query:"name"`                 // 名前
	Email       string `json:"email" query:"email"`               // Eメール
	FirebaseUID string `json:"firebase_uid" query:"firebase_uid"` // firebase id
	Page        int64  `json:"page" query:"page"`                 // ページ番号
	Limit       int64  `json:"limit" query:"limit"`               // リミット
	Offset      int64  `json:"-"`
}

func (in *InputGetUser) GetParam() *InputGetUser {
	if in.Limit <= 0 {
		in.Limit = 1000
	}
	if in.Page <= 1 {
		in.Page = 1
	}
	in.Offset = (in.Page - 1) * in.Limit
	return in
}

func (s *UserService) GetUser(in *InputGetUser) util.OutputBasicInterface {
	log.Info("GetUser start %#v", in)
	count, err := s.userDao.GetCount(nil, &entity.User{
		UserID:      in.UserID,
		Name:        in.Name,
		Email:       in.Email,
		FirebaseUID: in.FirebaseUID,
		DeleteFlag:  false,
	})
	if err != nil {
		log.Errorf("Error userDao.GetCount %v", err)
		return &util.OutputBasic{
			Code:    http.StatusInternalServerError,
			Result:  "NG",
			Message: err,
		}
	}
	list, err := s.userDao.GetEnable(nil, &entity.User{
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
