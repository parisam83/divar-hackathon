package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/handlers"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/database/db"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/pkg/provider"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/services"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/utils"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       //  'file' source driver
	"github.com/gorilla/mux"
)

func main() {

	conf, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load configurations: %s\n", err)
	}

	conPool, err := utils.ConnectToDatabase(
		context.Background(),
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Host,
		conf.Database.DBName,
		conf.Database.Port,
		conf.Database.SSLMode,
		conf.Database.MaxConns,
		conf.Database.MinConns,
		conf.Database.MaxConnLifetimeJitterMinutes,
		conf.Database.MaxConnLifetimeMinutes,
		conf.Database.MaxConnIdleTimeMinutes,
	)
	if err != nil {
		log.Fatal("Error in databse setup " + err.Error())
	}

	query := db.New(conPool)
	sessionStore := utils.NewSessionStore(&conf.Session)
	fmt.Println(conf.Session.AuthKey)
	fmt.Println(reflect.ValueOf(conf.Session.AuthKey).Kind())

	oauthService := services.NewOAuthService(conf.Kenar, query)
	oauthHandler := handlers.NewOAuthHandler(sessionStore, oauthService)
	// oauthHandler := handlers.NewOAuthHandler(oauthService)

	snapp := provider.NewSnapp(&conf.Snapp)
	// log.Println(conf.Snapp.ApiKey)
	tapsi := provider.NewTapsi(&conf.Tapsi)
	taxiService := services.NewRideService(snapp, tapsi)

	kenarService := services.NewKenarService(conf.Kenar.ApiKey, "https://api.divar.ir/v1/open-platform", query)
	kenarHandler := handlers.NewKenarHandler(sessionStore, kenarService, taxiService)

	r := mux.NewRouter()

	r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	r.HandleFunc("/poi", kenarHandler.Poi)
	r.HandleFunc("/addon/oauth", oauthHandler.AddonOauth)
	// r.HandleFunc("/remove", oauthHandler.Remove)
	r.HandleFunc("/oauth/callback", oauthHandler.OauthCallback)
	port := conf.Server.Port
	log.Printf("Server started on port %s", port)
	http.ListenAndServe(":"+port, r)
}
