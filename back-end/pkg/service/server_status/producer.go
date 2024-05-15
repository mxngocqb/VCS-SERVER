package service

import (
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
)

type ProducerService struct {
	producer sarama.AsyncProducer
}

func NewProducerService(config *config.Config) *ProducerService {
	// Set up Kafka producer configuration
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	// Create a new producer
	producer, err := sarama.NewAsyncProducer(config.KAFKA.Brokers, kafkaConfig)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
		return nil
	} else {
		log.Println("Kafka Producer created")
	}

	return &ProducerService{producer}
}

func (ps *ProducerService) SendServer(id uint,server model.Server) {
	serverSend := Server{
		ID: id,
		IP: server.IP,
		Status: server.Status,
	}
	// Create a new Kafka message
	message, err := json.Marshal(serverSend)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
	}

	// Send the message to the Kafka topic
	ps.producer.Input() <- &sarama.ProducerMessage{
		Topic: "test-topic3",
		Value: sarama.ByteEncoder(message),
	}
}

func (ps *ProducerService) DropServer(id uint) {
	// Send the message to the Kafka topic
	ps.producer.Input() <- &sarama.ProducerMessage{
		Topic: "test-topic3",
		Value: sarama.ByteEncoder("null"),
	}
}

