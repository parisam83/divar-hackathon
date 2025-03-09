package utils

import (
	"github.com/gorilla/sessions"
)

type sessionStore struct {
	store *sessions.CookieStore
}

func NewSessionStore(cfg *SessionConfig) *sessionStore {
	store := sessionStore{
		store: sessions.NewCookieStore(
			[]byte(cfg.AuthKey),
			[]byte(cfg.EncKey),
		),
	}
	return &store
}
