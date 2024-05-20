package kafka

import (
	"errors"
	"testing"

	"github.com/Shopify/sarama/mocks"
	"github.com/mxngocqb/VCS-SERVER/back-end/internal/model"
	"github.com/mxngocqb/VCS-SERVER/back-end/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestProducerService_SendServer(t *testing.T) {
	kafkaLogger := util.KafkaLogger()
	// Create a mock Kafka producer
	producerMock := mocks.NewSyncProducer(t, nil)

	// Create a ProducerService with the mock producer
	producerService := &ProducerService{producer: producerMock, logger: kafkaLogger}

	// Define a server to send
	server := model.Server{
		IP:"192.168.1.1",
		Status: true,
	}

	// Expect SendMessage to be called once and succeed
	producerMock.ExpectSendMessageAndSucceed()

	// Send the server
	err := producerService.SendServer(1, server)

	// Ensure no errors occurred
	assert.Nil(t, err)
}

func TestProducerService_DropServer(t *testing.T) {
	kafkaLogger := util.KafkaLogger()
	// Create a mock Kafka producer
	producerMock := mocks.NewSyncProducer(t, nil)

	// Create a ProducerService with the mock producer
	producerService := &ProducerService{producer: producerMock, logger: kafkaLogger}

	// Expect SendMessage to be called once and succeed
	producerMock.ExpectSendMessageAndSucceed()

	// Drop the server
	err := producerService.DropServer(1)

	// Ensure no errors occurred
	assert.Nil(t, err)
}

func TestProducerService_SendServer_Error(t *testing.T) {
	kafkaLogger := util.KafkaLogger()
	// Create a mock Kafka producer
	producerMock := mocks.NewSyncProducer(t, nil)

	// Create a ProducerService with the mock producer
	producerService := &ProducerService{producer: producerMock, logger: kafkaLogger}


	// Define a server to send
	server := model.Server{
		IP:"192.168.1.1",
		Status: true,
	}

	// Expect SendMessage to be called once and fail with a specific error
	producerMock.ExpectSendMessageAndFail(errors.New("expected error message"))

	// Send the server
	err := producerService.SendServer(1, server)

	// Ensure the expected error is returned
	assert.EqualError(t, err, "expected error message")
}
