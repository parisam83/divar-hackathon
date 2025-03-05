package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type KenarService struct {
	apiKey string
	client *http.Client
	domain string
}

func NewKenarService(apiKey, domain string) *KenarService {
	return &KenarService{
		apiKey: apiKey,
		client: http.DefaultClient,
		domain: domain, //https://api.divar.ir/v1/open-platform
	}
}
func (k *KenarService) doRequest(method, endpoint string) (*http.Request, error) {

	url := k.domain + endpoint
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", k.apiKey)
	return req, nil

}

func (k *KenarService) GetCoordinates(postToken string) {
	req, err := k.doRequest(http.MethodGet, "/finder/post/"+postToken)
	if err != nil {
		log.Println("Error creating request:", err)
		return
	}
	resp, err := k.client.Do(req)
	if err != nil {
		log.Println("Error sending request" + err.Error())
		return
	}
	defer resp.Body.Close()
	log.Println("Response.staus code:", resp.StatusCode)
	// if resp.StatusCode != http.StatusOK {
	// 	log.Println("Error response code:", resp.StatusCode)
	// 	return
	// }
	body, _ := io.ReadAll(resp.Body)
	var jsonData map[string]interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		log.Println("Error parsing response:", err)
		return
	}
	// fmt.Println(string(body))

	fmt.Println(jsonData["data"].(map[string]interface{})["latitude"])
	fmt.Println(jsonData["data"].(map[string]interface{})["longitude"])

	log.Println("Request successful")
}
