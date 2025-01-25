package processor

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/YESUBZERO/MAGI-Scraper/internal/config"
	"github.com/YESUBZERO/MAGI-Scraper/internal/kafka"
	"github.com/YESUBZERO/MAGI-Scraper/internal/parsership"
	"github.com/YESUBZERO/MAGI-Scraper/internal/scraper"
	"golang.org/x/net/context"
)

type AISMessage struct {
	MsgType int `json:"msg_type"`
	IMO     int `json:"imo"`
}

// Worker procesa mensajes de un canal y publica mensajes en Kafka
func Worker(ctx context.Context, wg *sync.WaitGroup, messageChan <-chan []byte, producer *kafka.KafkaProducer, scraperConfig config.ScraperConfig) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Println("Worker finalizado.")
			return
		case message, ok := <-messageChan:
			if !ok {
				return
			}

			// Procesar mensaje
			if err := ProcessMessage(producer, message, scraperConfig); err != nil {
				log.Printf("Error procesando mensaje: %v", err)
			}

			// Simular un retraso
			time.Sleep(time.Duration(scraperConfig.Delay) * time.Second)
		}
	}
}

// Procesar un mensaje Kafka
func ProcessMessage(producer *kafka.KafkaProducer, message []byte, scraperConfig config.ScraperConfig) error {
	var aisMessage AISMessage

	// Deserializar el mensaje
	if err := json.Unmarshal(message, &aisMessage); err != nil {
		return err
	}

	// Filtrar mensajes relevantes
	if aisMessage.MsgType == 5 || aisMessage.MsgType == 24 {
		if aisMessage.IMO != 0 {
			log.Printf("Procesando scraping con IMO: %d", aisMessage.IMO)

			// Logica del Scraper
			shipData := ScrapeHandler(aisMessage.IMO, scraperConfig)

			if shipData != "" {
				// Publicar mensaje procesado
				log.Println(shipData)
				processedMessage, err := json.Marshal(shipData)
				if err != nil {
					return err
				}
				return producer.PublishMessage(processedMessage)
			}
		}
	}
	return nil
}

func ScrapeHandler(imo int, scraperConfig config.ScraperConfig) string {
	// Realizar el scraping del código del buque
	rawHTML, err := scraper.ScrapeShipCode(imo, scraperConfig.UserAgent)
	if err != nil {
		log.Printf("Error al obtener HTML para IMO %d: %v", imo, err)
	}

	// Parsear el url del buque
	address, err := parsership.ParserShipCode(rawHTML)
	if err != nil {
		log.Printf("Error procesando URL para IMO %d: %v", imo, err)
	}

	// Realizar el scraping de las caracteristicas del buque
	rawData, err := scraper.ScrapeData(address, scraperConfig.UserAgent)
	if err != nil {
		log.Printf("Error al obtener HTML para URL %s: %v", address, err)
	}

	// Parsear la data de las caracteristicas del buque
	shipData, err := parsership.ParserShipData(rawData)
	if err != nil {
		log.Printf("Error procesando HTML para URL %s: %v", address, err)
	}

	return shipData
}
