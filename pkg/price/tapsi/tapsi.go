package tapsi

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func GetPrice(originLat, originLong, destinationLat, destinationLong string) int {
	err := godotenv.Load()
	if err != nil {
		return -1
	}

	client := &http.Client{}
	var data = strings.NewReader(`{"origin":{"latitude":` + originLat + `,"longitude":` + originLong + `},"destinations":[{"latitude":` + destinationLat + `,"longitude":` + destinationLong + `}],"rider":null,"hasReturn":false,"waitingTime":0,"gateway":"CAB","initiatedVia":"WEB","metadata":{"flowType":"DESTINATION_FIRST","previewType":"ORIGIN_FIRST"}}`)
	req, err := http.NewRequest("POST", "https://api.tapsi.cab/api/v3/ride/preview", data)
	if err != nil {
		return -1
	}
	req.Header.Set("Referer", "https://app.tapsi.cab/")
	req.Header.Set("x-agent", "v2.2|passenger|WEBAPP|7.13.4||5.0")
	req.Header.Set("Origin", "https://app.tapsi.cab/")
	req.Header.Set("Cookie", "_clck="+os.Getenv("_clck")+
		"accessToken="+os.Getenv("TAPSI_ACCESS_TOKEN")+
		"refreshToken="+os.Getenv("TAPSI_REFRESH_TOKEN")+
		"_clsk="+os.Getenv("_clsk"))

	resp, err := client.Do(req)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(bodyText, &jsonData)
	if err != nil {
		return -1
	}
	return int(jsonData["data"].(map[string]interface{})["categories"].([]interface{})[0].(map[string]interface{})["items"].([]interface{})[0].(map[string]interface{})["service"].(map[string]interface{})["prices"].([]interface{})[0].(map[string]interface{})["passengerShare"].(float64))
}
