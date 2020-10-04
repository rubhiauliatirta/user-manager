package config

import (
	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func init() {
	// authKeyOne := securecookie.GenerateRandomKey(64)
	// encryptionKeyOne := securecookie.GenerateRandomKey(32)

	Store = sessions.NewCookieStore([]byte("rubhiauliatirtarubhiauliatirta"))

	Store.Options = &sessions.Options{
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

}
