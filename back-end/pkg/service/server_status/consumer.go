package service

import (
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
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

func(cs *ConsumerService)ConsumerStart(sigchan chan os.Signal) {
    // Create a new partition consumer for the given topic
    partitionConsumer, err := cs.consumer.ConsumePartition("test-topic3", 0, sarama.OffsetNewest)
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
            log.Printf("Received message: key='%s' value='%s'\n", string(msg.Key), string(msg.Value))
        case err := <-partitionConsumer.Errors():
            log.Printf("Error consuming message: %v\n", err)
        }
    }
}