package kafka

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	service "github.com/mxngocqb/VCS-SERVER/back-end/pkg/service/server_status"
	"github.com/stretchr/testify/assert"
)

func TestConsumerService_ConsumerStart(t *testing.T) {
    // Mock Kafka consumer
    consumerMock := mocks.NewConsumer(t, nil)
    consumerService := &ConsumerService{consumer: consumerMock}

    // Set up servers map and signal channel
    servers := make(map[uint]service.Server)
    sigchan := make(chan os.Signal, 1)

    // Create a Kafka message
    serverJSON, _ := json.Marshal(service.Server{ID: 1, IP: "192.168.1.1"})
    msg := &sarama.ConsumerMessage{
        Value: serverJSON,
    }

    // Add message to mock consumer
    consumerMock.ExpectConsumePartition("Server", 0, sarama.OffsetNewest).YieldMessage(msg)

    // Run consumer start function
    go consumerService.ConsumerStart(&servers, sigchan)

    // Simulate receiving signal to stop
    sigchan <- os.Interrupt

    // Verify server has been added to the map
    assert.Equal(t, "192.168.1.1", servers[1].IP)

    // Close mock consumer
    consumerMock.Close()
}
