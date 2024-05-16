package service

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	service "github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/server_status"
)

type ConsumerService struct {
	consumer sarama.Consumer
}

func NewConsumerSevice(config *config.Config) *ConsumerService {
	// Set up Kafka consumer configuration
    kafkaConfig := sarama.NewConfig()
    kafkaConfig.Consumer.Return.Errors = true
    // Create a new consumer
    consumer, err := sarama.NewConsumer(config.KAFKA.Brokers, kafkaConfig)
    if err != nil {
        log.Fatalf("Error creating consumer: %v", err)
    } else{
        log.Println("Kafka Consumer created")
    }

	return &ConsumerService{consumer}
}

func(cs *ConsumerService)ConsumerStart(servers *map[uint]service.Server, sigchan chan os.Signal) {
    // Create a new partition consumer for the given topic
    partitionConsumer, err := cs.consumer.ConsumePartition("Server", 0, sarama.OffsetNewest)
    if err != nil {
        log.Fatalf("Error creating partition consumer: %v", err)
    }

    defer func() {
        if err := partitionConsumer.Close(); err != nil {
            log.Fatalf("Error closing partition consumer: %v", err)
        }
    }()

    // Consume messages from Kafka
    for {
        select {
        case <-sigchan:
            log.Println("Received signal, shutting down consumer...")
            return
        case msg := <-partitionConsumer.Messages():
            // Decode the Kafka message into the Server struct
            if string(msg.Value) == "null" {
                continue
            }

            var server service.Server
            if err := json.Unmarshal(msg.Value, &server); err != nil {
                log.Printf("Error decoding message: %v\n", err)
                continue // Skip to the next message
            }

            (*servers)[server.ID] = server
            log.Printf("Received message: %v\n", (*servers)[server.ID].IP)
            
        case err := <-partitionConsumer.Errors():
            log.Printf("Error consuming message: %v\n", err)
        }
    }
}