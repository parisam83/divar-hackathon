package handlers

import (
	"database/sql"
	"fmt"
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

func (k *KenarHandler) Poi(w http.ResponseWriter, r *http.Request) {
	log.Println("Kenar called")
	session, err := k.store.GetExistingSession(w, r)
	if err != nil {
		http.Error(w, "Failed to get session: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// sessionId := session.SessionKey
	// log.Println(sessionId)
	postToken := session.PostToken
	// oauth, err := k.kenarService.GetOAuthBySessionId(sessionId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "no session found", http.StatusNotFound)
		}
		http.Error(w, "Could not fetch data based on the sessionId "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("post token is %s\n", postToken)
	//fekr konam divar nmizare api bznm felan
	coordinates, err := k.kenarService.GetCoordinates(postToken)
	if err != nil {
		http.Error(w, "Couldn't fetch the coordinates "+err.Error(), http.StatusNotFound)
	}
	log.Println("origin: " + coordinates.Latitude)
	log.Println("origin: " + coordinates.Longitude)

	stationResult, err := k.transportService.FindNearestStation(coordinates.Latitude, coordinates.Longitude)
	// fmt.Println(stationResult)
	dest_lat := strconv.FormatFloat(stationResult.StationLat, 'f', -1, 64)
	dest_long := strconv.FormatFloat(stationResult.StationLong, 'f', -1, 64)
	// res, err := k.taxiService.GetPrice(coordinates.Latitude, coordinates.Longitude, dest_lat, dest_long)
	prices, err := k.transportService.GetPrice(coordinates.Latitude, coordinates.Longitude, dest_lat, dest_long)
	if err != nil {
		http.Error(w, "Could no fetch prices to the station", http.StatusNotFound)
	}
	// fmt.Println(prices)
	descriptionText := fmt.Sprintf(
		"نزدیک‌ترین ایستگاه مترو: %s\n"+
			"ههههههههفاصله تا ایستگاه: %s\n"+
			"زمان رسیدن به ایستگاه: %s\n\n"+
			"قیمت تخمینی تاکسی‌های آنلاین تا ایشتگاه مترو:\n"+
			"اسنپ: %v تومان\n"+
			"تپسی: %v تومان",
		stationResult.ClosestStation,
		stationResult.TotalDistance,
		stationResult.TotalDuration,
		prices["snapp"],
		prices["tapsi"],
	)
	// log.Println("=================================")
	log.Println(descriptionText)
	// k.kenarService.PostWidgets(postToken, oauth.AccessToken, descriptionText)

}
