package services

import (
	"context"
	"os"

	"golang.org/x/oauth2"
)

type OAuthService struct {
	conf *oauth2.Config
}

func NewOAuthService() *OAuthService {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("KENAR_APP_SLUG"),
		ClientSecret: os.Getenv("KENAR_OAUTH_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.divar.ir/oauth2/auth",
			TokenURL: "https://api.divar.ir/oauth2/token",
		},
	}
	return &OAuthService{
		conf: conf,
	}
}

func (s *OAuthService) GenerateAuthURL(scopes []string, state string) string {
	s.conf.Scopes = scopes
	return s.conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *OAuthService) ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.conf.Exchange(ctx, code)
}
