package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

func SetAuthCookie(w http.ResponseWriter, r *http.Request, token string) {
	cookie := &http.Cookie{
		Name:   "token",
		Value:  token,
		Path:   "/",
		MaxAge: 60 * 60 * 24 * 30,
	}
	http.SetCookie(w, cookie)
}

func ExtendAuthCookie(w http.ResponseWriter, r *http.Request, token string) {
	cookie := &http.Cookie{
		Name:   "token",
		Value:  token,
		Path:   "/",
		MaxAge: 60 * 60 * 24 * 30,
	}
	http.SetCookie(w, cookie)
}

func GetAuthCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func DeleteAuthCookie(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func HashToken(token string) string {
	sha256sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sha256sum[:])
}
