package utils

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type SessionStore struct {
	Store *sessions.CookieStore
}

func NewSessionStore(cfg *SessionConfig) *SessionStore {
	store := SessionStore{

		// Store: sessions.NewCookieStore([]byte("oauth-session-secret"))}
		Store: sessions.NewCookieStore([]byte(cfg.AuthKey))}
	// []byte(cfg.EncKey),
	// store.store.Options = &sessions.Options{}
	return &store
}

func (s *SessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return s.Store.Get(r, name)
}

func (s *SessionStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	return s.Store.Save(r, w, session)
}
