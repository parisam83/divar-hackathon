package config

// import (
// 	"log"
// 	"strings"

// 	"github.com/spf13/viper"
// )

// type Config struct {
// 	NeshanApiKey string `mapstructure:"neshan_api_key"`
// }

// var ConfigVar Config

// func init() {
// 	viper.AutomaticEnv()
// 	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

// 	if err := viper.Unmarshal(&ConfigVar); err != nil {
// 		log.Fatal(err)
// 	}
// }
