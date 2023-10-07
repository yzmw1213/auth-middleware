package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/yzmw1213/demo-api/conf"
	"github.com/yzmw1213/demo-api/dao"
	"github.com/yzmw1213/demo-api/entity"
)

func AuthAPI(next gin.HandlerFunc, authorities []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		auth, err := dao.Firebase().Auth(ctx)
		if err != nil {
			log.Errorf("Error Auth %v", err)
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "NG",
			})
			c.Abort()
			return
		}
		//	JWT取得
		authorizationHeader := c.Request.Header.Get("Authorization")
		idToken := strings.Replace(authorizationHeader, "Bearer ", "", 1)

		//	idTokenからid取得し、ユーザーの権限をチェックする
		token, err := auth.VerifyIDToken(ctx, idToken)
		if err != nil {
			log.Warnf("Error VerifyIDToken %v", err)
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": "NG",
			})
			c.Abort()
			return
		}

		uid := token.Claims["user_id"]
		if uid.(string) == "" {
			log.Warnf("Error ID empty")
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": "NG",
			})
			c.Abort()
			return
		}

		userList, err := dao.NewUserDao().GetEnable(nil, &entity.User{
			FirebaseUID: uid.(string),
		}, 1, 0)
		if err != nil {
			log.Errorf("Error User Get %v", err)
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": "NG",
			})
			c.Abort()
			return
		}
		if len(userList) != 1 {
			log.Warn("No User")
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": "NG",
			})
			c.Abort()
			return
		}
		user := userList[0]
		c.Set("user_name", user.Name)
		c.Set("user_id", user.UserID)
		c.Set("firebase_id", user.FirebaseUID)

		var isValue int64 = 1
		if len(authorities) > 0 {
			//	権限チェックが必要
			authCheck := false
			for _, authority := range authorities {
				if val, ok := token.Claims[authority]; ok {
					switch val.(type) {
					case int:
						if val.(int) == conf.CustomUserClaimON {
							authCheck = true
							if val == conf.CustomUserClaimAdmin {
								c.Set("is_admin", isValue) // adminの設定あり
							} else if val == conf.CustomUserClaimClient {
								c.Set("is_client", isValue) // clientの設定あり
							}
						}
					case float64:
						if val.(float64) == conf.CustomUserClaimON {
							authCheck = true
							if val == conf.CustomUserClaimAdmin {
								c.Set("is_admin", isValue) // adminの設定あり
							} else if val == conf.CustomUserClaimClient {
								c.Set("is_client", isValue) // clientの設定あり
							}
						}
					}

				}
			}
			if !authCheck {
				message := "error auth check"
				log.Warnf(message)
				c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"code":    http.StatusUnauthorized,
					"message": message,
				})
				c.Abort()
				return
			}
		}

		next(c)
	}

}
