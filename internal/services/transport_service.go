package services

import (
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/transport"
)

type TransportService struct {
	priceProviders  map[string]transport.PriceProvider
	loctionProvider transport.LocationProvider
}

func NewTransportService(snapp *transport.Snapp, tapsi *transport.Tapsi, neshan *transport.Neshan) *TransportService {
	return &TransportService{
		priceProviders: map[string]transport.PriceProvider{
			"snapp": snapp,
			"tapsi": tapsi,
		},
		loctionProvider: neshan,
	}
}

func (t *TransportService) GetPrice(originLat, originLong, destinationLat, destinationLong string) (map[string]int, error) {
	prices := make(map[string]int)
	for name, client := range t.priceProviders {
		price := client.GetPriceEstimation(originLat, originLong, destinationLat, destinationLong)
		prices[name] = price
	}
	return prices, nil
}

func (s *TransportService) FindNearestStation(originLat, originLong string) (*transport.StationResponse, error) {
	return s.loctionProvider.GetSubwayStation(originLat, originLong)
}
