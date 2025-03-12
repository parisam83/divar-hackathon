package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/services"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/utils"
)

type KenarHandler struct {
	store        *utils.SessionStore
	kenarService *services.KenarService
	taxiService  *services.RideService
}

func NewKenarHandler(store *utils.SessionStore, serv *services.KenarService, taxiService *services.RideService) *KenarHandler {
	return &KenarHandler{
		store:        store,
		kenarService: serv,
		taxiService:  taxiService,
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
	postToken := session.PostToken
	// oauth, err := k.kenarService.GetOAuthBySessionId(sessionId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "no session found", http.StatusNotFound)
		}
		http.Error(w, "Could not fetch data based on the sessionId"+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("post token is %s\n", postToken)
	//fekr konam divar nmizare api bznm felan
	coordinates, err := k.kenarService.GetCoordinates(postToken)
	if err != nil {
		http.Error(w, "Couldn't fetch the coordinates "+err.Error(), http.StatusNotFound)
	}
	// log.Println(coordinates.Latitude)
	result, err := pkg.GetSubwayStationHandler(coordinates.Latitude, coordinates.Longitude)

	// result, err := pkg.GetSubwayStationHandler("35.705080137369734", "51.3493")
	fmt.Println(result)

	dest_lat := strconv.FormatFloat(result.StationLat, 'f', -1, 64)
	dest_long := strconv.FormatFloat(result.StationLong, 'f', -1, 64)
	// res, err := k.taxiService.GetPrice(coordinates.Latitude, coordinates.Longitude, dest_lat, dest_long)
	res, err := k.taxiService.GetPrice(coordinates.Latitude, coordinates.Longitude, dest_lat, dest_long)
	if err != nil {
		http.Error(w, "Could no fetch prices to the station", http.StatusNotFound)
	}
	fmt.Println(res)

}
