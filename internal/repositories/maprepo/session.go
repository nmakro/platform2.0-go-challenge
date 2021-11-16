package maprepo

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

type Session struct {
	sessionConn *sessions.CookieStore
}

func NewSessionRepo() *Session {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	return &Session{
		sessionConn: store,
	}
}

func (s *Session) SaveSession(r *http.Request) error {
	//token := "adasdased"

	return nil
}

func (s *Session) LoadSession() error {
	return nil
}
