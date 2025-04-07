package handlers

import (
	"html/template"
	"log"
	"net/http"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/internal/services"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/utils"
)

type pageHandler struct {
	sessionStore *utils.SessionStore
	kenarService *services.KenarService
	taxiService  *services.TransportService
}

func NewPageHandler(
	sessionStore *utils.SessionStore,
	kenarService *services.KenarService,
	taxiService *services.TransportService,
) *pageHandler {
	return &pageHandler{
		sessionStore: sessionStore,
		kenarService: kenarService,
		taxiService:  taxiService,
	}
}
func (p *pageHandler) BuyerDashboardHandler(w http.ResponseWriter, r *http.Request) {
	postToken := r.URL.Query().Get("post_token")
	return_url := r.URL.Query().Get("return_url")

	if postToken == "" || return_url == "" {
		http.Error(w, "post_token and return_url are required", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Get property details
	property, err := p.kenarService.GetPropertyDetail(r.Context(), postToken)
	if err != nil {
		http.Error(w, "Failed to fetch property details", http.StatusInternalServerError)
		return
	}

	hasPurchased, err := p.kenarService.CheckUserPurchase(r.Context(), postToken, userID)

	// Render buyer template with data
	data := map[string]interface{}{
		"UserID":       userID,
		"PostToken":    postToken,
		"PropertyData": property,
		"IsPurchased":  hasPurchased,
	}
	tmp, err := template.ParseFiles("./web/buyer_landing.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	tmp.ExecuteTemplate(w, "buyer_landing.html", data)
	return
}
func (p *pageHandler) SellerDashboardHandler(w http.ResponseWriter, r *http.Request) {
	postToken := r.URL.Query().Get("post_token")
	return_url := r.URL.Query().Get("return_url")

	if postToken == "" || return_url == "" {
		http.Error(w, "post_token and return_url are required", http.StatusBadRequest)
		return
	}
	_, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	property, err := p.kenarService.GetPropertyDetail(r.Context(), postToken)
	if err != nil {
		http.Error(w, "Failed to fetch property details", http.StatusInternalServerError)
		return
	}

	tmp, err := template.ParseFiles("./web/landing.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Token":        postToken,
		"RedirectLink": return_url,
		"PropertyData": property,
	}
	tmp.ExecuteTemplate(w, "landing.html", data)
	return
}

func (p *pageHandler) AmenitiesPageHandler(w http.ResponseWriter, r *http.Request) {
	postToken := r.URL.Query().Get("post_token")
	latitude := r.URL.Query().Get("latitude")
	longitude := r.URL.Query().Get("longitude")
	title := r.URL.Query().Get("title")
	return_url := r.URL.Query().Get("return_url")
	log.Println(latitude)
	log.Println(longitude)
	log.Println(title)

	if postToken == "" || latitude == "" || longitude == "" || return_url == "" {
		http.Error(w, "post_token and latitude are longitude", http.StatusBadRequest)
		return
	}
	// Get user ID from context
	userId, ok := r.Context().Value("user_id").(string)
	log.Println(userId)
	log.Println(postToken)

	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	// check if user has privillage?
	IsOwner, err := p.kenarService.CheckPostOwnership(r.Context(), userId, postToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if !IsOwner {
		http.Error(w, "You dont have access to this page because you are not the owner", http.StatusUnauthorized)
		return
	}
	tmp, err := template.ParseFiles("./web/amenities_finder.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"PostToken": postToken,
		"ReturnUrl": return_url,
		"Latitude":  latitude,
		"Longitude": longitude,
		"Title":     title,
	}
	tmp.ExecuteTemplate(w, "amenities_finder.html", data)
	return

}
