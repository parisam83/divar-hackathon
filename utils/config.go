package utils

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Database      DatabaseConfig
	App           AppConfig
	SessionConfig SessionConfig
}
type AppConfig struct {
	AppSlug     string `mapstructure:"KENAR_APP_SLUG"`
	ApiKey      string `mapstructure:"KENAR_API_KEY"`
	OauthSecret string `mapstructure:"OAUTH_SECRET"`
	SessionKey  string `mapstructure:"OAUTH_SESSION_KEY"`
	ServerPort  string `mapstructure:"PORT"`
}

type DatabaseConfig struct {
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
type SessionConfig struct {
	AuthKey string
	EncKey  string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("/home/divar/Realestate-POI/utils/.env")
	if err != nil {
		log.Fatal("Error loading .env file" + err.Error())
	}
	viper.SetConfigName("db")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../Realestate-POI/configs")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// check if config file is not provided
	viper.SetConfigName("db")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../Realestate-POI/configs")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file")
	}
	for _, key := range viper.AllKeys() {
		// fmt.Println(key)
		value := viper.GetString(key)
		expanded := os.ExpandEnv(value)
		viper.Set(key, expanded)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	// log.Printf("Config Loaded: %+v", config)
	return &config, nil
}
