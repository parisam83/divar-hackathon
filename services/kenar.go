package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/database/db"
)

type coordinate struct {
	Latitude  string
	Longitude string
}

type KenarService struct {
	apiKey  string
	client  *resty.Client
	domain  string
	queries *db.Queries
}

type Widget struct {
	TitleRow       map[string]interface{} `json:"title_row,omitempty"`
	SubtitleRow    map[string]interface{} `json:"subtitle_row,omitempty"`
	DescriptionRow map[string]interface{} `json:"description_row,omitempty"`
}

type Row struct {
	Key  string
	Data map[string]interface{}
}

func (r Row) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]map[string]interface{}{
		r.Key: r.Data,
	})
}

type Payload struct {
	Widgets []Row `json:"widgets"`
}

func NewKenarService(apiKey, domain string, queries *db.Queries) *KenarService {
	return &KenarService{
		apiKey:  apiKey,
		client:  resty.New().SetHeader("Content-Type", "application/json").SetHeader("X-Api-Key", apiKey),
		domain:  domain, //https://api.divar.ir/v1/open-platform
		queries: queries,
	}
}

func (k *KenarService) GetOAuthBySessionId(sessionId string) (db.Oauth, error) {
	return k.queries.GetOAuthBySessionId(context.Background(), sessionId)

}

func (k *KenarService) PostWidgets(postToken, accessToke, description string) {
	log.Println("Post widgets")
	payload := Payload{
		Widgets: []Row{
			{"title_row", map[string]interface{}{"text": "🚇 دسترسی به مترو"}},
			{"subtitle_row", map[string]interface{}{"text": "Sample Subtitle"}},
			{"description_row", map[string]interface{}{"text": fmt.Sprintf("%s", description), "has_divider": false, "expandable": false}},
		},
	}

	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	resp, err := k.client.R().SetHeader("x-access-token", accessToke).SetBody(jsonData).Post(AddWidgetUrl + postToken)
	if err != nil {
		log.Println("failed to post widgets %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		log.Println(resp.StatusCode())
		log.Println("unexpected code bruh!")
	}

	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("x-access-token", "ory_at_oGndTdJE-Cfq8-fuAAmFuHI_itqopsk7Pr8zQiBkPEQ.2FxEb-7SlGNzhF5tnduTXgfUxqoFeOnNlEzwHnEbunw")
	// log.Println("=======================================")
	// fmt.Println(s.client)
	// res, err := s.cl
	// log.Println(res.Status)
	// log.Println(err)
	// log.Println("=======================================")
	// if err != nil {
	// 	log.Println("request failed: %w", err)
	// }
	// defer res.Body.Close()
	// body, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	log.Println("failed to read response body: %w", err)
	// }
	// var j map[string]interface{}
	// err = json.Unmarshal(body, &j)
	// log.Println(j)
	// if err != nil {
	// 	log.Println("Error parsing response:", err)
	// 	return
	// }

}

func (k *KenarService) GetCoordinates(postToken string) (*coordinate, error) {
	// k.PostWidgets(postToken)

	resp, err := k.client.R().Get(GetPostUrl + postToken)

	// req, err := k.doRequest(http.MethodGet, "/finder/post/"+postToken, nil)
	if err != nil {
		return nil, fmt.Errorf("error sending request %w", err)
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(resp.Body(), &jsonData)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: " + err.Error())
	}
	data, ok := jsonData["data"].(map[string]interface{})
	log.Println(data)
	if !ok {
		return nil, fmt.Errorf("invalid response format: 'data' field not found or invalid")
	}

	lat, ok := data["latitude"].(float64)
	if !ok {
		return nil, fmt.Errorf("latitude not found or invalid type")
	}

	long, ok := data["longitude"].(float64)
	if !ok {
		return nil, fmt.Errorf("longitude not found or invalid type")
	}

	coords := &coordinate{
		Latitude:  strconv.FormatFloat(lat, 'f', -1, 64),
		Longitude: strconv.FormatFloat(long, 'f', -1, 64),
	}
	return coords, nil
}
