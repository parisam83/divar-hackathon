package utils

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	Kenar    KenarConfig
	Session  SessionConfig
	Server   ServerConfig
}
type KenarConfig struct {
	AppSlug     string `mapstructure:"KenarAppSlug"`
	ApiKey      string `mapstructure:"KenarApiKey"`
	OauthSecret string `mapstructure:"OauthSecret"`
	BaseURL     string `mapstructure:"BaseUrl"`
}
type ServerConfig struct {
	Port string `mapstructure:"Port"`
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
	AuthKey string `mapstructure:"SessionAuthKey"`
	// EncKey  string `mapstructure:"SessionEncKey"`
}

func (cfg *KenarConfig) Validate() error {

	if cfg.AppSlug == "" || cfg.OauthSecret == "" {
		return fmt.Errorf("missing required OAuth configurations")
	}
	return nil
}

func LoadConfig() (*Config, error) {

	err := godotenv.Load("/home/divar/Realestate-POI/utils/.env")
	if err != nil {
		log.Fatal("Error loading .env file" + err.Error())
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../Realestate-POI/configs")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// check if config file is not provided
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file")
	}
	for _, key := range viper.AllKeys() {
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
