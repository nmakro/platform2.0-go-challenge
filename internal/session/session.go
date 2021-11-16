package session

import (
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
)

var sessionStore = sessions.NewCookieStore([]byte(viper.GetString("SESSION_AUTH_KEY")))

func GetSessionStore() *sessions.CookieStore {
	return sessionStore
}
