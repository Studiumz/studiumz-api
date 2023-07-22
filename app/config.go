package app

import (
	"fmt"
	"log"

	"github.com/cohere-ai/cohere-go"
	"github.com/spf13/viper"
)

type Config struct {
	Env  string `mapstructure:"ENV"`
	Port string `mapstructure:"PORT"`

	ClientAppUrl string `mapstructure:"CLIENT_APP_URL"`

	DbHost string `mapstructure:"DB_HOST"`
	DbPort string `mapstructure:"DB_PORT"`
	DbName string `mapstructure:"DB_NAME"`
	DbUser string `mapstructure:"DB_USER"`
	DbPwd  string `mapstructure:"DB_PWD"`
	DbDsn  string

	CohereApiKey string `mapstructure:"COHERE_API_KEY"`
	CohereClient *cohere.Client
}

func LoadConfig() (c Config) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to load config: ", err)
	}
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	c.setDBConfig()

	return c
}

func (c *Config) setDBConfig() {
	var ssl string
	if c.Env == "local" {
		ssl = "sslmode=disable"
	} else {
		//TODO: add ssl cert for dev and prod
		ssl = "sslmode=require"
	}

	c.DbDsn = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s %s", c.DbHost, c.DbPort, c.DbName, c.DbUser, c.DbPwd, ssl)
}
