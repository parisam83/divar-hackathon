package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5/pgtype"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/database/db"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/transport"
)

type KenarService struct {
	apiKey  string
	client  *resty.Client
	domain  string
	queries *db.Queries
}

func NewKenarService(apiKey, domain string, queries *db.Queries) *KenarService {
	return &KenarService{
		apiKey:  apiKey,
		client:  resty.New().SetHeader("Content-Type", "application/json").SetHeader("X-Api-Key", apiKey),
		domain:  domain, //https://api.divar.ir/v1/open-platform
		queries: queries,
	}
}

func (k *KenarService) GetUserDetail(accessToken string) (*userInfo, error) {
	k.client.SetHeader("x-access-token", accessToken)
	resp, err := k.client.R().Get("https://api.divar.ir/v1/open-platform/users")
	if err != nil {
		return nil, fmt.Errorf("failed to execute user information from divar: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to fetch user information from divar: %s", resp.String())
	}
	var UserInfo userInfo
	if err := json.Unmarshal(resp.Body(), &UserInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user information: %w", err)
	}
	log.Printf("fetched User info from Divar: %+v", UserInfo)
	return &UserInfo, nil
}

func formatDistance(distance int32) string {
	if distance >= 1000 {
		return fmt.Sprintf("%.1f کیلومتر", float64(distance)/1000)
	}
	return fmt.Sprintf("%d متر", distance)
}

func (r Row) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]map[string]interface{}{
		r.Key: r.Data,
	})
}

func (k *KenarService) PostLocationWidget(ctx context.Context, userId string, postToken string, amenities transport.NearbyPOIsResponse) error {

	log.Printf("Posting information widget for post: %s", postToken)
	token, err := k.queries.GetAccessTokenByUserIdPostId(ctx, db.GetAccessTokenByUserIdPostIdParams{
		ID:     userId,
		PostID: postToken,
	})
	if err != nil {
		return fmt.Errorf("could not fetch access token from database")
	}

	widgets := []Row{
		{"title_row", map[string]interface{}{"text": "🏙️ دسترسی به امکانات شهری"}},
		{"subtitle_row", map[string]interface{}{
			"text":        "امکانات نزدیک به این ملک:",
			"has_divider": true,
		}},
	}

	if amenities.Subway != nil && len(amenities.Subway.POIs) > 0 {
		text := ""
		for i, subway := range amenities.Subway.POIs {
			if i > 0 {
				text += "\n\n"
			}
			distanceText := formatDistance(subway.Distance)
			text += fmt.Sprintf("(%d) 🚉 ایستگاه مترو: %s\n", i+1, subway.Name) +
				fmt.Sprintf("📍 فاصله: %s\n", distanceText) +
				fmt.Sprintf("🚗 مدت زمان با خودرو: %d دقیقه", subway.Duration)
		}
		widgets = append(widgets,
			Row{"description_row", map[string]interface{}{
				"text":        text,
				"has_divider": true,
				"expandable":  true,
			}},
		)
	}

	if amenities.BusStation != nil && len(amenities.BusStation.POIs) > 0 {
		text := ""
		for i, bus := range amenities.BusStation.POIs {
			distanceText := formatDistance(bus.Distance)
			if i > 0 {
				text += "\n\n"
			}
			text += fmt.Sprintf("(%d) 🚌 ایستگاه اتوبوس: %s\n", i+1, bus.Name) +
				fmt.Sprintf("📍 فاصله: %s\n", distanceText) +
				fmt.Sprintf("🚗 مدت زمان با خودرو: %d دقیقه", bus.Duration)
		}

		widgets = append(widgets,
			Row{"description_row", map[string]interface{}{
				"text":        text,
				"has_divider": true,
				"expandable":  true,
			}},
		)
	}

	if amenities.Hospital != nil && len(amenities.Hospital.POIs) > 0 {
		text := ""
		for i, hospital := range amenities.Hospital.POIs {
			distanceText := formatDistance(hospital.Distance)
			if i > 0 {
				text += "\n\n"
			}
			text += fmt.Sprintf("(%d) 🏥 بیمارستان: %s\n", i+1, hospital.Name) +
				fmt.Sprintf("📍 فاصله: %s\n", distanceText) +
				fmt.Sprintf("🚗 مدت زمان با خودرو: %d دقیقه", hospital.Duration)
		}

		widgets = append(widgets,
			Row{"description_row", map[string]interface{}{
				"text":        text,
				"has_divider": true,
				"expandable":  true,
			}},
		)
	}

	if amenities.Supermarket != nil && len(amenities.Supermarket.POIs) > 0 {
		text := ""
		for i, market := range amenities.Supermarket.POIs {
			distanceText := formatDistance(market.Distance)
			if i > 0 {
				text += "\n\n"
			}
			text += fmt.Sprintf("(%d) 🛒 سوپرمارکت: %s\n", i+1, market.Name) +
				fmt.Sprintf("📍 فاصله: %s\n", distanceText) +
				fmt.Sprintf("🚗 مدت زمان با خودرو: %d دقیقه", market.Duration)
		}

		widgets = append(widgets,
			Row{"description_row", map[string]interface{}{
				"text":        text,
				"has_divider": true,
				"expandable":  true,
			}},
		)
	}
	payload := Payload{widgets}
	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}
	resp, err := k.client.R().
		SetHeader("x-access-token", token.AccessToken).
		SetBody(jsonData).
		Post(AddWidgetUrl + postToken)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("failed to post widgets: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("failed to set poi information on user's ad: %s", resp.String())
	}
	if resp.StatusCode() != http.StatusOK {
		log.Println(resp.StatusCode())
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil
}

func (k *KenarService) GetPropertyDetail(ctx context.Context, postToken string) (*propertyInfo, error) {
	post, err := k.queries.GetPost(ctx, postToken)
	if err == nil {
		log.Printf("Post %s found in database: location (lat: %f, lng: %f)", postToken, post.Latitude, post.Longitude)
		return &propertyInfo{
			PostID:    post.PostID,
			Latitude:  post.Latitude,
			Longitude: post.Longitude,
			Title:     post.Title.String,
		}, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to query property detail from database: %w", err)
	}
	propertyInfo, err := k.fetchPropertyInfoFromDivar(ctx, postToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch property info: %w", err)
	}
	return propertyInfo, nil
}

func (k *KenarService) fetchPropertyInfoFromDivar(ctx context.Context, postToken string) (*propertyInfo, error) {

	resp, err := k.client.R().SetContext(ctx).Get(GetPostUrl + postToken)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch property request %w", err)
	}
	var apiResponse propertyApiResponse
	err = json.Unmarshal(resp.Body(), &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse property response: %v", err.Error())
	}
	// if the ad has no coordinates, we should not save it
	if apiResponse.Data.Latitude == 0 && apiResponse.Data.Longitude == 0 {
		return nil, fmt.Errorf("property has missing or invalid coordinates (0,0), Can not have this post in application")
	}
	propertyInfo := &propertyInfo{
		PostID:    postToken,
		Latitude:  apiResponse.Data.Latitude,
		Longitude: apiResponse.Data.Longitude,
		Title:     apiResponse.Data.Title,
	}
	log.Printf("fetched property info from Divar: %+v", propertyInfo)
	return propertyInfo, nil
}
func (k *KenarService) InsertPostPurchase(ctx context.Context, postToken string, userID string) error {
	result, err := k.queries.InsertPostPurchase(ctx, db.InsertPostPurchaseParams{
		UserID:    userID,
		PostToken: postToken,
	})
	if err != nil {
		return fmt.Errorf("failed to record post purchase: %w", err)
	}
	if result.RowsAffected() == 0 {
		// this part should never be called actually
		log.Printf("Purchase record already exists for user %s and post %s",
			userID, postToken)
		return fmt.Errorf("purchase record already exists for user %s and post %s",
			userID, postToken)

	} else {
		log.Printf("Successfully recorded purchase for user %s and post %s", userID, postToken)
	}
	return nil

}

func (k *KenarService) CheckUserPurchase(ctx context.Context, postToken string, userID string) (bool, error) {
	hasPurchased, err := k.queries.CheckUserPurchase(ctx, db.CheckUserPurchaseParams{
		UserID:    userID,
		PostToken: postToken,
	})
	if err != nil {
		return false, fmt.Errorf("could not fetch the realtion of user from post-purchase")
	}
	return hasPurchased, nil
}

func (k *KenarService) CheckPostOwnership(ctx context.Context, userId, postId string) (bool, error) {
	log.Println("Checking if user is the owner of the post")
	isOwner, err := k.queries.CheckPostOwnership(ctx, db.CheckPostOwnershipParams{
		OwnerID: pgtype.Text{String: userId, Valid: true},
		PostID:  postId,
	})
	if err != nil {
		return false, fmt.Errorf("failed to get post owner: %w", err)
	}
	return isOwner, nil
}
