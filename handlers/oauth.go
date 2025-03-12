package handlers

import (
	"fmt"
	"log"
	"net/http"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/services"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/utils"

	"github.com/google/uuid"
)

func NewOAuthHandler(store *utils.SessionStore, serv *services.OAuthService) *oAuthHandler {

	return &oAuthHandler{
		store:        store,
		oauthService: serv,
	}
}

// func (h *oAuthHandler) getExistingSession(w http.ResponseWriter, r *http.Request) (*OAuthSession, error) {
// 	oauthSession, err := h.store.Get(r, SessionName)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s", "Failed to get session: "+err.Error())
// 	}

// 	data, ok := oauthSession.Values[SessionKey].([]byte)
// 	if !ok {
// 		return nil, fmt.Errorf("%s", "No session data found ")

// 	}

// 	session := &OAuthSession{}
// 	if err := json.Unmarshal(data, session); err != nil {
// 		return nil, fmt.Errorf("%s", "failed to decode session:"+err.Error())
// 	}
// 	return session, nil
// }

// func (h *oAuthHandler) saveSession(w http.ResponseWriter, r *http.Request, session *OAuthSession) (*OAuthSession, error) {
// 	oauthSession, err := h.store.Get(r, SessionName)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get session: %w", err)
// 	}

// 	sessionJson, err := json.Marshal(session)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s", "Failed to marshal session: "+err.Error())
// 	}

// 	oauthSession.Values[SessionKey] = sessionJson
// 	err = h.store.Save(r, w, oauthSession)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s", "Failed to save session: "+err.Error())
// 	}
// 	log.Printf("Session saved successfully with state: %s", session.State)
// 	return session, nil

// }

// func (h *oAuthHandler) createNewSession(w http.ResponseWriter, r *http.Request, postToken string) (*OAuthSession, error) {
// 	state := uuid.New().String()
// 	session := &OAuthSession{
// 		PostToken: postToken,
// 		State:     state,
// 	}
// 	return h.saveSession(w, r, session)
// }

func (h *oAuthHandler) buildScopes(postToken string) []string {
	oauthScopes := []Scope{
		{resourceType: POST_ADDON_CREATE, resourceID: postToken},
		{resourceType: USER_PHONE},
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
	// Add this function to your handlers
	log.Println("AddonOauth called")

	postToken := r.URL.Query().Get("post_token")
	callback_url := r.URL.Query().Get("return_url")

	if postToken == "" || callback_url == "" {
		http.Error(w, "post_token and return_url are required", http.StatusBadRequest)
		return
	}
	// check existing session
	session, err := h.store.GetExistingSession(w, r)
	if err == nil && session != nil {
		log.Println(session.SessionKey)
		log.Println(session.PostToken)
		log.Println("User has entered before?!")
		url := fmt.Sprintf("https://oryx-meet-elf.ngrok-free.app/poi")
		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	}
	//create new session
	session, err = h.store.CreateNewSession(w, r, postToken)
	if err != nil {
		http.Error(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("New session created with state: %s", session.State)

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

	//update session
	// deleting state from session because we dont need it after oauth
	session.State = ""
	//add sessionKey to the reuqest
	session.SessionKey = uuid.New().String()

	//save the new session
	h.store.SaveSession(w, r, session)

	// Save token in database
	if err := h.oauthService.InsertOAuthData(
		session.SessionKey,
		token.AccessToken,
		token.RefreshToken,
		session.PostToken,
		// time.Unix(token.Expiry.Unix(), 0),
		token.Expiry,
	); err != nil {
		http.Error(w, "Failed to save token in database"+err.Error(), http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("https://oryx-meet-elf.ngrok-free.app/poi")
	http.Redirect(w, r, url, http.StatusSeeOther)

}
