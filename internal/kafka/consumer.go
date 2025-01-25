package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
	"github.com/YESUBZERO/MAGI-Scraper/internal/config"
)

// KafkaHandler es un manejador de mensajes Kafka
type KafkaHandler struct {
	MessageHandler func(message []byte) error
}

// Implementación de la interfaz sarama.ConsumerGroupHandler
func (h *KafkaHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *KafkaHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *KafkaHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if err := h.MessageHandler(message.Value); err != nil {
			log.Printf("Error al procesar mensaje: %v", err)
		}
		session.MarkMessage(message, "")
	}
	return nil
}

// Iniciar el consumidor Kafka
func StartKafkaConsumer(ctx context.Context, Config config.KafkaConfig, messageHandler func(message []byte) error) error {
	consumerGroup, err := sarama.NewConsumerGroup(Config.Brokers, Config.GroupID, nil)
	if err != nil {
		return err
	}
	defer consumerGroup.Close()

	handler := &KafkaHandler{MessageHandler: messageHandler}

	// Ciclo principal del consumidor
	go func() {
		for {
			if ctx.Err() != nil {
				return
			}
			if err := consumerGroup.Consume(ctx, []string{Config.ProcessedTopic}, handler); err != nil {
				log.Printf("Error al consumir mensaje: %v", err)
			}
		}
	}()

	log.Printf("Consumidor Kafka iniciado para el topico: %s", Config.ProcessedTopic)
	<-ctx.Done()
	return nil
}
