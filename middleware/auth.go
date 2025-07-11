package middleware

import (
	"context"
	"net/http"
	"subalertor/logger"
	"subalertor/storage"
	"subalertor/utils"
	"time"
)

type UserCtxType string

const UserCtxKey UserCtxType = "user"

func AuthMiddleware(store *storage.Storage) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := utils.GetAuthCookie(r)
			if err != nil {
				logger.Log.Error(err)
				utils.Unauthorized(w)
				return
			}

			session, err := store.SessionStore.FindSession(utils.HashToken(cookie))
			if err != nil {
				logger.Log.Error(err)
				utils.Unauthorized(w)
				return
			}
			if time.Now().Unix() > int64(session.ExpiresAt) {
				logger.Log.Error(err)
				utils.Unauthorized(w)
				return
			}

			if session.ExpiresAt-int(time.Now().Unix()) <= 60*60*24*15 {
				_, err := store.SessionStore.ExtendSession(utils.HashToken(cookie), time.Now().Add(time.Hour*24*15))
				if err != nil {
					logger.Log.Error(err)
				} else {
					utils.ExtendAuthCookie(w, r, cookie)
				}
			}

			user, err := store.UserStore.FindUserByID(session.UserID)
			if err != nil {
				logger.Log.Error(err)
				utils.Unauthorized(w)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserCtxKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) storage.UserModel {
	return ctx.Value(UserCtxKey).(storage.UserModel)
}
