package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

// KafkaProducer envuelve la funcionalidad de un productor Kafka
type KafkaProducer struct {
	Producer sarama.SyncProducer
	Topic    string
}

// Nueva instancia de KafkaProducer
func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		Producer: producer,
		Topic:    topic,
	}, nil
}

// Publicar un mensaje en Kafka
func (kp *KafkaProducer) PublishMessage(message []byte) error {
	_, _, err := kp.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: kp.Topic,
		Value: sarama.ByteEncoder(message),
	})
	if err != nil {
		log.Printf("Error publicando mensaje enriquecido: %v", err)
		return err
	}
	log.Printf("Mensaje enriquecido publicado en el tópico '%s'", kp.Topic)
	return nil
}

// Cerrar el productor Kafka
func (kp *KafkaProducer) Close() {
	if err := kp.Producer.Close(); err != nil {
		log.Printf("Error cerrando el productor Kafka: %v", err)
	}
}
