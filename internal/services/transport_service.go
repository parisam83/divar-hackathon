package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/database/db"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/transport"
)

type TransportService struct {
	query            *db.Queries
	priceProviders   map[string]transport.PriceProvider
	locationProvider transport.LocationProvider
}

func NewTransportService(snapp *transport.Snapp, tapsi *transport.Tapsi, neshan *transport.Neshan, query *db.Queries) *TransportService {
	return &TransportService{
		priceProviders: map[string]transport.PriceProvider{
			"snapp": snapp,
			"tapsi": tapsi,
		},
		locationProvider: neshan,
		query:            query,
	}
}

func (t *TransportService) GetPrice(originLat, originLong, destinationLat, destinationLong string) (map[string]int, error) {
	prices := make(map[string]int)
	for name, client := range t.priceProviders {
		log.Println(name)
		price := client.GetPriceEstimation(originLat, originLong, destinationLat, destinationLong)
		prices[name] = price
	}
	return prices, nil
}

func (s *TransportService) FindNearestStation(postToken, originLat, originLong string) (*transport.StationResponse, error) {
	ctx := context.Background()

	// Convert string coordinates to float64
	lat, err := strconv.ParseFloat(originLat, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid latitude: %w", err)
	}

	lng, err := strconv.ParseFloat(originLong, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid longitude: %w", err)
	}

	result, err := s.query.GetToSubwayInfo(ctx, db.GetToSubwayInfoParams{
		Latitude:  lat,
		Longitude: lng,
	})

	// If we found a cached result, use it
	if err == nil {
		log.Printf("Cache hit for coordinates (%s, %s)", originLat, originLong)

		return &transport.StationResponse{
			ClosestStation: result.StationName,
			TotalDistance:  fmt.Sprintf("%d", result.Distance),
			TotalDuration:  fmt.Sprintf("%d", result.Duration),
		}, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		log.Printf("Database error when checking cache: %v", err)
		return nil, fmt.Errorf("error checking cache: %w", err)
	}

	apiResult, err := s.locationProvider.GetSubwayStation(originLat, originLong)
	if err != nil {
		return nil, fmt.Errorf("error calling external API: %w", err)
	}

	go s.cacheStationResult(ctx, postToken, apiResult)

	return apiResult, nil

}

func (t *TransportService) cacheStationResult(ctx context.Context, postId string, apiResult *transport.StationResponse) {
	poiId, err := t.query.UpsertSubwayStation(ctx, db.UpsertSubwayStationParams{
		Latitude:  apiResult.StationLat,
		Longitude: apiResult.StationLong,
		Name:      apiResult.ClosestStation,
	})
	if err != nil {
		log.Printf("Error getting/creating POI: %v", err)
		return
	}
	distance, err := strconv.Atoi(apiResult.TotalDistance)
	if err != nil {
		log.Printf("Error converting distance to int: %v", err)
		return
	}
	duration, err := strconv.Atoi(apiResult.TotalDuration)
	if err != nil {
		log.Printf("Error converting distance to int: %v", err)
		return
	}
	result, err := t.query.SaveTravelMetrics(ctx, db.SaveTravelMetricsParams{
		Distance:      int32(distance),
		Duration:      int32(duration),
		OriginID:      postId,
		DestinationID: poiId,
	})
	if err != nil {
		log.Printf("Error saving travel metrics: %v", err)
		return
	}
	if result.RowsAffected() == 0 {
		log.Printf("Travel metrics already exist for origin %s and destination %d", postId, poiId)
	}
	log.Printf("Successfully cached subway result for (%f, %f)", postId, apiResult.ClosestStation)
	//we need to save station into poi + travel metrics
}
