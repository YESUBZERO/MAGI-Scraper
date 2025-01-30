package config

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Kafka   KafkaConfig
	Scraper ScraperConfig
}

type KafkaConfig struct {
	Brokers        []string `envconfig:"KAFKA_BROKERS" required:"true"`
	ProcessedTopic string   `envconfig:"KAFKA_CONSUMER_TOPIC" required:"true"`
	EnrichedTopic  string   `envconfig:"KAFKA_PRODUCER_TOPIC" required:"true"`
	GroupID        string   `envconfig:"KAFKA_GROUP_ID" required:"true"`
}

type ScraperConfig struct {
	UserAgent string `envconfig:"SCRAPER_USER_AGENT" default:"Mozilla/5.0 (Windows NT 10.0; Win64; x64)"`
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

// Descargar el certificado del servidor y guardarlo como "shipdb.pem"
func DownloadCertificate(host string, outputPath string) error {
	// Crear el comando para ejecutar OpenSSL
	cmd := exec.Command("sh", "-c",
		fmt.Sprintf("openssl s_client -showcerts -connect %s </dev/null 2>/dev/null | openssl x509 -outform PEM > %s", host, outputPath))

	// Establecer salida estándar y error estándar para depuración
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Ejecutar el comando
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error ejecutando OpenSSL: %w", err)
	}

	log.Printf("📜 Certificado descargado y guardado como %s", outputPath)
	return nil
}
