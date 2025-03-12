package services

import (
	"fmt"
	"log"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/provider"
)

type RideService struct {
	snapp provider.RideProvider
	tapsi provider.RideProvider
}

func NewRideService(snapp, tapsi provider.RideProvider) *RideService {
	return &RideService{
		snapp: snapp,
		tapsi: tapsi,
	}
}

func (rs *RideService) GetPrice(originLat, originLong, destinationLat, destinationLong string) (
	map[string]int, error) {
	prices := make(map[string]int)
	fmt.Println("doooooooooooooooooooooooo in ride.go")
	prices["snapp"] = rs.snapp.GetPriceEstimation(originLat, originLong, destinationLat, destinationLong)
	log.Println(prices["snapp"])
	prices["tapsi"] = rs.tapsi.GetPriceEstimation(originLat, originLong, destinationLat, destinationLong)
	log.Println(prices["tapsi"])
	return prices, nil
}
