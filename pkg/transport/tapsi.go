package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/configs"
	"github.com/andybalholm/brotli"
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

type Tapsi struct {
	clck         string
	accessToken  string
	refreshToken string
	clsk         string
}

func NewTapsi(s *configs.TapsiConfig) *Tapsi {
	return &Tapsi{
		refreshToken: s.RefreshToken,
		accessToken:  s.AccessToken,
		clck:         s.Clck,
		clsk:         s.Clsk,
	}

}

func (t *Tapsi) GetPriceEstimation(stroriginLat, stroriginLong, strdestinationLat, strdestinationLong string) int {

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
	// log.Println("121 in tapsi.go")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	reader := brotli.NewReader(resp.Body)
	// log.Println(resp.StatusCode)

	// Read the decompressed data
	bodyText, err := io.ReadAll(reader)
	if err != nil {
		log.Println(bodyText)
		log.Fatal(err)
	}
	// log.Println("128 in tapsi.go")
	// log.Println("134 in tapsi.go")
	// fmt.Println(bodyText)
	// fmt.Println("Raw response:", string(bodyText))
	var jsonData tapsiResponse
	err = json.Unmarshal(bodyText, &jsonData)
	if err != nil {
		fmt.Println("test")
		log.Fatal(err)
	}
	// log.Println("140 in tapsi.go")
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

func (t *Tapsi) SetHeader(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:136.0) Gecko/20100101 Firefox/136.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://app.tapsi.cab/")
	req.Header.Set("x-agent", "v2.2|passenger|WEBAPP|7.14.7||5.0")
	req.Header.Set("Origin", "https://app.tapsi.cab/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Priority", "u=4")
	req.Header.Set("TE", "trailers")
	// req.Header.Set("Accept-Encoding", "br")

	req.Header.Set("Cookie", "_clck="+t.clck+"; "+
		"accessToken="+t.accessToken+"; "+
		"refreshToken="+t.refreshToken+"; "+
		"_clsk="+t.clsk)
	// for key, values := range req.Header {
	// 	for _, value := range values {
	// 		fmt.Printf("%s: %s\n", key, value)
	// 	}
	// }
}
