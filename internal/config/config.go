package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServHost   string
	ServPort   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	MIURL      string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		ServHost:   os.Getenv("SERV_HOST"),
		ServPort:   os.Getenv("SERV_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		MIURL:      os.Getenv("MI_URL"),
	}
}

func (c *Config) DBConnectionString() string {
	return "host=" + c.DBHost + " port=" + c.DBPort + " user=" + c.DBUser + " password=" + c.DBPassword + " dbname=" + c.DBName + " sslmode=disable"
}
