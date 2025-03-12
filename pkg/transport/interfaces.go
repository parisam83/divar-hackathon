package transport

type PriceProvider interface {
	GetPriceEstimation(originLat, originLong, destinationLat, destinationLong string) int
}

type LocationProvider interface {
	GetSubwayStation(startLatstr, startLongstr string) (*StationResponse, error)
}
