package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/database/db"
)

type KenarService struct {
	apiKey  string
	client  *http.Client
	domain  string
	queries *db.Queries
}

func NewKenarService(apiKey, domain string, queries *db.Queries) *KenarService {
	return &KenarService{
		apiKey:  apiKey,
		client:  http.DefaultClient,
		domain:  domain, //https://api.divar.ir/v1/open-platform
		queries: queries,
	}
}
func (k *KenarService) doRequest(method, endpoint string, payload io.Reader) (*http.Request, error) {

	url := k.domain + endpoint
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", k.apiKey)
	return req, nil

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

// We implement json.Marshaler interface for custom marshaling
func (r Row) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]map[string]interface{}{
		r.Key: r.Data,
	})
}

type Payload struct {
	Widgets []Row `json:"widgets"`
}

// type Payload struct {
// 	Widgets []Widget `json:"widgets"`
// }

func (s *KenarService) PostWidgets(postToken string) {
	log.Println("Post widgets")
	payload := Payload{
		Widgets: []Row{
			{"title_row", map[string]interface{}{"text": "Sample Title"}},
			{"subtitle_row", map[string]interface{}{"text": "Sample Subtitle"}},
			{"description_row", map[string]interface{}{"text": "Sample Description", "has_divider": false, "expandable": false}},
		},
	}

	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	// change ittttttttt
	url := "https://api.divar.ir/v2/open-platform/addons/post/" + postToken
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(jsonData)))
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}
	// req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", s.apiKey)
	req.Header.Set("x-access-token", "ory_at_oGndTdJE-Cfq8-fuAAmFuHI_itqopsk7Pr8zQiBkPEQ.2FxEb-7SlGNzhF5tnduTXgfUxqoFeOnNlEzwHnEbunw")
	log.Println("=======================================")
	fmt.Println(s.client)
	res, err := s.client.Do(req)
	log.Println(res.Status)
	log.Println(err)
	log.Println("=======================================")
	if err != nil {
		log.Println("request failed: %w", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("failed to read response body: %w", err)
	}
	var j map[string]interface{}
	err = json.Unmarshal(body, &j)
	log.Println(j)
	if err != nil {
		log.Println("Error parsing response:", err)
		return
	}

}

func (k *KenarService) GetCoordinates(postToken string) {
	k.PostWidgets(postToken)

	// req, err := k.doRequest(http.MethodGet, "/finder/post/"+postToken, nil)
	// if err != nil {
	// 	log.Println("Error creating request:", err)
	// 	return
	// }
	// resp, err := k.client.Do(req)
	// if err != nil {
	// 	log.Println("Error sending request" + err.Error())
	// 	return
	// }
	// defer resp.Body.Close()
	// log.Println("Response.staus code:", resp.StatusCode)
	// if resp.StatusCode != http.StatusOK {
	// 	log.Println("Error response code:", resp.StatusCode)
	// 	return
	// // }
	// body, _ := io.ReadAll(resp.Body)
	// var jsonData map[string]interface{}
	// err = json.Unmarshal(body, &jsonData)
	// if err != nil {
	// 	log.Println("Error parsing response:", err)
	// 	return
	// }

	// fmt.Println(jsonData["data"].(map[string]interface{})["latitude"])
	// fmt.Println(jsonData["data"].(map[string]interface{})["longitude"])

	log.Println("Request successful")
}
