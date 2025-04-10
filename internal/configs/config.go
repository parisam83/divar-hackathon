package configs

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
	Jwt      JWTConfig
	Server   ServerConfig
	Neshan   NeshanConfig
	Snapp    SnappConfig
	Tapsi    TapsiConfig
}

type NeshanConfig struct {
	NeshanApiKey string `mapstructure:"NeshanApiKey"`
}

type KenarConfig struct {
	AppSlug     string `mapstructure:"KenarAppSlug"`
	ApiKey      string `mapstructure:"KenarApiKey"`
	OauthSecret string `mapstructure:"OauthSecret"`
	BaseURL     string `mapstructure:"BaseUrl"`
	OpenPlatformApi string `mapstructure:"OpenPlatformApi"`
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

type JWTConfig struct {
	JwtSecret string `mapstructure:"JwtSecret"`
}

type SnappConfig struct {
	ApiKey        string `mapstructure:"access_token"`
	Clck          string `mapstructure:"clck"`
	Clsk          string `mapstructure:"clsk"`
	GA            string `mapstructure:"ga"`
	GATracking    string `mapstructure:"ga_tracking"`
	YandexDate    string `mapstructure:"ym_d"`
	YandexAd      string `mapstructure:"ym_isad"`
	YandexUID     string `mapstructure:"ym_uid"`
	CookieSession string `mapstructure:"cookie_session"`
}
type TapsiConfig struct {
	Clck         string `mapstructure:"clck"`
	Clsk         string `mapstructure:"clsk"`
	AccessToken  string `mapstructure:"access_token"`
	RefreshToken string `mapstructure:"refresh_token"`
}

func (cfg *KenarConfig) Validate() error {
	if cfg.AppSlug == "" || cfg.OauthSecret == "" || cfg.OpenPlatformApi == "" {
		return fmt.Errorf("missing required OAuth configurations")
	}
	return nil
}

func (cfg *Config) Validate() error {
	// Database validation
    if cfg.Database.Host == "" {
        return fmt.Errorf("database host is empty")
    }
    if cfg.Database.Port == 0 {
        return fmt.Errorf("database port is zero")
    }
    if cfg.Database.Username == "" {
        return fmt.Errorf("database username is empty")
    }
    if cfg.Database.Password == "" {
        return fmt.Errorf("database password is empty")
    }
    if cfg.Database.DBName == "" {
        return fmt.Errorf("database name is empty")
    }

	// TODO: Clean this part
    // Kenar validation
    if err := cfg.Kenar.Validate(); err != nil {
        return fmt.Errorf("kenar config validation failed: %w", err)
    }

    // Server validation
    if cfg.Server.Port == "" {
        return fmt.Errorf("server port is empty")
    }

    // Session validation
    if cfg.Session.AuthKey == "" {
        return fmt.Errorf("session auth key is empty")
    }

    // JWT validation
    if cfg.Jwt.JwtSecret == "" {
        return fmt.Errorf("jwt secret is empty")
    }

    // Neshan validation
    if cfg.Neshan.NeshanApiKey == "" {
        return fmt.Errorf("neshan api key is empty")
    }

    // Snapp validation
    if cfg.Snapp.ApiKey == "" {
        return fmt.Errorf("snapp api key is empty")
    }

    // Tapsi validation
    if cfg.Tapsi.AccessToken == "" {
        return fmt.Errorf("tapsi access token is empty")
    }

    return nil
}

func LoadConfig() (*Config, error) {
	if os.Getenv("ENV") == "development" {
		err := godotenv.Load("./internal/configs/.env")
		if err != nil {
			log.Printf("Error loading .env file in LoadConfig: %v", err)
			return nil, err
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/configs")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// check if config file is not provided
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
		return nil, err
	}

	for _, key := range viper.AllKeys() {
		value := viper.GetString(key)
		expanded := os.ExpandEnv(value)
		viper.Set(key, expanded)
	}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Printf("Error unmarshalling config file: %v", err)
		return nil, err
	}

	// Validate the configuration
    if err := config.Validate(); err != nil {
        log.Printf("Configuration validation failed: %v", err)
        return nil, err
    }

	// log.Printf("Config Loaded: %+v", config)
	return &config, nil
}
