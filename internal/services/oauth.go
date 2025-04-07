package services

import (
	"context"
	"fmt"
	"log"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/configs"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/database/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"
)

type OAuthService struct {
	appConf   configs.KenarConfig
	oauthConf *oauth2.Config
	queries   *db.Queries
	db        *pgxpool.Pool
}

func NewOAuthService(appConfig configs.KenarConfig, queries *db.Queries, db *pgxpool.Pool) *OAuthService {
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
		db:        db,
	}
}

func (t *Transaction) Validate() error {
	log.Printf("internal/services/Validate called")
	if t.PropertyDetail == nil || t.UserDetail == nil {
		return fmt.Errorf("missing required transaction components")
	}

	// Basic validation for all transactions
	switch {
	case t.UserDetail.UserId == "":
		return fmt.Errorf("invalid user_id: cannot be empty")
	case t.PropertyDetail.PostID == "":
		return fmt.Errorf("invalid post_id: cannot be empty")
	case t.PropertyDetail.Title == "":
		return fmt.Errorf("invalid title: cannot be empty")
	case t.PropertyDetail.Latitude == 0 || t.PropertyDetail.Longitude == 0:
		return fmt.Errorf("invalid coordinates: lat=%f, long=%f", t.PropertyDetail.Latitude, t.PropertyDetail.Longitude)
	}

	// For non-buyers, validate TokenInfo
	if !t.IsBuyer {
		if t.TokenInfo == nil {
			return fmt.Errorf("missing token info for non-buyer transaction")
		}

		switch {
		case t.TokenInfo.AccessToken == "":
			return fmt.Errorf("invalid access_token: cannot be empty")
		case t.TokenInfo.RefreshToken == "":
			return fmt.Errorf("invalid refresh_token: cannot be empty")
		case t.TokenInfo.ExpiresIn.IsZero():
			return fmt.Errorf("invalid expires_in: cannot be zero")
		}
	}

	return nil
}

func (s *OAuthService) RegisterAuthData(ctx context.Context, input *Transaction) error {
	log.Printf("internal/services/RegisterAuthData called")
	err := input.Validate()
	if err != nil {
		return fmt.Errorf("invalid transaction data: %w", err)
	}

	log.Printf("Starting Transaction after OAuth for user id:%s and post_id: %s", input.UserDetail.UserId, input.PropertyDetail.PostID)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := s.queries.WithTx(tx)

	// Insert user
	result, err := qtx.InsertUser(ctx, input.UserDetail.UserId)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	if result.RowsAffected() == 0 {
		log.Printf("user already exists: %+v", input.UserDetail)
	} else {
		log.Printf("Successfully added user to the database: %+v", input.UserDetail)
	}

	// Insert post
	ownerID := pgtype.Text{Valid: false} // Default to NULL

	// Set owner_id only if this is not a buyer (i.e., it's a seller)
	if !input.IsBuyer {
		ownerID = pgtype.Text{
			String: input.UserDetail.UserId,
			Valid:  true,
		}
	}

	result, err = qtx.InsertPost(ctx, db.InsertPostParams{
		PostID:    input.PropertyDetail.PostID,
		Latitude:  input.PropertyDetail.Latitude,
		Longitude: input.PropertyDetail.Longitude,
		Title:     pgtype.Text{String: input.PropertyDetail.Title, Valid: true},
		OwnerID:   ownerID,
	})
	if err != nil {
		return fmt.Errorf("failed to insert post: %w", err)
	}
	if result.RowsAffected() == 0 {
		log.Printf("post already exists: %+v", input.PropertyDetail)
	} else {
		log.Printf("Successfully added post to the database: %+v", input.PropertyDetail)
	}

	// Only insert token data for post owners
	if !input.IsBuyer && input.TokenInfo != nil {
		result, err = qtx.InsertToken(ctx, db.InsertTokenParams{
			PostID:       input.PropertyDetail.PostID,
			UserID:       input.UserDetail.UserId,
			AccessToken:  input.TokenInfo.AccessToken,
			RefreshToken: input.TokenInfo.RefreshToken, // Note: You had AccessToken here which is likely a bug
			ExpiresAt:    pgtype.Timestamp{Time: input.TokenInfo.ExpiresIn, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to insert token: %w", err)
		}
		if result.RowsAffected() == 0 {
			log.Printf("Token already exists: %+v", input.TokenInfo)
		} else {
			log.Printf("Successfully added token to the database: %+v", input.TokenInfo)
		}
	}

	return tx.Commit(ctx)
}

func (s *OAuthService) IsUserPostOwner(ctx context.Context, userId, postId string) (bool, error) {
	log.Printf("internal/services/IsUserPostOwner called")
	isOwner, err := s.queries.CheckPostOwnership(ctx, db.CheckPostOwnershipParams{
		OwnerID: pgtype.Text{String: userId, Valid: true},
		PostID:  postId,
	})
	if err != nil {
		return false, fmt.Errorf("failed to get post owner: %w", err)
	}
	return isOwner, nil
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
