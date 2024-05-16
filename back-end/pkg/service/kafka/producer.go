package kafka

import (
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	service "github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/server_status"
)

type ProducerService struct {
	producer sarama.SyncProducer
}

func NewProducerService(config *config.Config) *ProducerService {
	// Set up Kafka producer configuration
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	// Create a new producer
	producer, err := sarama.NewSyncProducer(config.KAFKA.Brokers, kafkaConfig)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
		return nil
	} else {
		log.Println("Kafka Producer created")
	}

	return &ProducerService{producer}
}

func (ps *ProducerService) SendServer(id uint,server model.Server) {
	serverSend := service.Server{
		ID: id,
		IP: server.IP,
		Status: server.Status,
	}
	// Create a new Kafka message
	message, err := json.Marshal(serverSend)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
	}

	messageSend := &sarama.ProducerMessage{
		Topic: "Server",
		Value: sarama.ByteEncoder(message),
	}

	// Send the message
	partition, offset, err := ps.producer.SendMessage(messageSend)
	
	if err != nil {
		log.Printf("Error sending message:", err)
	} else {
		log.Printf("Message sent successfully, partition=%d, offset=%d\n", partition, offset)
	}

	
}

func (ps *ProducerService) DropServer(id uint) {
	dropMessage := service.DropServer{
		ID: id,
		Message: "drop",
	}
	// Create a new Kafka message
	message, err := json.Marshal(dropMessage)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
	}

	// Send the message to the Kafka topic
	messageSend := &sarama.ProducerMessage{
		Topic: "Server",
		Value: sarama.ByteEncoder(message),
	}

	// Send the message
	partition, offset, err := ps.producer.SendMessage(messageSend)
	
	if err != nil {
		log.Printf("Error sending message:", err)
	} else {
		log.Printf("Message sent successfully, partition=%d, offset=%d\n", partition, offset)
	}
}

