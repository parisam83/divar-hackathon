package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
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

func addon_oauth(w http.ResponseWriter, r *http.Request) {

	post_token := r.URL.Query().Get("post_token")

	if post_token == "" {
		http.Error(w, "post_token is required", http.StatusBadRequest)
		return
	}
	callback_url := r.URL.Query().Get("return_url") //the adress the user will be redirected after oauth
	if callback_url == "" {
		http.Error(w, "return_url is required", http.StatusBadRequest)
		return
	}

	// TODO
	// create a post with it's token in databse

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

	conf := &oauth2.Config{
		ClientID:     os.Getenv("KENAR_APP_SLUG"),
		ClientSecret: os.Getenv("KENAR_OAUTH_SECRET"),
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.divar.ir/oauth2/auth",
			TokenURL: "https://provider.com/o/oauth2/auth",
		},
	}
	state := uuid.New().String()
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	//create a specifc session for user
	//redirect to the

}
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	http.HandleFunc("/addon/oauth", addon_oauth)
	http.ListenAndServe(":8080", nil)
}
