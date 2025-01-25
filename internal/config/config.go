package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Kafka   KafkaConfig
	Scraper ScraperConfig
}

type KafkaConfig struct {
	Brokers        []string `envconfig:"KAFKA_BROKERS" required:"true"`
	ProcessedTopic string   `envconfig:"KAFKA_PROCESSED_TOPIC" required:"true"`
	EnrichedTopic  string   `envconfig:"KAFKA_ENRICHED_TOPIC" required:"true"`
	GroupID        string   `envconfig:"KAFKA_GROUP_ID" required:"true"`
}

type ScraperConfig struct {
	UserAgent string `envconfig:"SCRAPER_USER_AGENT" default:"Mozilla/5.0"`
	Delay     int    `envconfig:"SCRAPER_DELAY" required:"true"`
}

// Carga la configuración desde las variables de entorno
func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
