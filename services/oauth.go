package services

import (
	"context"
	"log"
	"os"
	"time"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/database/db"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/oauth2"
)

type OAuthService struct {
	queries *db.Queries
	conf    *oauth2.Config
}

func NewOAuthService(queries *db.Queries) *OAuthService {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("KENAR_APP_SLUG"),
		ClientSecret: os.Getenv("KENAR_OAUTH_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.divar.ir/oauth2/auth",
			TokenURL: "https://api.divar.ir/oauth2/token",
		},
	}
	return &OAuthService{
		conf:    conf,
		queries: queries,
	}
}

func (s *OAuthService) AddOAuth(sessionId, accessToken, refreshToken, postToken string, expiresIn time.Time) error {
	err := s.queries.AddOAuthData(context.Background(), db.AddOAuthDataParams{
		SessionID:    sessionId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn: pgtype.Timestamp{
			Time: expiresIn,
			// InfinityModifier: pgtype.NegativeInfinity,
			Valid: true, // Key: Mark as non-NULL

		},
		PostToken: postToken,
	})
	log.Println(err)
	log.Println("new oauth added")
	return nil
}

func (s *OAuthService) GenerateAuthURL(scopes []string, state string) string {
	s.conf.Scopes = scopes
	return s.conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *OAuthService) ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.conf.Exchange(ctx, code)
}
