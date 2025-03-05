package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/services"
)

type KenarHandler struct {
	kenarService *services.KenarService
}

func NewKenarHandler(serv *services.KenarService) *KenarHandler {

	return &KenarHandler{
		kenarService: serv,
	}
}

func (k *KenarHandler) Poi(w http.ResponseWriter, r *http.Request) {
	log.Println("Kenar called")
	// how should i get the post token? using database or the sesion

	oauthSession, err := store.Get(r, "oauth-session")
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
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

	log.Println("post token is %s", session.PostToken)
	k.kenarService.GetCoordinates(session.PostToken)

}
