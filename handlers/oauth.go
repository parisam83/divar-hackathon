package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/services"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/utils"

	"github.com/google/uuid"
)

type OauthResourceType string

const (
	POST_ADDON_CREATE OauthResourceType = "POST_ADDON_CREATE"
	USER_PHONE        OauthResourceType = "USER_PHONE"
	OFFLINE_ACCESS    OauthResourceType = "offline_access"
)

type Scope struct {
	resourceType OauthResourceType
	resourceID   string
}

type oAuthHandler struct {
	oauthService *services.OAuthService
	store        *utils.SessionStore
}

func NewOAuthHandler(store *utils.SessionStore, serv *services.OAuthService) *oAuthHandler {
	return &oAuthHandler{
		oauthService: serv,
		store:        store,
	}
}

type OAuthSession struct {
	PostToken   string `json:"post_token"`
	State       string `json:"state"`
	CallbackURL string `json:"callback_url"`
	SessionKey  string `json:"session_key"`
}

const (
	SessionName = "Auth"
	SessionKey  = "data"
)

func (h *oAuthHandler) getExistingSession(w http.ResponseWriter, r *http.Request) (*OAuthSession, error) {
	oauthSession, err := h.store.Get(r, SessionName)
	if err != nil {
		return nil, fmt.Errorf("%s", "Failed to get session: "+err.Error())
	}

	data, ok := oauthSession.Values[SessionKey].([]byte)
	if !ok {
		return nil, fmt.Errorf("%s", "No session data found ")

	}

	var session *OAuthSession
	if err := json.Unmarshal(data, session); err != nil {
		return nil, fmt.Errorf("%s", "failed to decode session:"+err.Error())
	}
	return session, nil
}

func (h *oAuthHandler) createNewSession(w http.ResponseWriter, r *http.Request, postToken string) (*OAuthSession, error) {
	oauthSession, err := h.store.Get(r, SessionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	state := uuid.New().String()
	session := &OAuthSession{
		PostToken: postToken,
		State:     state,
	}
	sessionJson, err := json.Marshal(session)
	if err != nil {
		return nil, fmt.Errorf("%s", "Failed to marshal session: " + err.Error())
	}

	oauthSession.Values["data"] = sessionJson
	err = h.store.Save(r, w, oauthSession)
	if err != nil {
		return nil, fmt.Errorf("%s", "Failed to set session: " + err.Error())
	}
	return session, nil

}

func (h *oAuthHandler) AddonOauth(w http.ResponseWriter, r *http.Request) {
	log.Println("AddonOauth called")

	postToken := r.URL.Query().Get("post_token")
	callback_url := r.URL.Query().Get("return_url")

	if postToken == "" || callback_url == "" {
		http.Error(w, "post_token and return_url are required", http.StatusBadRequest)
		return
	}
	// check existing session
	session, err := h.getExistingSession(w, r)
	if err == nil && session != nil {
		//redirect
	}

	//create new session
	session, err = h.createNewSession(w, r, postToken)
	if err != nil {
		http.Error(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// choosing scopes
	oauthScopes := []Scope{
		{resourceType: POST_ADDON_CREATE, resourceID: postToken},
		{resourceType: USER_PHONE},
		// {resourceType: OFFLINE_ACCESS},
	}
	var scopes []string

	for _, scope := range oauthScopes {
		if scope.resourceID != "" {
			scopes = append(scopes, fmt.Sprintf("%s.%s", scope.resourceType, scope.resourceID))
		} else {
			scopes = append(scopes, string(scope.resourceType))
		}
	}

	redirect_url := h.oauthService.GenerateAuthURL(scopes, session.State)
	log.Println(redirect_url)
	http.Redirect(w, r, redirect_url, http.StatusFound)

}

func (h *oAuthHandler) OauthCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("OauthCallback called")
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	oauthSession, err := h.store.Get(r, "Session")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// if code or state is not present, we redirect to the add page
	//can we redirect to the user's add page?
	if code == "" || state == "" {
		// http.Error(w, "code and state are required", http.StatusBadRequest)
		// post_token := oauthSession.Values["post_token"].(string)
		redirectURL := fmt.Sprintf("https://divar.ir/")
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
	data, ok := oauthSession.Values["data"].([]byte)
	if !ok {
		http.Error(w, "no session data found", http.StatusBadRequest)
		return

	}

	var session OAuthSession
	if err := json.Unmarshal(data, &session); err != nil {
		http.Error(w, "failed to decode session:"+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("session state is %s", session.State)

	//check if state is the same as the one in the session
	if session.State != state {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}
	//sending code to get the token
	token, err := h.oauthService.ExchangeToken(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	accessToken := token.AccessToken
	expires_in := time.Unix(token.Expiry.Unix(), 0)
	refreshToken := token.RefreshToken
	log.Println(accessToken)
	log.Println(expires_in)

	// deleting state from session because we dont need it after oauth
	session.State = ""
	//add sessionKey to the reuqest
	session.SessionKey = uuid.New().String()
	//update
	updatedData, err := json.Marshal(session)
	if err != nil {
		log.Println("Session encoding failed:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	oauthSession.Values["data"] = updatedData
	if err := oauthSession.Save(r, w); err != nil {
		log.Println("Session save failed:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.oauthService.InsertOAuthData(session.SessionKey, accessToken, refreshToken, session.PostToken, expires_in)

	// url := fmt.Sprintf("https://oryx-meet-elf.ngrok-free.app/poi")
	// http.Redirect(w, r, url, http.StatusSeeOther)

	//redirect to api for poi
}
