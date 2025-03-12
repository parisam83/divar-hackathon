package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/configs"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

const (
	SessionName = "auth_session"
	SessionKey  = "data"
)

type SessionStore struct {
	Store *sessions.CookieStore
}

type OAuthSession struct {
	PostToken   string `json:"post_token"`
	State       string `json:"state"`
	CallbackURL string `json:"callback_url"`
	SessionKey  string `json:"session_key"`
}

func NewSessionStore(cfg *configs.SessionConfig) *SessionStore {
	store := SessionStore{
		Store: sessions.NewCookieStore([]byte(cfg.AuthKey)),
	}
	return &store
}

func (s *SessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return s.Store.Get(r, name)
}

func (s *SessionStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	return s.Store.Save(r, w, session)
}

func (h *SessionStore) GetExistingSession(w http.ResponseWriter, r *http.Request) (*OAuthSession, error) {
	oauthSession, err := h.Get(r, SessionName)
	if err != nil {
		return nil, fmt.Errorf("%s", "Failed to get session: "+err.Error())
	}

	data, ok := oauthSession.Values[SessionKey].([]byte)
	if !ok {
		return nil, fmt.Errorf("%s", "No session data found ")
	}

	session := &OAuthSession{}
	if err := json.Unmarshal(data, session); err != nil {
		return nil, fmt.Errorf("%s", "failed to decode session:"+err.Error())
	}
	return session, nil
}

func (h *SessionStore) SaveSession(w http.ResponseWriter, r *http.Request, session *OAuthSession) (*OAuthSession, error) {
	oauthSession, err := h.Get(r, SessionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	sessionJson, err := json.Marshal(session)
	if err != nil {
		return nil, fmt.Errorf("%s", "Failed to marshal session: "+err.Error())
	}

	oauthSession.Values[SessionKey] = sessionJson
	err = h.Save(r, w, oauthSession)
	if err != nil {
		return nil, fmt.Errorf("%s", "Failed to save session: "+err.Error())
	}
	log.Printf("Session saved successfully with state: %s", session.State)
	return session, nil

}

func (h *SessionStore) CreateNewSession(w http.ResponseWriter, r *http.Request, postToken string) (*OAuthSession, error) {
	state := uuid.New().String()
	session := &OAuthSession{
		PostToken: postToken,
		State:     state,
	}
	return h.SaveSession(w, r, session)
}
