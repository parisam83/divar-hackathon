package handlers

import (
	"fmt"
	"log"
	"net/http"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/internal/services"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/utils"
)

type OauthResourceType string

const (
	POST_ADDON_CREATE  OauthResourceType = "POST_ADDON_CREATE"
	USER_PHONE         OauthResourceType = "USER_PHONE"
	USER_ID            OauthResourceType = "USER_ID"
	OFFLINE_ACCESS     OauthResourceType = "offline_access"
	DefaultRedirectURL                   = "https://divar.ir/"
)

type Scope struct {
	resourceType OauthResourceType
	resourceID   string
}

type oAuthHandler struct {
	oauthService *services.OAuthService
	kenarService *services.KenarService
	store        *utils.SessionStore
}

func NewOAuthHandler(store *utils.SessionStore, serv *services.OAuthService, kenar *services.KenarService) *oAuthHandler {
	if store == nil {
		log.Fatal("cookie store can not be nil")
	}
	if serv == nil {
		log.Fatal("oauth service can not be nil")
	}

	return &oAuthHandler{
		store:        store,
		oauthService: serv,
		kenarService: kenar,
	}
}

func (h *oAuthHandler) buildScopes(postToken string) []string {
	oauthScopes := []Scope{
		{resourceType: POST_ADDON_CREATE, resourceID: postToken},
		{resourceType: USER_ID},
		{resourceType: OFFLINE_ACCESS},
	}

	var scopes []string
	for _, scope := range oauthScopes {
		if scope.resourceID != "" {
			scopes = append(scopes, fmt.Sprintf("%s.%s", scope.resourceType, scope.resourceID))
		} else {
			scopes = append(scopes, string(scope.resourceType))
		}
	}
	return scopes
}

func (h *oAuthHandler) AddonOauth(w http.ResponseWriter, r *http.Request) {
	log.Println("AddonOauth called")

	postToken := r.URL.Query().Get("post_token")
	callback_url := r.URL.Query().Get("return_url")

	if postToken == "" || callback_url == "" {
		http.Error(w, "post_token and return_url are required", http.StatusBadRequest)
		return
	}
	session, err := h.store.GetExistingSession(w, r)
	if err != nil || session.PostToken != postToken {
		session, err = h.store.CreateNewSession(w, r, postToken)
		if err != nil {
			http.Error(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// creating scopes
	scopes := h.buildScopes(postToken)
	redirect_url := h.oauthService.GenerateAuthURL(scopes, session.State)
	http.Redirect(w, r, redirect_url, http.StatusFound)

}

// call back
func (h *oAuthHandler) validateCallback(r *http.Request) (string, string, error) {

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		return "", "", fmt.Errorf("code and state are required")
	}
	return code, state, nil
}

func (h *oAuthHandler) OauthCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("OauthCallback called")
	code, state, err := h.validateCallback(r)
	if err != nil {
		log.Printf("Invalid callback parameters: %v", err)
		http.Redirect(w, r, DefaultRedirectURL, http.StatusSeeOther)
		return
	}

	//get existing session
	session, err := h.store.GetExistingSession(w, r)
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		return
	}

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
	userDetail, err := h.kenarService.GetUserInformation(token.AccessToken)
	err = h.oauthService.InsertUser(userDetail.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = h.oauthService.InsertPost(session.PostToken, userDetail.UserId, token.AccessToken, token.RefreshToken, token.Expiry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//update session
	// deleting state from session because we dont need it after oauth
	// session.State = ""
	// add sessionKey to the reuqest
	// session.SessionKey = uuid.New().String()
	// log.Println("new session id is " + session.SessionKey)
	// log.Println("This is the access token " + token.AccessToken)
	// add the user to database if the user existed
	// add the post to the database

	// //save the new session
	// h.store.SaveSession(w, r, session)
	// log.Println(h.store.GetExistingSession(w, r))

	// // Save token in database
	// if err := h.oauthService.InsertOAuthData(
	// 	session.SessionKey,
	// 	token.AccessToken,
	// 	token.RefreshToken,
	// 	session.PostToken,
	// 	// time.Unix(token.Expiry.Unix(), 0),
	// 	token.Expiry,
	// ); err != nil {
	// 	http.Error(w, "Failed to save token in database"+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// url := fmt.Sprintf("https://oryx-meet-elf.ngrok-free.app/poi")
	// http.Redirect(w, r, url, http.StatusSeeOther)

}
