package snapp

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type snappPoint struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type snappRequest struct {
	Points         []snappPoint `json:"points"`
	VoucherCode    any          `json:"voucher_code"`
	ServiceTypes   []int        `json:"service_types"`
	Priceriderecom bool         `json:"priceriderecom"`
	Tag            string       `json:"tag"`
	HurryRaised    int          `json:"hurryRaised"`
}

type snappPrices struct {
	Final int `json:"final"`
}

type snappData struct {
	Prices []snappPrices `json:"prices"`
}
type snappResponse struct {
	Data snappData `json:"data"`
}

func GetSnappPriceEstimation(originLat, originLong, destinationLat, destinationLong string) int {
	data := snappRequest{
		Points: []snappPoint{
			{Lat: originLat, Lng: originLong},
			{Lat: destinationLat, Lng: destinationLong},
		},
		VoucherCode:    nil,
		ServiceTypes:   []int{1, 2, 24},
		Priceriderecom: false,
		Tag:            "0",
		HurryRaised:    0,
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	body := bytes.NewReader(dataBytes)
	req, err := http.NewRequest("POST", "https://app.snapp.taxi/api/api-base/v2/passenger/newprice/s/6/0", body)
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

	var jsonData snappResponse
	err = json.Unmarshal(bodyText, &jsonData)
	if err != nil {
		log.Fatal(err)
	}

	return jsonData.Data.Prices[0].Final / 10
}

func setHeader(req *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://app.snapp.taxi/pre-ride?utm_source=landing&utm_medium=request-button&utm_campaign=taxi&_gl=1*6bvi14*_gcl_au*MTEzNjQxNTUwMy4xNzQwNTc4NzI0")
	req.Header.Set("Origin", "https://app.snapp.taxi")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SNAPP_ACCESS_TOKEN"))
	req.Header.Set("Cookie", "cookiesession1="+os.Getenv("cookiesession1")+
		"_clck="+os.Getenv("_clck")+"_ga_Y4QV007ERR="+os.Getenv("_ga_Y4QV007ERR")+"_ga="+os.Getenv("_ga")+
		"_ym_uid="+os.Getenv("_ym_uid")+"_ym_d="+os.Getenv("_ym_d")+"_ym_isad="+os.Getenv("_ym_isad")+
		"_clsk="+os.Getenv("_clsk"))
}
