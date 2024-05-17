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
	logger *log.Logger
}

func NewProducerService(config *config.Config, logger *log.Logger) *ProducerService {
	// Set up Kafka producer configuration
	sarama.Logger = log.New(logger.Writer(), "[Sarama] ", log.LstdFlags)
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	// Create a new producer
	producer, err := sarama.NewSyncProducer(config.KAFKA.Brokers, kafkaConfig)
	if err != nil {
		logger.Fatalf("Error creating producer: %v", err)
		log.Println("Error creating producer")
		return nil
	} else {
		logger.Println("Kafka Producer created")
		log.Println("Kafka Producer created")
	}

	return &ProducerService{producer, logger}
}

func (ps *ProducerService) SendServer(id uint,server model.Server)  error{
	serverSend := service.Server{
		ID: id,
		IP: server.IP,
		Status: server.Status,
	}
	// Create a new Kafka message
	message, err := json.Marshal(serverSend)
	if err != nil {
		ps.logger.Printf("Error marshalling message: %v", err)
		log.Printf("Error marshalling message: %v", err)
		return err
	}

	messageSend := &sarama.ProducerMessage{
		Topic: "Server",
		Value: sarama.ByteEncoder(message),
	}

	// Send the message
	partition, offset, err := ps.producer.SendMessage(messageSend)
	
	if err != nil {
		ps.logger.Printf("Error sending message:", err)
		log.Printf("Error sending message:", err)
		return err
	} else {
		ps.logger.Printf("Message sent successfully, partition=%d, offset=%d\n", partition, offset)
	}
	return nil
}

func (ps *ProducerService) DropServer(id uint) error {
	dropMessage := service.DropServer{
		ID: id,
		Message: "drop",
	}
	// Create a new Kafka message
	message, err := json.Marshal(dropMessage)
	if err != nil {
		ps.logger.Printf("Error marshalling message: %v", err)
		return err
	}

	// Send the message to the Kafka topic
	messageSend := &sarama.ProducerMessage{
		Topic: "Server",
		Value: sarama.ByteEncoder(message),
	}

	// Send the message
	partition, offset, err := ps.producer.SendMessage(messageSend)
	
	if err != nil {
		ps.logger.Printf("Error sending message:", err)
		return err
	} else {
		ps.logger.Printf("Message sent successfully, partition=%d, offset=%d\n", partition, offset)
		log.Printf("Message sent successfully, partition=%d, offset=%d\n", partition, offset)
	}
	return nil
}

