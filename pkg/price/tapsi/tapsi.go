package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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

func GetTapsiPriceEstimation(originLat, originLong, destinationLat, destinationLong float64) int {
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
	setHeader(req)

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

	return jsonData.Data.Categories[0].Items[0].Service.Prices[0].PassengerShare
}

func setHeader(req *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Referer", "https://app.tapsi.cab/")
	req.Header.Set("x-agent", "v2.2|passenger|WEBAPP|7.13.4||5.0")
	req.Header.Set("Origin", "https://app.tapsi.cab/")
	req.Header.Set("Cookie", "_clck="+os.Getenv("_clck")+
		"accessToken="+os.Getenv("TAPSI_ACCESS_TOKEN")+
		"refreshToken="+os.Getenv("TAPSI_REFRESH_TOKEN")+
		"_clsk="+os.Getenv("_clsk"))
}

func main() {
	fmt.Println(GetTapsiPriceEstimation(35.6895, 51.3890, 35.7741, 51.5112))
}
