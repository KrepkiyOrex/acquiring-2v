package integration

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/KrepkiyOrex/acquiring/internal/kafkaclient"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/mock"
)

type MockReader struct {
	mock.Mock
}

func (m *MockReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	args := m.Called(ctx)
	return args.Get(0).(kafka.Message), args.Error(1)
}

func (m *MockReader) Close() error {
	return nil
}

type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) SendPaymentProcessed(payment kafkaclient.PaymentProcessed) {
	m.Called(payment)
}

func (m *MockProducer) Close() {
	// заглушка для совместимости с интерфейсом
}

func TestConsumerOrders(t *testing.T) {
	mockReader := new(MockReader)
	mockProducer := new(MockProducer)

	order := kafkaclient.OrderCreated{
		OrderID:    "123",
		TotalPrice: 1000,
	}
	orderBytes, _ := json.Marshal(order)
	
	mockReader.On("ReadMessage", mock.Anything).Return(kafka.Message{Value: orderBytes}, nil).Once()
	mockReader.On("ReadMessage", mock.Anything).Return(kafka.Message{}, context.Canceled).Once()

	mockProducer.On("SendPaymentProcessed", mock.Anything).Return()

	consumer := &kafkaclient.Consumer{Client: mockReader}

	go consumer.ConsumerOrders(mockProducer)
	time.Sleep(100 * time.Millisecond)

	mockProducer.AssertCalled(t, "SendPaymentProcessed", mock.Anything)
}
