package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/internal/services"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/transport"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/utils"
)

type KenarHandler struct {
	store            *utils.SessionStore
	kenarService     *services.KenarService
	transportService *services.TransportService
}

type SubwayInfo struct {
	Distance string `json:"distance"`
	Name     string `json:"name"`
	Duration string `json:"duration"`
}

type AddToListingRequest struct {
	PostToken string                       `json:"post_token"`
	Amenity   transport.NearbyPOIsResponse `json:"amenities"`
}

func NewKenarHandler(store *utils.SessionStore, serv *services.KenarService, transportService *services.TransportService) *KenarHandler {
	return &KenarHandler{
		store:            store,
		kenarService:     serv,
		transportService: transportService,
	}
}

func (k *KenarHandler) GetPrice(w http.ResponseWriter, r *http.Request) {
	log.Printf("internal/handlers/GetPrice called")

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
		log.Printf("failed to decode request body: %v", err)
		utils.HanleError(w, r, http.StatusInternalServerError, "خطا در پردازش درخواست", "فرمت درخواست نامعتبر است", err.Error())
		return
	}

	price, err := k.transportService.GetPrice(r.Context(), strconv.FormatFloat(req.Origin.Lat, 'f', -1, 64), strconv.FormatFloat(req.Origin.Long, 'f', -1, 64), strconv.FormatFloat(req.Destination.Lat, 'f', -1, 64), strconv.FormatFloat(req.Destination.Long, 'f', -1, 64))
	if err != nil {
		log.Printf("failed to get price: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"snapp_price": price.SnappPrice,
		"tapsi_price": price.TapsiPrice,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode response: %v", err)
		utils.HanleError(w, r, http.StatusInternalServerError, "خطای سیستمی", "خطا در تولید پاسخ", err.Error())
		return
	}
}

func (k *KenarHandler) Poi(w http.ResponseWriter, r *http.Request) {
	log.Printf("internal/handlers/Poi called")

	userId, ok := r.Context().Value("user_id").(string)
	if !ok {
		log.Printf("User ID not found in context")
		utils.HanleError(w, r, http.StatusInternalServerError, "خطای احراز هویت", "کاربر شناسایی نشد", "User ID not found in context")
		return
	}
	log.Printf("userId: %s", userId)

	var req struct {
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lng"`
		PostToken string  `json:"post_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("failed to decode request body: %v", err)
		utils.HanleError(w, r, http.StatusInternalServerError, "خطا در پردازش درخواست", "فرمت درخواست نامعتبر است", err.Error())
		return
	}

	stationResult, err := k.transportService.FindNearestStation(r.Context(), userId, req.PostToken, strconv.FormatFloat(req.Latitude, 'f', -1, 64), strconv.FormatFloat(req.Longitude, 'f', -1, 64))
	if err != nil {
		log.Printf("failed to find nearest station: %v", err)
		http.Error(w, "failed to find nearest station: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(stationResult); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *KenarHandler) AddLocationWidget(w http.ResponseWriter, r *http.Request) {
	log.Printf("internal/handlers/AddLocationWidget called")

	userId, ok := r.Context().Value("user_id").(string)
	if !ok {
		log.Printf("User ID not found in context")
		utils.HanleError(w, r, http.StatusInternalServerError, "خطای احراز هویت", "کاربر شناسایی نشد", "User ID not found in context")
		return
	}

	var req AddToListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("failed to decode request body: %v", err)
		utils.HanleError(w, r, http.StatusBadRequest, "درخواست نامعتبر", "فرمت درخواست نامعتبر است", err.Error())
		return
	}

	if req.PostToken == "" {
		log.Printf("missing required fields: post_token")
		utils.HanleError(w, r, http.StatusBadRequest, "درخواست نامعتبر", "فیلدهای ضروری وارد نشده‌اند", "Missing post_token field")
		return
	}

	err := h.kenarService.PostLocationWidget(r.Context(), userId, req.PostToken, req.Amenity)
	if err != nil {
		log.Printf("failed to post widget: %v", err)
		utils.HanleError(w, r, http.StatusInternalServerError, "خطا در ثبت ویجت", "خطا در ثبت اطلاعات مکانی", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Location widget added successfully",
	})
	log.Printf("Location widget added successfully")
}

func (k *KenarHandler) GetOriginCoordinates(w http.ResponseWriter, r *http.Request) {
	log.Printf("internal/handlers/GetOriginCoordinates called")

	var req struct {
		PostToken string `json:"post_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("failed to decode request body: %v", err)
		utils.HanleError(w, r, http.StatusBadRequest, "درخواست نامعتبر", "خطا در خواندن اطلاعات درخواست", err.Error())
		return
	}
	post, err := k.kenarService.GetPropertyDetail(r.Context(), req.PostToken)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("post not found: %v", err)
			utils.HanleError(w, r, http.StatusNotFound, "یافت نشد", "آگهی موردنظر یافت نشد", err.Error())
			return
		}
		log.Printf("failed to get post data: %v", err)
		utils.HanleError(w, r, http.StatusInternalServerError, "خطای سیستمی", "خطا در دریافت اطلاعات آگهی", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct {
		Latitude  float64 `json:"Latitude"`
		Longitude float64 `json:"Longitude"`
		Title     string  `json:"title"`
	}{
		Title:     post.Title,
		Latitude:  post.Latitude,
		Longitude: post.Longitude,
	}); err != nil {
		log.Printf("failed to encode response: %v", err)
		utils.HanleError(w, r, http.StatusInternalServerError, "خطای سیستمی", "خطا در تولید پاسخ", err.Error())
		return
	}
}
