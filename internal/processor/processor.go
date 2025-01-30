package processor

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/YESUBZERO/MAGI-Scraper/internal/config"
	"github.com/YESUBZERO/MAGI-Scraper/internal/kafka"
	"github.com/YESUBZERO/MAGI-Scraper/internal/models"
	"github.com/YESUBZERO/MAGI-Scraper/internal/parsership"
	"github.com/YESUBZERO/MAGI-Scraper/internal/scraper"
	"golang.org/x/net/context"
)

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
				log.Println(err)
			}

			// Simular un retraso
			time.Sleep(time.Duration(scraperConfig.Delay) * time.Second)
		}
	}
}

// Procesar un mensaje Kafka
func ProcessMessage(producer *kafka.KafkaProducer, message []byte, scraperConfig config.ScraperConfig) error {
	var aisMessage models.AISMessage

	// Deserializar el mensaje
	if err := json.Unmarshal(message, &aisMessage); err != nil {
		return err
	}

	// Filtrar mensajes relevantes
	if aisMessage.MsgType == 5 || aisMessage.MsgType == 24 {
		if aisMessage.IMO != 0 {
			//log.Printf("Procesando scraping con IMO: %d", aisMessage.IMO)

			// Logica del Scraper
			shipData, err := ScrapeHandler(aisMessage.IMO, scraperConfig)
			if err != nil {
				return fmt.Errorf("❌ error al obtener datos del buque con IMO %d: %v", aisMessage.IMO, err)
			}

			if strings.TrimSpace(shipData) != "" && strings.TrimSpace(shipData) != "{\n}" || strings.TrimSpace(shipData) != "{" {
				// Deserializar mensaje procesado
				var scrapedShip models.Ship
				if err := json.Unmarshal([]byte(shipData), &scrapedShip); err != nil {
					return err
				}

				// Crear mensaje enriquecido
				enrichedMessage := models.Ship{
					IMO:            aisMessage.IMO,
					MMSI:           aisMessage.MMSI,
					Callsign:       aisMessage.CALLSIGN,
					Shipname:       aisMessage.SHIPNAME,
					ShipType:       aisMessage.SHIP_TYPE,
					BuiltYear:      scrapedShip.BuiltYear,
					Shipyard:       scrapedShip.Shipyard,
					HullNumber:     scrapedShip.HullNumber,
					KeelLaying:     scrapedShip.KeelLaying,
					LaunchDate:     scrapedShip.LaunchDate,
					DeliveryDate:   scrapedShip.DeliveryDate,
					GT:             scrapedShip.GT,
					NT:             scrapedShip.NT,
					CarryingCapTDW: scrapedShip.CarryingCapTDW,
					LengthOverall:  scrapedShip.LengthOverall,
					Breadth:        scrapedShip.Breadth,
					Depth:          scrapedShip.Depth,
					Propulsion:     scrapedShip.Propulsion,
					Power:          scrapedShip.Power,
					Screws:         scrapedShip.Screws,
					Speed:          scrapedShip.Speed,
				}

				// Publicar mensaje procesado
				processedMessage, err := json.Marshal(enrichedMessage)
				if err != nil {
					return err
				}
				return producer.PublishMessage(processedMessage, enrichedMessage.IMO)
			}
		}
	}
	return nil
}

func ScrapeHandler(imo int, scraperConfig config.ScraperConfig) (string, error) {
	// Realizar el scraping del código del buque
	rawHTML, err := scraper.ScrapeShipCode(imo, scraperConfig.UserAgent)
	if err != nil {
		//log.Printf("Error al obtener HTML para IMO %d: %v", imo, err)
		return "", err
	}

	// Parsear el url del buque
	address, _ := parsership.ParserShipCode(rawHTML)

	// Realizar el scraping de las caracteristicas del buque
	rawData, err := scraper.ScrapeData(address, scraperConfig.UserAgent)
	if err != nil {
		//log.Printf("Error al obtener HTML para URL %s: %v", address, err)
		return "", err
	}

	// Parsear la data de las caracteristicas del buque
	shipData, _ := parsership.ParserShipData(rawData)

	return shipData, nil
}
