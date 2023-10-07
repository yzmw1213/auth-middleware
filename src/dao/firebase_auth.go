package dao

import (
	"context"
	"errors"

	"github.com/yzmw1213/demo-api/conf"

	"firebase.google.com/go/auth"
	log "github.com/sirupsen/logrus"
)

type AuthDao struct {
}

func NewAuthDao() *AuthDao {
	return &AuthDao{}
}

func getCustomClaims(uid string) (map[string]interface{}, error) {
	ctx := context.Background()
	client, errA := Firebase().Auth(ctx)
	if errA != nil {
		log.Warnf("Warn FB().Auth(ctx) %v", errA)
		return nil, errA
	}
	userRecord, err := client.GetUser(context.Background(), uid)
	if err != nil {
		log.Warnf("Warn client.GetUser(context.Background(), uid) %v", err)
		return nil, err
	}
	// nilを削除する
	result := map[string]interface{}{}
	for claim, value := range userRecord.CustomClaims {
		if value != nil {
			result[claim] = value
		}
	}
	log.Debugf("user GetCustomClaims : Ok")
	return result, nil
}

func (a *AuthDao) CreateUser(name, password, email string, defaultAuthority []string) (firebaseID string, err error) {
	ctx := context.Background()
	client, err := Firebase().Auth(ctx)
	if err != nil {
		return
	}
	u, errCreate := client.GetUserByEmail(ctx, email)
	if errCreate != nil {
		log.Infof("no User data %v", errCreate)
		params := (&auth.UserToCreate{}).Email(email).EmailVerified(false).DisplayName(name).Password(password)
		user, createErr := client.CreateUser(ctx, params)
		if createErr != nil {
			err = createErr
			log.Errorf("Error client.CreateUser %v", createErr)
			return
		}
		firebaseID = user.UID
	} else {
		log.Infof("Update User data %v", u)
		firebaseID = u.UID
		if u.DisplayName != name {
			params := (&auth.UserToUpdate{}).
				DisplayName(name)
			u, errU := client.UpdateUser(ctx, u.UID, params)
			if errU != nil {
				log.Errorf("error update user: %v", errU)
				err = errU
				return
			}
			log.Infof("Successfully update user: %v", u)
		}
	}
	_, err = a.AddCustomClaim(firebaseID, defaultAuthority)
	return
}

func (a *AuthDao) AddCustomClaim(uid string, keys []string) (map[string]interface{}, error) {
	if len(keys) < 1 {
		return nil, errors.New("CustomClaims Add keys 0")
	}
	claims, err := getCustomClaims(uid)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	client, errA := Firebase().Auth(ctx)
	if errA != nil {
		return nil, errA
	}
	// auth 追加
	for _, key := range keys {
		if key == "" {
			continue
		}
		claims[key] = conf.CustomUserClaimON
	}
	err = client.SetCustomUserClaims(context.Background(), uid, claims)
	if err != nil {
		return nil, err
	}
	log.Infof("user AddCustomClaims : %v", keys)
	return claims, err
}
