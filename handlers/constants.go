package handlers

import (
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/services"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/utils"
)

type OauthResourceType string

const (
	POST_ADDON_CREATE  OauthResourceType = "POST_ADDON_CREATE"
	USER_PHONE         OauthResourceType = "USER_PHONE"
	OFFLINE_ACCESS     OauthResourceType = "offline_access"
	SessionName                          = "auth_session"
	SessionKey                           = "data"
	DefaultRedirectURL                   = "https://divar.ir/"
)

type Scope struct {
	resourceType OauthResourceType
	resourceID   string
}

type oAuthHandler struct {
	oauthService *services.OAuthService
	store        *utils.SessionStore
}

type OAuthSession struct {
	PostToken   string `json:"post_token"`
	State       string `json:"state"`
	CallbackURL string `json:"callback_url"`
	SessionKey  string `json:"session_key"`
}
