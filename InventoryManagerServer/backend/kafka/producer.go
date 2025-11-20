package kafka

import (
	"fmt"
	"inventory-manager-server/config"
	"inventory-manager-server/dto"
	"log"

	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
)

var Producer sarama.SyncProducer

// InitProducer initializes Kafka producer
func InitProducer(brokers []string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	var err error
	Producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return fmt.Errorf("failed to create kafka producer: %w", err)
	}

	log.Println("Kafka producer initialized successfully")
	return nil
}

// PublishMessage publishes a message to Kafka
func PublishMessage(topic string, message interface{}) error {
	if Producer == nil {
		return fmt.Errorf("kafka producer not initialized")
	}
	messageDTO := dto.Message{
		SenderInstanceID: config.CONFIG.InstanceID,
		Payload:          message,
	}
	// Serialize the message to JSON
	payload, err := sonic.Marshal(messageDTO)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Construct Kafka message
	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(payload),
	}

	// Send the message
	_, _, err = Producer.SendMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("failed to send kafka message: %w", err)
	}

	return nil
}

// CloseProducer closes the producer
func CloseProducer() error {
	if Producer != nil {
		return Producer.Close()
	}
	return nil
}
