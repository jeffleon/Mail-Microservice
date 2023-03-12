package config

import (
	"log"

	"github.com/spf13/viper"
)

var EnvConfigs *envConfigs

// We will call this in main.go to load the env variables
func InitEnvConfigs() {
	EnvConfigs = loadEnvVariables()
}

type envConfigs struct {
	AppEnvironment   string `mapstructure:"APP_ENV"`
	AppPort          string `mapstructure:"APP_PORT"`
	EmailHost        string `mapstructure:"EMAIL_HOST" dafault:"prueba"`
	EmailPort        int    `mapstructure:"EMAIL_PORT"`
	EmailDomain      string `mapstructure:"EMAIL_DOMAIN"`
	EmailPassword    string `mapstructure:"EMAIL_PASSWORD"`
	EmailUserName    string `mapstructure:"EMAIL_USERNAME"`
	EmailEncription  string `mapstructure:"EMAIL_ENCRIPTION"`
	EmailFromName    string `mapstructure:"EMAIL_FROM_NAME"`
	EmailFromAddress string `mapstructure:"EMAIL_FROM_ADDRESS"`
	RPCPort          string `mapstructure:"RPC_PORT"`
	KafkaHost        string `mapstructure:"KAFKA_HOST" default:"localhost"`
	KafkaUsername    string `mapstructure:"KAFKA_USERNAME"`
	KafkaPassword    string `mapstructure:"KAFKA_PASSWORD"`
	KafkaUserTopic   string `mapstructure:"KAFKA_USER_TOPIC"`
	KafkaPort        string `mapstructure:"KAFKA_PORT" deafault:"9092"`
}

func loadEnvVariables() (config *envConfigs) {
	// Tell viper the path/location of your env file. If it is root just add "."
	viper.AddConfigPath(".")

	// Tell viper the name of your file
	viper.SetConfigName(".env")

	// Tell viper the type of your file
	viper.SetConfigType("env")

	// Viper reads all the variables from env file and log error if any found
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	// Viper unmarshals the loaded env varialbes into the struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	return
}
