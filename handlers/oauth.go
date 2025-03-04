package handlers

import (
	"Realestate-POI/services"
	"fmt"
	"log"
	"net/http"

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

var store = sessions.NewCookieStore([]byte("oauth-session-secret"))

func (h *oAuthHandler) AddonOauth(w http.ResponseWriter, r *http.Request) {

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
	state := uuid.New().String()

	// if there was no session appointed to the user and post, create a new session
	oauthSession.Values["post_token"] = post_token
	oauthSession.Values["callback_url"] = callback_url
	oauthSession.Values["state"] = state
	err = oauthSession.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to set session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	oauthScopes := []Scope{
		{resourceType: POST_ADDON_CREATE, resourceID: post_token},
		{resourceType: USER_PHONE},
	}
	var scopes []string

	for _, scope := range oauthScopes {
		if scope.resourceID != "" {
			scopes = append(scopes, fmt.Sprintf("%s.%s", scope.resourceType, scope.resourceID))
		} else {
			scopes = append(scopes, string(scope.resourceType))
		}
	}
	// create a post with token in database????????/

	redirect_url := h.oauthService.GenerateAuthURL(scopes, state)
	log.Println(redirect_url)
	http.Redirect(w, r, redirect_url, http.StatusFound)
}

func (h *oAuthHandler) OauthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" || state == "" {
		http.Error(w, "code and state are required", http.StatusBadRequest)
		return
	}

	oauthSession, err := store.Get(r, "oauth-session")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//check if state is the same as the one in the session
	sessionState, ok := oauthSession.Values["state"].(string)
	if !ok || sessionState != state {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}
	// deleting state from session because we dont need it after oauth
	delete(oauthSession.Values, "state")
	err = oauthSession.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session after deleting state: "+err.Error(), http.StatusInternalServerError)
		return
	}

	//sending code to get the token
	token, err := h.oauthService.ExchangeToken(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	accessToken := token.AccessToken
	expires_in := token.Expiry
	log.Println(accessToken)
	log.Println(expires_in)

	//create the oauth object with session key , tokens and post in database
	sessionKey := uuid.New().String()
	oauthSession.Values["sessionKey"] = sessionKey
	oauthSession.Save(r, w)

	//redirect to api for poi

}
