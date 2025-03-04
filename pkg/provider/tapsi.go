package provider

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/joho/godotenv"
)

type origin struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type destinations struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type metadata struct {
	FlowType    string `json:"flowType"`
	PreviewType string `json:"previewType"`
}

type tapsiRequest struct {
	Origin       *origin         `json:"origin"`
	Destinations []*destinations `json:"destinations"`
	Rider        *string         `json:"rider"`
	HasReturn    bool            `json:"hasReturn"`
	WaitingTime  int             `json:"waitingTime"`
	Gateway      string          `json:"gateway"`
	InitiatedVia string          `json:"initiatedVia"`
	Metadata     *metadata       `json:"metadata"`
}

type tapsiPrices struct {
	PassengerShare int `json:"passengerShare"`
}

type tapsiService struct {
	Prices []*tapsiPrices `json:"prices"`
}

type tapsiItems struct {
	Service *tapsiService `json:"service"`
}

type tapsiCategories struct {
	Items []*tapsiItems `json:"items"`
}

type tapsiData struct {
	Categories []*tapsiCategories `json:"categories"`
}

type tapsiResponse struct {
	Data *tapsiData `json:"data"`
}

type tapsi struct {
	_clck        string
	AccessToken  string
	RefreshToken string
	_clsk        string
}

func (t *tapsi) GetPriceEstimation(stroriginLat, stroriginLong, strdestinationLat, strdestinationLong string) int {

	//
	destinationLong, _ := strconv.ParseFloat(strdestinationLong, 64)
	destinationLat, _ := strconv.ParseFloat(strdestinationLat, 64)
	originLat, _ := strconv.ParseFloat(stroriginLat, 64)
	originLong, _ := strconv.ParseFloat(stroriginLong, 64)

	data := tapsiRequest{
		Origin: &origin{
			Latitude:  originLat,
			Longitude: originLong,
		},
		Destinations: []*destinations{
			{
				Latitude:  destinationLat,
				Longitude: destinationLong,
			},
		},
		Rider:        nil,
		HasReturn:    false,
		WaitingTime:  0,
		Gateway:      "CAB",
		InitiatedVia: "WEB",
		Metadata: &metadata{
			FlowType:    "DESTINATION_FIRST",
			PreviewType: "ORIGIN_FIRST",
		},
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	body := bytes.NewReader(dataBytes)
	req, err := http.NewRequest("POST", "https://api.tapsi.cab/api/v3/ride/preview", body)
	if err != nil {
		log.Fatal(err)
	}
	t.SetHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var jsonData tapsiResponse
	err = json.Unmarshal(bodyText, &jsonData)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("Status code is not 200")
	}

	if jsonData.Data == nil {
		log.Fatal("Data is nil")
	}

	if len(jsonData.Data.Categories) == 0 {
		log.Fatal("No categories found")
	}

	if len(jsonData.Data.Categories[0].Items) == 0 {
		log.Fatal("No items found")
	}

	if jsonData.Data.Categories[0].Items[0].Service == nil {
		log.Fatal("No service found")
	}

	if len(jsonData.Data.Categories[0].Items[0].Service.Prices) == 0 {
		log.Fatal("No prices found")
	}

	if jsonData.Data.Categories[0].Items[0].Service.Prices[0].PassengerShare == 0 {
		log.Fatal("Price is 0")
	}

	return jsonData.Data.Categories[0].Items[0].Service.Prices[0].PassengerShare
}

func (t *tapsi) SetHeader(req *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Referer", "https://app.tapsi.cab/")
	req.Header.Set("x-agent", "v2.2|passenger|WEBAPP|7.13.4||5.0")
	req.Header.Set("Origin", "https://app.tapsi.cab/")
	req.Header.Set("Cookie", "_clck="+t._clck+
		"accessToken="+t.AccessToken+
		"refreshToken="+t.RefreshToken+
		"_clsk="+t._clck)
}
