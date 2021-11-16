package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/nmakro/platform2.0-go-challenge/internal/app"
	"github.com/nmakro/platform2.0-go-challenge/internal/app/user"
	gwihttp "github.com/nmakro/platform2.0-go-challenge/pkg/http"
	"github.com/nmakro/platform2.0-go-challenge/pkg/security"
)

var store = sessions.NewCookieStore([]byte("!@saz#asdasd@DaazSq@3"))

type UsersModule struct {
	service *user.UserService
}

func Setup(router *mux.Router, service *user.UserService) {
	m := &UsersModule{
		service: service,
	}

	router.HandleFunc("/signup", m.SignUp).Methods("POST")
	router.HandleFunc("/login", m.SignIn).Methods("POST")

	users := router.PathPrefix("/users").Subrouter()
	users.HandleFunc("/", m.ListUsers).Methods("GET")
	users.HandleFunc("/user/{id}", m.GetUser).Methods("GET")
	users.HandleFunc("/user/{id}", m.DeleteUser).Methods("DELETE")
}

func (m *UsersModule) SignUp(w http.ResponseWriter, r *http.Request) {
	req := user.AddUserCommand{}
	if err := gwihttp.ValidateRequest(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := m.service.AddUser(r.Context(), req); err != nil {
		var notValid *user.ErrValidation
		if errors.As(err, &notValid) {
			gwihttp.ResponseWithJSON(http.StatusUnprocessableEntity, map[string]interface{}{"error": err.Error()}, w)
			return
		}

		var duplicateUser *app.ErrDuplicateEntry
		if errors.As(err, &duplicateUser) {
			gwihttp.ResponseWithJSON(http.StatusConflict, map[string]interface{}{"error": err.Error()}, w)
			return
		}

		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	session, _ := store.Get(r, "gwi-cookie")
	session.Values["authenticated"] = true
	session.Values["user_email"] = req.Email

	if err := session.Save(r, w); err != nil {
		errSession := fmt.Errorf("error in user sign: %w", err)
		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": errSession.Error()}, w)
		return
	}
	fmt.Println(session.Values)

	gwihttp.ResponseWithJSON(http.StatusOK, nil, w)
}

func (m *UsersModule) SignIn(w http.ResponseWriter, r *http.Request) {
	req := user.AddUserCommand{}
	if err := gwihttp.ValidateRequest(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := m.service.GetUserWithPassword(r.Context(), req.Email)
	if err != nil {
		var notFound *app.ErrEntityNotFound
		if errors.As(err, &notFound) {
			gwihttp.ResponseWithJSON(http.StatusNotFound, map[string]interface{}{"error": err.Error()}, w)
			return
		}

		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	if !security.CheckPasswordHash(req.Password, user.Password) {
		gwihttp.ResponseWithJSON(http.StatusUnauthorized, map[string]interface{}{"error": "wrong password"}, w)
		return
	}

	session, _ := store.Get(r, "gwi-cookie")
	session.Values["authenticated"] = true
	session.Values["user_email"] = user.Email
	if err = session.Save(r, w); err != nil {
		errSession := fmt.Errorf("error in user sign: %w", err)
		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": errSession.Error()}, w)
		return
	}

	gwihttp.ResponseWithJSON(http.StatusOK, nil, w)
}

func (m *UsersModule) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userEmail, ok := vars["id"]
	if !ok {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "you must specify a user email"}, w)
		return
	}

	user, err := m.service.GetUser(r.Context(), userEmail)
	if err != nil {
		var notFound *app.ErrEntityNotFound
		if errors.As(err, &notFound) {
			gwihttp.ResponseWithJSON(http.StatusNotFound, map[string]interface{}{"error": err.Error()}, w)
			return
		}

		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return
	}

	gwihttp.ResponseWithJSON(http.StatusOK, user, w)
}

func (m *UsersModule) DeleteUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "gwi-cookie")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		gwihttp.ResponseWithJSON(http.StatusUnauthorized, nil, w)
		return
	}

	vars := mux.Vars(r)
	userEmail, ok := vars["id"]
	if !ok {
		gwihttp.ResponseWithJSON(http.StatusBadRequest, map[string]interface{}{"error": "you must specify a user email"}, w)
		return
	}

	if sessionMail, ok := session.Values["user_email"].(string); !ok || sessionMail != userEmail {
		gwihttp.ResponseWithJSON(http.StatusUnauthorized, map[string]interface{}{"error": "user cannot delele another's user account"}, w)
		return
	}

	err := m.service.DeleteUser(r.Context(), userEmail)
	if err != nil {
		var notFound *app.ErrEntityNotFound
		if !errors.As(err, &notFound) {
			gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
			return
		}
	}

	gwihttp.ResponseWithJSON(http.StatusNoContent, nil, w)
}

func (m *UsersModule) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := m.service.ListUsers(r.Context())
	if err != nil {
		gwihttp.ResponseWithJSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()}, w)
		return

	}
	gwihttp.ResponseWithJSON(http.StatusOK, users, w)
}
