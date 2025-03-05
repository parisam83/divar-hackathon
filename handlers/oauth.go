package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/services"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
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
}

func NewOAuthHandler(serv *services.OAuthService) *oAuthHandler {

	return &oAuthHandler{
		oauthService: serv,
	}
}

type OAuthSession struct {
	PostToken   string `json:"post_token"`
	State       string `json:"state"`
	CallbackURL string `json:"callback_url"`
	SessionKey  string `json:"session_key"`
}

var store = sessions.NewCookieStore([]byte("oauth-session-secret"))

func (h *oAuthHandler) AddonOauth(w http.ResponseWriter, r *http.Request) {
	log.Println("AddonOauth called")

	post_token := r.URL.Query().Get("post_token")
	callback_url := r.URL.Query().Get("return_url")

	if post_token == "" || callback_url == "" {
		http.Error(w, "post_token and return_url are required", http.StatusBadRequest)
		return
	}
	oauthSession, err := store.Get(r, "oauth-session")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//check if session existed before....using database
	//TODO

	// if there was no session appointed to the user and post, create a new session and state
	state := uuid.New().String()
	session := &OAuthSession{
		PostToken:   post_token,
		State:       state,
		CallbackURL: callback_url,
	}
	sessionJson, err := json.Marshal(session)
	if err != nil {
		http.Error(w, "Failed to marshal session: "+err.Error(), http.StatusInternalServerError)
		return
	}
	oauthSession.Values["data"] = sessionJson
	err = oauthSession.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to set session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// choosing scopes
	oauthScopes := []Scope{
		{resourceType: POST_ADDON_CREATE, resourceID: post_token},
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

	redirect_url := h.oauthService.GenerateAuthURL(scopes, state)
	log.Println(redirect_url)
	http.Redirect(w, r, redirect_url, http.StatusFound)

}

func (h *oAuthHandler) OauthCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("OauthCallback called")
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	oauthSession, err := store.Get(r, "oauth-session")
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
	h.oauthService.AddOAuth(session.SessionKey, accessToken, refreshToken, session.PostToken, expires_in)

	url := fmt.Sprintf("https://oryx-meet-elf.ngrok-free.app/poi")
	http.Redirect(w, r, url, http.StatusSeeOther)

	//redirect to api for poi
}
