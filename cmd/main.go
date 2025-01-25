package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/YESUBZERO/MAGI-Scraper/internal/config"
	"github.com/YESUBZERO/MAGI-Scraper/internal/kafka"
	"github.com/YESUBZERO/MAGI-Scraper/internal/processor"
)

func main() {
	// Cargar configuración desde variables de entorno
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	// Contexto para manejar señales del sistema (Ctrl+C, etc.)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Capturar señales del sistema
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		<-signalChan
		log.Println("Cerrando servicio...")
		cancel()
	}()

	// Asegurar que el tópico exista
	err = kafka.EnsureTopicExists(cfg.Kafka.Brokers, cfg.Kafka.EnrichedTopic, 3, 1)
	if err != nil {
		log.Fatalf("Error asegurando el tópico '%s': %v", cfg.Kafka.EnrichedTopic, err)
	}

	// Crear productor Kafka
	producer, err := kafka.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.EnrichedTopic)
	if err != nil {
		log.Fatalf("Error inicializando productor Kafka: %v", err)
	}
	defer producer.Close()

	// Iniciar el consumidor Kafka
	err = kafka.StartKafkaConsumer(ctx, cfg.Kafka, func(message []byte) error {
		// Procesar mensajes y publicar en el tópico
		return processor.ProcessMessage(producer, message, cfg.Scraper)
	})
	if err != nil {
		log.Fatalf("Error iniciando consumidor Kafka: %v", err)
	}
}