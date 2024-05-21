package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/config"
	service "github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/server_status"
)

type ConsumerService struct {
    consumerGroup sarama.ConsumerGroup
}

func NewConsumerService(config *config.Config) *ConsumerService {
    // Set up Kafka consumer group configuration
    kafkaConfig := sarama.NewConfig()
    kafkaConfig.Consumer.Return.Errors = true
    log.Printf("Kafka brokers: %v\n", config.KAFKA.GroupID)
    // Create a new consumer group
    consumerGroup, err := sarama.NewConsumerGroup(config.KAFKA.Brokers, config.KAFKA.GroupID, kafkaConfig)
    if err != nil {
        log.Fatalf("Error creating consumer group: %v", err)
    } else {
        log.Println("Kafka Consumer Group created")
    }

    return &ConsumerService{consumerGroup}
}

func (cs *ConsumerService) ConsumerStart(servers *map[uint]service.Server, sigchan chan os.Signal) {
    // Consume messages from Kafka
    go func() {
        for {
            // Join the consumer group
            err := cs.consumerGroup.Consume(context.Background(), []string{"Server2"}, cs)
            if err != nil {
                log.Printf("Error from consumer group: %v\n", err)
            }

            // Check for errors
            if err != nil {
                log.Printf("Error from consumer group: %v\n", err)
            }
        }
    }()

    // Wait for signals
    <-sigchan
    log.Println("Received signal, shutting down consumer...")
}

func (cs *ConsumerService) Setup(sarama.ConsumerGroupSession) error {
    // Setup is called before ConsumeClaim
    return nil
}

func (cs *ConsumerService) Cleanup(sarama.ConsumerGroupSession) error {
    // Cleanup is called after ConsumeClaim
    return nil
}

func (cs *ConsumerService) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    // Consume messages from the claim
    for message := range claim.Messages() {
        // Decode the Kafka message into the Server struct
        if message.Value == nil {
            continue
        }

        if strings.Contains(string(message.Value), "drop") {
            var dropServer service.DropServer
            if err := json.Unmarshal(message.Value, &dropServer); err != nil {
                log.Printf("Error decoding message: %v\n", err)
                continue // Skip to the next message
            }
            // delete(*servers, dropServer.ID)
            log.Printf("Received message: %v\n", dropServer)
        } else {
            var server service.Server
            if err := json.Unmarshal(message.Value, &server); err != nil {
                log.Printf("Error decoding message: %v\n", err)
                continue // Skip to the next message
            }

            // (*servers)[server.ID] = server
            // log.Printf("Received message: %v\n", (*servers)[server.ID].IP)
            log.Printf("Received message: %v\n", server.IP)
        }

        // Mark message as processed
        session.MarkMessage(message, "")
    }
    return nil
}
