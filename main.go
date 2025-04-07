package main

import (
	"context"
	"fmt"
	"log"
	"mime"
	"net/http"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/internal/handlers"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/internal/services"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/configs"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/database"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/database/db"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/transport"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/utils"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	router     *mux.Router
	config     *configs.Config
	dbConn     *pgxpool.Pool
	handlers   *appHandlers
	services   *appServices
}

type appServices struct {
	kenarService *services.KenarService
	oauthService *services.OAuthService
}

type appHandlers struct {
	kenar *handlers.KenarHandler
	oauth *handlers.OAuthHandler
	page  *handlers.PageHandler
	jwt   *utils.JWTManager
}

func NewApp() *App {
	return &App{
		router: mux.NewRouter(),
	}
}

func (a *App) Initialize() error {
	conf, err := configs.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configurations: %v", err)
	}
	a.config = conf

	if err := a.initializeDatabase(); err != nil {
		return err
	}

	if err := a.initializeServices(); err != nil {
		return err
	}

	a.setupRoutes()
	return nil
}

func (a *App) initializeDatabase() error {
	dbConn, err := database.ConnectToDatabase(context.Background(), a.config.Database)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}
	a.dbConn = dbConn
	return nil
}

func (a *App) initializeServices() error {
	query := db.New(a.dbConn)
	sessionStore := utils.NewSessionStore(&a.config.Session)
	jwtManager := utils.NewJWTManager(&a.config.Jwt)

	// Initialize transport services
	snapp := transport.NewSnapp(&a.config.Snapp)
	tapsi := transport.NewTapsi(&a.config.Tapsi)
	neshan := transport.NewNeshan(&a.config.Neshan)
	taxiService := services.NewTransportService(snapp, tapsi, neshan, query)

	// Initialize Kenar services
	kenarService := services.NewKenarService(
		a.config.Kenar.ApiKey,
		a.config.Kenar.OpenPlatformApi,
		query,
	)
	oauthService := services.NewOAuthService(a.config.Kenar, query, a.dbConn)
	a.services = &appServices{
		kenarService: kenarService,
		oauthService: oauthService,
	}

	// Initialize handlers
	kenarHandler := handlers.NewKenarHandler(sessionStore, kenarService, taxiService)
	pageHandler := handlers.NewPageHandler(sessionStore, kenarService, taxiService)
	oauthHandler := handlers.NewOAuthHandler(sessionStore, oauthService, kenarService, jwtManager)

	// Store handlers in the app struct for route setup
	a.handlers = &appHandlers{
		kenar: kenarHandler,
		oauth: oauthHandler,
		page:  pageHandler,
		jwt:   jwtManager,
	}

	return nil
}

func (a *App) setupRoutes() {
	// Set up MIME types
	mime.AddExtensionType(".css", "text/css")

	// API routes
	a.router.HandleFunc("/", homeHandler)
	a.router.HandleFunc("/poi", a.handlers.kenar.Poi)
	a.router.HandleFunc("/addon/oauth", a.handlers.oauth.AddonOauth)
	a.router.HandleFunc("/api/calculate-fare", a.handlers.kenar.GetPrice)
	a.router.HandleFunc("/api/find-amenities", a.handlers.jwt.JWTMiddlewear(a.handlers.kenar.Poi))
	a.router.HandleFunc("/api/add-to-ad", a.handlers.jwt.JWTMiddlewear(a.handlers.kenar.AddLocationWidget))
	a.router.HandleFunc("/api/record-purchase", a.handlers.jwt.JWTMiddlewear(a.handlers.kenar.RecordPurchase))
	a.router.HandleFunc("/api/get-origin", a.handlers.kenar.GetOriginCoordinates).Methods("POST")
	a.router.HandleFunc("/oauth/callback", a.handlers.oauth.OauthCallback)

	// Frontend routes
	a.router.HandleFunc("/api/seller/landing", a.handlers.jwt.JWTMiddlewear(a.handlers.page.SellerDashboardHandler)).Methods("GET")
	a.router.Handle("/api/buyer/landing", a.handlers.jwt.JWTMiddlewear(a.handlers.page.BuyerDashboardHandler)).Methods("GET")
	a.router.Handle("/api/serve/amenities-page", a.handlers.jwt.JWTMiddlewear(a.handlers.page.AmenitiesPageHandler)).Methods("GET")
	

	// Static file server
	htmlFileServer := http.FileServer(http.Dir("./web"))
	a.router.PathPrefix("/web/").Handler(http.StripPrefix("/web/", htmlFileServer))

	// Error handling
	a.router.HandleFunc("/error", utils.RenderErrorPage)
}

func (a *App) Run() error {
	port := a.config.Server.Port
	log.Printf("Server started on port %s", port)

	err := http.ListenAndServe(":"+port, a.router)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	return nil
}

func (a *App) Cleanup() {
	if a.dbConn != nil {
		a.dbConn.Close()
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func main() {
	app := NewApp()

	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer app.Cleanup()

	if err := app.Run(); err != nil {
		log.Fatalf("Error running app: %v", err)
	}
}
