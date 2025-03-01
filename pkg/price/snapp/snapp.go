package main

import (
	"encoding/json"
	"fmt"
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
	var data = strings.NewReader(`{"points":[{"lat":"` + originLat + `","lng":"` + originLong + `"},{"lat":"` + destinationLat + `","lng":"` + destinationLong + `"},null],"voucher_code":null,"service_types":[1,2,24],"priceriderecom":false,"tag":"0","hurryRaised":0}`)
	req, err := http.NewRequest("POST", "https://app.snapp.taxi/api/api-base/v2/passenger/newprice/s/6/0", data)
	if err != nil {
		return -1
	}
	setHeader(req)
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
	return int(jsonData["data"].(map[string]interface{})["prices"].([]interface{})[0].(map[string]interface{})["final"].(float64)) / 10
}

func setHeader(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://app.snapp.taxi/pre-ride?utm_source=landing&utm_medium=request-button&utm_campaign=taxi&_gl=1*6bvi14*_gcl_au*MTEzNjQxNTUwMy4xNzQwNTc4NzI0")
	req.Header.Set("Origin", "https://app.snapp.taxi")
	req.Header.Set("Authorization", "Bearer " + os.Getenv("SNAPP_ACCESS_TOKEN"))
	req.Header.Set("Cookie", "cookiesession1=" + os.Getenv("cookiesession1") +
	"_clck=" + os.Getenv("_clck") + "_ga_Y4QV007ERR=" + os.Getenv("_ga_Y4QV007ERR") + "_ga=" + os.Getenv("_ga") +
	"_ym_uid=" + os.Getenv("_ym_uid") + "_ym_d=" + os.Getenv("_ym_d") + "_ym_isad=" + os.Getenv("_ym_isad") + 
	"_clsk=" + os.Getenv("_clsk"))
}

func main() {
	fmt.Println(GetPrice("35.70427000000001", "51.344799999999964", "35.70507152245129", "51.35158062440473"))
}
