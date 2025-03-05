package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"git.divar.cloud/divar/girls-hackathon/realestate-poi/handlers"
	"git.divar.cloud/divar/girls-hackathon/realestate-poi/services"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Import the Postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Import the 'file' source driver
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type config struct {
	Database databaseConfig
}

type databaseConfig struct {
	Host                         string
	Port                         int
	Username                     string
	Password                     string
	DBName                       string
	SSLMode                      string
	MaxConns                     int32
	MinConns                     int32
	MaxConnLifetimeJitterMinutes int
	MaxConnLifetimeMinutes       int
	MaxConnIdleTimeMinutes       int
}

func loadConfig() (*config, error) {
	var config config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func ConnectToDatabase(
	serverContext context.Context,
	databaseUsername string,
	databasePassword string,
	databaseHost string,
	databaseName string,
	databasePort int,
	databaseSSLMode string,
	databaseMaxConns int32,
	databaseMinConns int32,
	databaseMaxConnLifetimeJitterMinutes int,
	databaseMaxConnLifetimeMinutes int,
	databaseMaxConnIdleTimeMinutes int,
) (*pgxpool.Pool, error) {

	//check if config is provided
	if databaseHost == "" || databaseName == "" || databaseUsername == "" || databasePassword == "" || databasePort == 0 {
		return nil, errors.Errorf("incomplete database config: userName: %s password: %s, host: %s, name: %s, port: %d",
			databaseUsername, databasePassword, databaseHost, databaseName, databasePort)
	}

	//create database url
	// example URL: postgres://username:password@localhost:5432/database_name
	pgURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(databaseUsername, databasePassword),
		Host:   fmt.Sprintf("%s:%d", databaseHost, databasePort),
		Path:   "/" + databaseName,
	}
	query := pgURL.Query()
	query.Add("sslmode", databaseSSLMode)
	// fmt.Printf(query.Encode() + "\n")
	pgURL.RawQuery = query.Encode()
	pgxPoolConfig, err := pgxpool.ParseConfig(pgURL.String())
	if err != nil {
		return nil, errors.Errorf("failed to parse postgres config URL: %w", err)
	}
	// Configure pgxpool Parameters
	pgxPoolConfig.MaxConns = databaseMaxConns
	pgxPoolConfig.MinConns = databaseMinConns
	pgxPoolConfig.MaxConnLifetime = time.Duration(databaseMaxConnLifetimeMinutes) * time.Minute
	pgxPoolConfig.MaxConnIdleTime = time.Duration(databaseMaxConnIdleTimeMinutes) * time.Minute
	pgxPoolConfig.MaxConnLifetimeJitter = time.Duration(databaseMaxConnLifetimeJitterMinutes) * time.Minute

	// Establish Connection Pool
	pool, err := pgxpool.NewWithConfig(serverContext, pgxPoolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create PGX connection pool: %w", err)
	}

	// Ping the Database to Verify Connection
	if err := pool.Ping(serverContext); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to reach database: %w", err)
	}

	fmt.Println("Connected to the database successfully.")

	m, err := migrate.New("file://./pkg/database/migrations", pgURL.String())
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("migration initialization failed: %w", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		pool.Close()
		return nil, fmt.Errorf("migration failed: %w", err)
	}
	return pool, nil

}
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// configFilePath := "./configs/db.yaml"

	// // var err error
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// viper.AutomaticEnv()

	// // check if config file is not provided
	// if configFilePath != "" {
	// 	viper.SetConfigFile(configFilePath)
	// 	if err := viper.ReadInConfig(); err != nil {
	// 		log.Fatalln("failed to read config file.")
	// 	}
	// }

	// var conf *config
	// conf, err = loadConfig()
	// if err != nil {
	// 	log.Fatalf("failed to load configurations: %s\n", err)
	// }
	// conPool, err := ConnectToDatabase(
	// 	context.Background(),
	// 	conf.Database.Username,
	// 	conf.Database.Password,
	// 	conf.Database.Host,
	// 	conf.Database.DBName,
	// 	conf.Database.Port,
	// 	conf.Database.SSLMode,
	// 	conf.Database.MaxConns,
	// 	conf.Database.MinConns,
	// 	conf.Database.MaxConnLifetimeJitterMinutes,
	// 	conf.Database.MaxConnLifetimeMinutes,
	// 	conf.Database.MaxConnIdleTimeMinutes,
	// )
	// if err != nil {
	// 	log.Fatal("Error in databse setup" + err.Error())
	// }
	// db.New(conPool)

	oauthService := services.NewOAuthService()
	kenarService := services.NewKenarService(
		"eyJhbGciOiJSUzI1NiIsImtpZCI6InByaXZhdGVfa2V5XzIiLCJ0eXAiOiJKV1QifQ.eyJhcHBfc2x1ZyI6InBsYW5ldC1yaXBwbGUtbGVnZW5kIiwiYXVkIjoic2VydmljZXByb3ZpZGVycyIsImV4cCI6MTc0NjM1NDM3NiwianRpIjoiNDU4YzM0MTMtZjlhYy0xMWVmLTlhMTktZmFkODI3M2I1OGM1IiwiaWF0IjoxNzQxMTcwMzc2LCJpc3MiOiJkaXZhciIsInN1YiI6ImFwaWtleSJ9.KabES1DIR9hv6FEqbQ_G2jMm06Rb3pDxcKOLNDiC67q_9uoandGrWBg1hFNBDcp7_vkgXF77If9vGdYCSgoYXbRgVtESPPYdmY8v3tpl8LiayYR4a3PA9omvpeJIGLMe_YQoL3AeyIcaf2mLf5tTuRqIunbCJxTk1fWOGe5RddwZk44g2-0uZ69tm9YUM7cqL2r8I8stgYDB8d6OWSUeXzY3NxhQK_OH09qFQyHlO0UrpY0luyWgj2RrZiA3zMQYJMJLCJSOS-Ea1a-siSIyo_5LULi2JdRkbjJsMqNNY9AkH-vn4D7knBqRPSzwBxDymLhdt9MyEf1mh6L6UFeh3A", "https://api.divar.ir/v1/open-platform")
	kenarHandler := handlers.NewKenarHandler(kenarService)
	oauthHandler := handlers.NewOAuthHandler(oauthService)

	r := mux.NewRouter()

	r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	r.HandleFunc("/poi", kenarHandler.Poi)
	r.HandleFunc("/addon/oauth", oauthHandler.AddonOauth)
	r.HandleFunc("/oauth/callback", oauthHandler.OauthCallback)
	port := os.Getenv("PORT")
	log.Printf("Server started on port %s", port)
	http.ListenAndServe(":"+port, r)
}
