package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/configs"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/database/db"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/oauth2"
)

type OAuthService struct {
	appConf   configs.KenarConfig
	oauthConf *oauth2.Config
	queries   *db.Queries
}

func NewOAuthService(appConfig configs.KenarConfig, queries *db.Queries) *OAuthService {
	conf := &oauth2.Config{
		ClientID:     appConfig.AppSlug,
		ClientSecret: appConfig.OauthSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.divar.ir/oauth2/auth",
			TokenURL: "https://api.divar.ir/oauth2/token",
		},
	}

	return &OAuthService{
		appConf:   appConfig,
		oauthConf: conf,
		queries:   queries,
	}
}

func (s *OAuthService) InsertOAuthData(sessionId, accessToken, refreshToken, postToken string, expiresIn time.Time) error {
	err := s.queries.AddOAuthData(context.Background(), db.AddOAuthDataParams{
		SessionID:    sessionId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn: pgtype.Timestamp{
			Time:  expiresIn,
			Valid: true,
		},
		PostToken: postToken,
	})
	if err != nil {
		return fmt.Errorf("faild to add oauth data into the database: %w", err)
	}
	log.Println("New oauth data added successfully")
	return nil
}

func (s *OAuthService) GenerateAuthURL(scopes []string, state string) string {
	s.oauthConf.Scopes = scopes
	return s.oauthConf.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *OAuthService) ExchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := s.oauthConf.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %w", err)
	}
	return token, nil
}
