package transport

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/configs"
)

const (
	searchApi    = "https://api.neshan.org/v1/search?term=%s&lat=%f&lng=%f"
	directionApi = "https://api.neshan.org/v4/direction?origin=%s,&destination=%s"
)

type SearchResult struct {
	Items []Items `json:"items"`
}

type DirectionResult struct {
	Route []Routes `json:"routes"`
}

type Routes struct {
	Legs []Leg `json:"legs"`
}

type Leg struct {
	Duration Duration `json:"duration"`
	Distance Distance `json:"distance"`
	Steps    []Step   `json:"steps"`
}

type Step struct {
	Instruction string    `json:"instruction"`
	Distance    *Distance `json:"distance"`
	Duration    *Duration `json:"duration"`
}

type Distance struct {
	Value float64 `json:"value"`
	Text  string  `json:"text"`
}

type Duration struct {
	Value float64 `json:"value"`
	Text  string  `json:"text"`
}

type Items struct {
	Title         string         `json:"title"`
	Address       string         `json:"address"`
	Category      string         `json:"category"`
	Type          string         `json:"type"`
	Region        string         `json:"region"`
	Neighbourhood string         `json:"neighbourhood"`
	Location      SearchLocation `json:"location"`
}

type SearchLocation struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type StationResponse struct {
	ClosestStation string  `json:"closest"`
	StationLat     float64 `json:"stationLat"`
	StationLong    float64 `json:"stationLong"`
	TotalDuration  string  `json:"totalDuration"`
	TotalDistance  string  `json:"totalDistance"`
}

type Neshan struct {
	apiKey string
}

func NewNeshan(s *configs.NeshanConfig) *Neshan {
	return &Neshan{
		apiKey: s.NeshanApiKey,
	}

}

func (n *Neshan) GetSearchResult(startLat, startLong float64) ([]Items, error) {
	searchURL := fmt.Sprintf(searchApi, "%D9%85%D8%AA%D8%B1%D9%88", startLat, startLong)
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Api-Key", n.apiKey)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s", resp.Status)
	}

	var searchResult SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, err
	}

	if len(searchResult.Items) == 0 {
		return nil, fmt.Errorf("No station found")
	}

	var nearbyStations []Items
	for _, item := range searchResult.Items {
		if distance(startLat, startLong, item.Location.Y, item.Location.X) <= 4 {
			nearbyStations = append(nearbyStations, item)
		}
	}

	if len(nearbyStations) == 0 {
		return nil, fmt.Errorf("No station found within 4 km radius")
	}

	return nearbyStations, nil
}

func (n *Neshan) GetDirectionResult(origin, destination string) (DirectionResult, error) {
	directionUrl := fmt.Sprintf(directionApi, origin, destination)
	req, err := http.NewRequest("GET", directionUrl, nil)
	if err != nil {
		return DirectionResult{}, err
	}

	req.Header.Set("Api-Key", n.apiKey)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return DirectionResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DirectionResult{}, fmt.Errorf("Error: %s", resp.Status)
	}

	var directionResult DirectionResult
	if err := json.NewDecoder(resp.Body).Decode(&directionResult); err != nil {
		return DirectionResult{}, err
	}

	return directionResult, nil
}

func distance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Radius of the Earth in km
	dLat := (lat2 - lat1) * (math.Pi / 180)
	dLon := (lon2 - lon1) * (math.Pi / 180)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*(math.Pi/180))*math.Cos(lat2*(math.Pi/180))*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func (n *Neshan) GetSubwayStation(startLatstr, startLongstr string) (*StationResponse, error) {
	if startLatstr == "" || startLongstr == "" {
		return nil, fmt.Errorf("startLat and startLong are required")
	}
	startLat, err := strconv.ParseFloat(startLatstr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid startLat parameter")
	}

	startLong, err := strconv.ParseFloat(startLongstr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid startLong parameter")
	}

	nearbyStations, err := n.GetSearchResult(startLat, startLong)
	if err != nil {
		return nil, err
	}

	closest := nearbyStations[0]
	stationLat := closest.Location.Y
	stationLong := closest.Location.X
	origin := fmt.Sprintf("%f", startLat) + "," + fmt.Sprintf("%f", startLong)
	destination := fmt.Sprintf("%f", closest.Location.Y) + "," + fmt.Sprintf("%f", closest.Location.X)

	directionResult, err := n.GetDirectionResult(origin, destination)
	if err != nil {
		return nil, err
	}

	totalDuration := directionResult.Route[0].Legs[0].Duration.Text
	totalDistance := directionResult.Route[0].Legs[0].Distance.Text

	response := &StationResponse{
		ClosestStation: closest.Title,
		StationLat:     stationLat,
		StationLong:    stationLong,
		TotalDuration:  totalDuration,
		TotalDistance:  totalDistance,
	}
	return response, nil
}
