package middleware

import (
	"context"
	"net/http"
	"time"
	"twithoauth/storage"
	"twithoauth/utils"
)

type UserCtxType string

const UserCtxKey UserCtxType = "user"

func AuthMiddleware(store *storage.Storage) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := utils.GetAuthCookie(r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			session, err := store.SessionStore.FindSession(utils.HashToken(cookie))
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if time.Now().Unix() > int64(session.ExpiresAt) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// extend session if time more then 15 days
			if time.Now().Unix()-int64(session.ExpiresAt) > 60*60*24*15 {
				_, err := store.SessionStore.ExtendSession(utils.HashToken(cookie), time.Now().Add(time.Hour*24*15))
				if err != nil {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}
			user, err := store.UserStore.FindUser(cookie)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			w.Header().Set("X-User-Id", user.ID)
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserCtxKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) storage.UserModel {
	return ctx.Value(UserCtxKey).(storage.UserModel)
}
