package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type config struct {
	DatabaseDriver string
	DatabaseUrl    string
	PORT           string
	KafkaBrokers   []string
	KafkaEnabled   bool
}

func LoadConfig() *config {
	err := godotenv.Load()
	if err != nil {
		log.Panicln("not able to find the .env file")
	}

	cfg := &config{
		DatabaseDriver: os.Getenv("dbDriver"),
		DatabaseUrl:    os.Getenv("dbSource"),
		PORT:           os.Getenv("port"),
		KafkaEnabled:   os.Getenv("kafkaenabled") == "true",
	}

	brokers := os.Getenv("kafkabrokers")
	if brokers != "" {
		cfg.KafkaBrokers = strings.Split(brokers, ",")
	} else {
		cfg.KafkaBrokers = []string{"localhost:9092"}
	}

	if cfg.DatabaseDriver == "" {
		cfg.DatabaseDriver = "postgres"
		log.Println("No database driver specified, using default: postgres")
	}
	if cfg.DatabaseUrl == "" {
		log.Panicln("DATABASE_URL is required")
	}
	if cfg.PORT == "" {
		cfg.PORT = "8080"
		log.Println("No PORT specified, using default: 8080")
	}
	return cfg

}
