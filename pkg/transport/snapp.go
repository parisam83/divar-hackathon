package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/configs"
)

type snappPoint struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type snappRequest struct {
	Points         []*snappPoint `json:"points"`
	VoucherCode    *string       `json:"voucher_code"`
	ServiceTypes   []int         `json:"service_types"`
	Priceriderecom bool          `json:"priceriderecom"`
	Tag            string        `json:"tag"`
	HurryRaised    int           `json:"hurryRaised"`
}

type snappPrices struct {
	Final int `json:"final"`
}

type snappData struct {
	Prices []*snappPrices `json:"prices"`
}
type snappResponse struct {
	Data *snappData `json:"data"`
}

type Snapp struct {
	accessToken    string
	cookiesession1 string
	clck           string
	ga_Y4QV007ERR  string
	ga             string
	ym_d           string
	ym_uid         string
	ym_isad        string
	clsk           string
}

func NewSnapp(s *configs.SnappConfig) *Snapp {
	return &Snapp{
		accessToken:    s.ApiKey,
		cookiesession1: s.CookieSession,
		clck:           s.Clck,
		clsk:           s.Clsk,
		ym_d:           s.YandexDate,
		ym_uid:         s.YandexUID,
		ym_isad:        s.YandexAd,
		ga_Y4QV007ERR:  s.GATracking,
		ga:             s.GA,
	}

}

func (s *Snapp) GetPriceEstimation(originLat, originLong, destinationLat, destinationLong string) int {
	fmt.Println("dibididbididbi in snapp.go")
	data := snappRequest{
		Points: []*snappPoint{
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
	fmt.Println("85 in snapp.go")
	body := bytes.NewReader(dataBytes)
	req, err := http.NewRequest("POST", "https://app.snapp.taxi/api/api-base/v2/passenger/newprice/s/6/0", body)
	if err != nil {
		log.Fatal(err)
	}
	s.SetHeader(req)
	log.Println("92 in snapp.go")
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
	log.Println("104 in snapp.go")
	var jsonData snappResponse
	err = json.Unmarshal(bodyText, &jsonData)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp.StatusCode)

	if resp.StatusCode != 200 {
		log.Fatal("Status code is not 200")
	}

	if jsonData.Data == nil {
		log.Fatal("No data found")
	}

	if len(jsonData.Data.Prices) == 0 {
		log.Fatal("No price found")
	}

	if jsonData.Data.Prices[0].Final == 0 {
		log.Fatal("Price is 0")
	}

	return jsonData.Data.Prices[0].Final / 10
}

func (s *Snapp) SetHeader(req *http.Request) {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://app.snapp.taxi/pre-ride?utm_source=landing&utm_medium=request-button&utm_campaign=taxi&_gl=1*6bvi14*_gcl_au*MTEzNjQxNTUwMy4xNzQwNTc4NzI0")
	req.Header.Set("Origin", "https://app.snapp.taxi")
	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	req.Header.Set("Cookie", "cookiesession1="+s.cookiesession1+
		"_clck="+s.clck+"_ga_Y4QV007ERR="+s.ga_Y4QV007ERR+"_ga="+s.ga+
		"_ym_uid="+s.ym_uid+"_ym_d="+s.ym_d+"_ym_isad="+s.ym_isad+
		"_clsk="+s.clsk)
	// Print all headers
	// for key, values := range req.Header {
	// 	for _, value := range values {
	// 		fmt.Printf("%s: %s\n", key, value)
	// 	}
	// }

}
