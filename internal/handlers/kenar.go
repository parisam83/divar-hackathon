package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/internal/services"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/utils"
)

type KenarHandler struct {
	store            *utils.SessionStore
	kenarService     *services.KenarService
	transportService *services.TransportService
}

func NewKenarHandler(store *utils.SessionStore, serv *services.KenarService, transportService *services.TransportService) *KenarHandler {
	return &KenarHandler{
		store:            store,
		kenarService:     serv,
		transportService: transportService,
	}
}

func (k *KenarHandler) GetPrice(w http.ResponseWriter, r *http.Request) {
	log.Println("Get price called")
	var req struct {
		PostToken string `json:"post_token"`

		Origin struct {
			Lat  float64 `json:"lat"`
			Long float64 `json:"lng"`
		} `json:"origin"`
		Destination struct {
			Lat  float64 `json:"lat"`
			Long float64 `json:"lng"`
		} `json:"destination"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	price, err := k.transportService.GetPrice(strconv.FormatFloat(req.Origin.Lat, 'f', -1, 64), strconv.FormatFloat(req.Origin.Long, 'f', -1, 64), strconv.FormatFloat(req.Destination.Lat, 'f', -1, 64), strconv.FormatFloat(req.Destination.Long, 'f', -1, 64))
	log.Println(price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"snappPrice": price["snapp"],
		"tapsiPrice": price["tapsi"],
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (k *KenarHandler) Poi(w http.ResponseWriter, r *http.Request) {
	log.Println("Kenar called")
	var req struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lng"`
		PostToken string  `json:"post_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	stationResult, err := k.transportService.FindNearestStation(req.PostToken, strconv.FormatFloat(req.Latitude, 'f', -1, 64), strconv.FormatFloat(req.Longitude, 'f', -1, 64))
	if err != nil {
		http.Error(w, "failed to find nearest station: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	response := PoiResponse{
		Subway: SubwayInfo{
			Distance: stationResult.TotalDistance,
			Name:     stationResult.ClosestStation,
			Duration: stationResult.TotalDuration,
		},
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}

}

type SubwayInfo struct {
	Distance string `json:"distance"`
	Name     string `json:"name"`
	Duration string `json:"duration"`
}

type PoiResponse struct {
	Subway SubwayInfo `json:"subway"`
}

type AddToListingRequest struct {
	PostToken string     `json:"post_token"`
	Subway    SubwayInfo `json:"subway"`
	Hospital  string     `json:"hospital,omitempty"` // omitempty since hospital isn't implemented yet
}

func (h *KenarHandler) AddLocationWidget(w http.ResponseWriter, r *http.Request) {
	//WE GET THE JWT THING AND USER ID
	var req AddToListingRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.PostToken == "" || req.Subway.Name == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// sample userId until we use jwt
	userId := "GuBFy0p90aKOX2nT-ptZ1-0jCm-pgGm-QH750nb56pY="
	err := h.kenarService.PostLocationWidget(r.Context(), userId, &services.PoiDetail{
		PostToken: req.PostToken,
		Subway: services.SubwayInfo{
			Distance: req.Subway.Distance,
			Name:     req.Subway.Name,
			Duration: req.Subway.Duration,
		},
		Hospital: req.Hospital,
	})

	if err != nil {
		http.Error(w, "Failed to post widget: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Location widget added successfully",
	})
}

func (k *KenarHandler) GetOriginCoordinates(w http.ResponseWriter, r *http.Request) {
	// this works with database
	var req struct {
		PostToken string `json:"post_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	post, err := k.kenarService.GetPropertyDetail(req.PostToken)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "no post found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get post: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct {
		Latitude  float64 `json:"Latitude"`
		Longitude float64 `json:"Longitude"`
	}{
		Latitude:  post.Latitude,
		Longitude: post.Longitude,
	}); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}

}
