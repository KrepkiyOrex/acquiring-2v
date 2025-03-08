package kafkaclient

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type OrderCreated struct {
	OrderID    string `json:"order_id"`
	TotalPrice int    `json:"total_price"`
}

type Reader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

type Consumer struct {
	Client Reader
}

func NewConsumer(config kafka.ReaderConfig) (*Consumer, error) {
	client := kafka.NewReader(config)
	return &Consumer{Client: client}, nil
}

func (c *Consumer) Close() {
	c.Client.Close()
}

func (c *Consumer) ConsumerOrders(p ProducerInterface) {
	for {
		msg, err := c.Client.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error receiving message: %v\n", err)
			if err == context.Canceled {
				break
			}
			continue
		}

		var order OrderCreated
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Error unmarshalling message: %v\n", err)
			continue
		}

		log.Printf("Received order: %+v\n", order)

		// имитация обработки платежа
		payment := PaymentProcessed{
			OrderID:       order.OrderID,
			Status:        "success",
			TransactionID: "tx789",
			ProcessedAt:   time.Now().Format(time.RFC3339),
		}
		
		p.SendPaymentProcessed(payment)
	}
}
