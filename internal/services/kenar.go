package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/database/db"
)

type KenarService struct {
	apiKey  string
	client  *resty.Client
	domain  string
	queries *db.Queries
}

func (r Row) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]map[string]interface{}{
		r.Key: r.Data,
	})
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
		return nil, fmt.Errorf("failed to fetch user information from divar: %w", err)
	}
	var UserInfo userInfo
	if err := json.Unmarshal(resp.Body(), &UserInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user information: %w", err)
	}
	return &UserInfo, nil
}

func getLocationTitle(locationType string) string {
	switch locationType {
	case "subway":
		return "🚇 دسترسی به مترو"
	case "hospital":
		return "🏥 دسترسی به مراکز درمانی"
	case "park":
		return "🌳 دسترسی به پارک"
	default:
		return "📍 دسترسی"
	}
}

type PoiDetail struct {
	PostToken string     `json:"post_token"`
	Subway    SubwayInfo `json:"subway"`
	Hospital  string     `json:"hospital,omitempty"` // omitempty since hospital isn't implemented yet
}
type SubwayInfo struct {
	Distance string `json:"distance"`
	Name     string `json:"name"`
	Duration string `json:"duration"`
}

func (k *KenarService) PostLocationWidget(ctx context.Context, userId string, poi_detail *PoiDetail) error {

	log.Printf("Posting information widget for post: %s", poi_detail.PostToken)
	token, err := k.queries.GetAccessTokenByUserIdPostId(ctx, db.GetAccessTokenByUserIdPostIdParams{
		ID:     userId,
		PostID: poi_detail.PostToken,
	})
	if err != nil {
		return fmt.Errorf("could not fetch access token from database")
	}

	// Parse distance for formatting
	distanceValue, err := strconv.ParseFloat(poi_detail.Subway.Distance, 64)
	if err != nil {
		return fmt.Errorf("error parsing distance: %w", err)
	}

	// Format distance text
	var distanceText string
	if distanceValue >= 1000 {
		distanceText = fmt.Sprintf("%.1f کیلومتر", distanceValue/1000)
	} else {
		distanceText = fmt.Sprintf("%.0f متر", distanceValue)
	}

	// Create a more structured widget with multiple rows
	payload := Payload{
		Widgets: []Row{
			{"title_row", map[string]interface{}{
				"text": "🚇 دسترسی به مترو",
			}},
			{"subtitle_row", map[string]interface{}{
				"text":        "نزدیک‌ترین ایستگاه مترو به این ملک:",
				"has_divider": true,
			}},
			{"description_row", map[string]interface{}{
				"text":        "🚉 " + poi_detail.Subway.Name,
				"has_divider": true,
			}},
			{"description_row", map[string]interface{}{
				"text":        "📍 فاصله تا ایستگاه: " + distanceText,
				"has_divider": true,
			}},
			{"description_row", map[string]interface{}{
				"text":        "🚗 مدت زمان با خودرو: " + poi_detail.Subway.Duration + " دقیقه",
				"has_divider": false,
			}},
		},
	}
	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	resp, err := k.client.R().
		SetHeader("x-access-token", token.AccessToken).
		SetBody(jsonData).
		Post(AddWidgetUrl + poi_detail.PostToken)

	if err != nil {
		return fmt.Errorf("failed to post widgets: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil
}

func (k *KenarService) GetPropertyDetail(postToken string) (*propertyInfo, error) {
	post, err := k.queries.GetPost(context.Background(), postToken)
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
	propertyInfo, err := k.fetchPropertyInfoFromDivar(postToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch property info: %w", err)
	}
	return propertyInfo, nil
}

func (k *KenarService) fetchPropertyInfoFromDivar(postToken string) (*propertyInfo, error) {

	resp, err := k.client.R().Get(GetPostUrl + postToken)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch property request %w", err)
	}
	var apiResponse propertyApiResponse
	err = json.Unmarshal(resp.Body(), &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse property response: %v", err.Error())
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
