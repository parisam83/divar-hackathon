package provider

import "net/http"

type RideProvider interface {
	GetPriceEstimation(originLat, originLong, destinationLat, destinationLong string) int
	SetHeader(req *http.Request)
}
