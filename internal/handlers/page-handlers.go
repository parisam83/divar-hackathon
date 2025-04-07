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
		http.Error(w, "post_token and return_url are required", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		log.Printf("User ID not found in context")
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Get property details
	property, err := p.kenarService.GetPropertyDetail(r.Context(), postToken)
	if err != nil {
		log.Printf("Failed to fetch property details: %v", err)
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
		log.Printf("Template error: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	tmp.ExecuteTemplate(w, "buyer_landing.html", data)
}

func (p *PageHandler) SellerDashboardHandler(w http.ResponseWriter, r *http.Request) {
	postToken := r.URL.Query().Get("post_token")
	return_url := r.URL.Query().Get("return_url")

	if postToken == "" || return_url == "" {
		log.Printf("post_token and return_url are required")
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
		log.Printf("Template error: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Token":        postToken,
		"RedirectLink": return_url,
		"PropertyData": property,
	}

	tmp.Execute(w, data)
}
