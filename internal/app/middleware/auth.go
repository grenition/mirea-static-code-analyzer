package middleware

import (
	"net/http"
	"github.com/gorilla/sessions"
)

const sessionName = "auth_session"
const userIDKey = "user_id"
const usernameKey = "username"

func UserIDKey() string {
	return userIDKey
}

func UsernameKey() string {
	return usernameKey
}

var store *sessions.CookieStore

func InitAuth(secretKey string) {
	store = sessions.NewCookieStore([]byte(secretKey))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	}
}

func GetSession(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, sessionName)
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := GetSession(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		userID := session.Values[userIDKey]
		if userID == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next(w, r)
	}
}

func GetUserID(r *http.Request) (string, bool) {
	session, err := GetSession(r)
	if err != nil {
		return "", false
	}
	userID, ok := session.Values[userIDKey].(string)
	return userID, ok
}

func GetUsername(r *http.Request) (string, bool) {
	session, err := GetSession(r)
	if err != nil {
		return "", false
	}
	username, ok := session.Values[usernameKey].(string)
	return username, ok
}

func GetUserIDFromRequest(r *http.Request) (string, bool) {
	return GetUserID(r)
}

