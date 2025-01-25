package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

// Asegurar que el tópico exista
func EnsureTopicExists(brokers []string, topic string, partitions, replicationFactor int) error {
	admin, err := sarama.NewClusterAdmin(brokers, nil)
	if err != nil {
		return err
	}
	defer admin.Close()

	topics, err := admin.ListTopics()
	if err != nil {
		return err
	}

	if _, exists := topics[topic]; exists {
		log.Printf("El tópico '%s' ya existe", topic)
		return nil
	}

	err = admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     int32(partitions),
		ReplicationFactor: int16(replicationFactor),
	}, false)
	if err != nil {
		return err
	}
	log.Printf("Tópico '%s' creado con éxito", topic)
	return nil
}
