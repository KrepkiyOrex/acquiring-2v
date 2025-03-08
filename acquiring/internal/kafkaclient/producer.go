package kafkaclient

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type ProducerInterface interface {
	SendPaymentProcessed(payment PaymentProcessed)
	Close()
}

type Producer struct {
	Client *kafka.Writer
}

type PaymentProcessed struct {
	OrderID       string `json:"order_id"`
	Status        string `json:"status"`
	TransactionID string `json:"transaction_id"`
	ProcessedAt   string `json:"processed_at"`
}

func NewProducer(broker string) (*Producer, error) {
	client := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   "payment_processed",
	})
	return &Producer{Client: client}, nil
}

func (p *Producer) SendPaymentProcessed(payment PaymentProcessed) {
	value, err := json.Marshal(payment)
	if err != nil {
		log.Printf("[Acquiring] Error marshalling payment: %v\n", err)
		return
	}

	err = p.Client.WriteMessages(context.Background(), kafka.Message{
		Value: value,
	})
	if err != nil {
		log.Printf("[Acquiring] Error sending message: %v\n", err)
		return
	}

	log.Printf("[Acquiring] Payment processed event sent")
}

func (p *Producer) Close() {
	p.Client.Close()
}
