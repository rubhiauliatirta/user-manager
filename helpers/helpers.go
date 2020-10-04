package helpers

import (
	"errors"
	"net/http"
	"strconv"

	"../config"
	"../models"
	"github.com/gorilla/sessions"
)

func GetUser(s *sessions.Session) (models.AuthUser, error) {
	val := s.Values["user"]
	var user = models.AuthUser{}
	user, ok := val.(models.AuthUser)
	if !ok {
		return user, errors.New("Authentication Failed")
	}
	return user, nil
}

func Auth(r *http.Request) (models.AuthUser, error) {
	var authenticatedUser models.AuthUser
	session, err := config.Store.Get(r, "auth")
	if err != nil {
		return authenticatedUser, err
	}

	authUser, err := GetUser(session)
	if err != nil || authUser.IsAuthenticated == false {
		return authenticatedUser, err
	}

	return authUser, nil

}

func IsAuthorized(user models.AuthUser, r *http.Request) bool {
	r.ParseForm()
	editedId := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(editedId)

	if user.User.Role == "admin" || user.User.Id == id {
		return true
	}
	return false
}
