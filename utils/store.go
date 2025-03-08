package utils

import "github.com/gorilla/sessions"

type sessionStore struct {
	store *sessions.CookieStore
}

// func NewSessionStore() {

// 	sessionStore{
// 		store: sessions.NewCookieStore(
// 			cfg.
// 		),
// 	}
// }
// func GetS
