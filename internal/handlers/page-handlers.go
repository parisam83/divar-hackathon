package handlers

import (
	"html/template"
	"log"
	"net/http"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/internal/services"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/utils"
)

type PageHandler struct {
	sessionStore *utils.SessionStore
	kenarService *services.KenarService
	taxiService  *services.TransportService
}

func NewPageHandler(
	sessionStore *utils.SessionStore,
	kenarService *services.KenarService,
	taxiService *services.TransportService,
) *PageHandler {
	return &PageHandler{
		sessionStore: sessionStore,
		kenarService: kenarService,
		taxiService:  taxiService,
	}
}

func (p *PageHandler) BuyerDashboardHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("internal/handlers/BuyerDashboardHandler called")

	postToken := r.URL.Query().Get("post_token")
	return_url := r.URL.Query().Get("return_url")

	if postToken == "" || return_url == "" {
		log.Printf("post_token and return_url are required")
		utils.HanleError(w, r, http.StatusBadRequest, "درخواست نامعتبر", "خطا در خواندن اطلاعات درخواست", "post_token and return_url are required")
		return
	}
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		log.Printf("User ID not found in context")
		utils.HanleError(w, r, http.StatusUnauthorized, "خطای احراز هویت", "کاربر شناسایی نشد", "User ID not found in context")
		return
	}

	// Get property details
	property, err := p.kenarService.GetPropertyDetail(r.Context(), postToken)
	if err != nil {
		log.Printf("Failed to fetch property details: %v", err)
		utils.HanleError(w, r, http.StatusInternalServerError, "خطای سیستمی", "خطا در پردازش درخواست", err.Error())
		return
	}

	hasPurchased, err := p.kenarService.CheckUserPurchase(r.Context(), postToken, userID)
	if err != nil {
		log.Printf("Failed to check purchase status: %v", err)
		utils.HanleError(w, r, http.StatusInternalServerError, "خطای سیستمی", "خطا در پردازش درخواست", err.Error())
	}
	// Render buyer template with data
	data := map[string]interface{}{
		"UserID":       userID,
		"PostToken":    postToken,
		"PropertyData": property,
		"IsPurchased":  hasPurchased,
	}

	tmp, err := template.ParseFiles("./web/buyer_landing.html")
	if err != nil {
		log.Printf("Template error: %v", err)
		utils.HanleError(w, r, http.StatusInternalServerError, "خطای سیستمی", "خطا در پردازش درخواست", err.Error())
		return
	}
	tmp.ExecuteTemplate(w, "buyer_landing.html", data)
}

func (p *PageHandler) SellerDashboardHandler(w http.ResponseWriter, r *http.Request) {
	postToken := r.URL.Query().Get("post_token")
	return_url := r.URL.Query().Get("return_url")

	if postToken == "" || return_url == "" {
		log.Printf("post_token and return_url are required")
		utils.HanleError(w, r, http.StatusBadRequest, "درخواست نامعتبر", "خطا در خواندن اطلاعات درخواست", "post_token and return_url are required")
		return
	}
	_, ok := r.Context().Value("user_id").(string)
	if !ok {
		utils.HanleError(w, r, http.StatusUnauthorized, "خطای احراز هویت", "کاربر شناسایی نشد", "User ID not found in context")
		return
	}
	property, err := p.kenarService.GetPropertyDetail(r.Context(), postToken)
	if err != nil {
		utils.HanleError(w, r, http.StatusInternalServerError, "خطای سیستمی", "خطا در پردازش درخواست", err.Error())
		return
	}

	tmp, err := template.ParseFiles("./web/landing.html")
	if err != nil {
		log.Printf("Template error: %v", err)
		utils.HanleError(w, r, http.StatusInternalServerError, "خطای سیستمی", "خطا در پردازش درخواست", err.Error())
		return
	}
	data := map[string]interface{}{
		"Token":        postToken,
		"RedirectLink": return_url,
		"PropertyData": property,
	}
	tmp.ExecuteTemplate(w, "landing.html", data)
}

func (p *PageHandler) AmenitiesPageHandler(w http.ResponseWriter, r *http.Request) {
	postToken := r.URL.Query().Get("post_token")
	latitude := r.URL.Query().Get("latitude")
	longitude := r.URL.Query().Get("longitude")
	title := r.URL.Query().Get("title")
	return_url := r.URL.Query().Get("return_url")

	if postToken == "" || latitude == "" || longitude == "" || return_url == "" {
		log.Printf("post_token, latitude, longitude and return_url are required")
		utils.HanleError(w, r, http.StatusBadRequest, "درخواست نامعتبر", "خطا در خواندن اطلاعات درخواست", "post_token, latitude, longitude and return_url are required")
		return
	}
	// Get user ID from context
	userId, ok := r.Context().Value("user_id").(string)
	if !ok {
		utils.HanleError(w, r, http.StatusInternalServerError, "خطای احراز هویت", "کاربر شناسایی نشد", "User ID not found in context")
		return
	}
	// check if user has privillage?
	IsOwner, err := p.kenarService.CheckPostOwnership(r.Context(), userId, postToken)
	if err != nil {
		utils.HanleError(w, r, http.StatusInternalServerError, "خطای سیستمی", "خطا در پردازش درخواست", err.Error())
	}
	if !IsOwner {
		utils.HanleError(w, r, http.StatusUnauthorized, "خطای احراز هویت", "شما به این صفحه دسترسی ندارید", "You dont have access to this page because you are not the owner")
		return
	}
	tmp, err := template.ParseFiles("./web/amenities_finder.html")
	if err != nil {
		log.Printf("Template error: %v", err)
		utils.HanleError(w, r, http.StatusInternalServerError, "خطای سیستمی", "خطا در پردازش درخواست", err.Error())
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
}
