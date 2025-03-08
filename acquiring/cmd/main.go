package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/segmentio/kafka-go"

	"github.com/KrepkiyOrex/acquiring/handlers"
	"github.com/KrepkiyOrex/acquiring/internal/database/postgres"
	"github.com/KrepkiyOrex/acquiring/internal/kafkaclient"
	"github.com/KrepkiyOrex/acquiring/internal/repository"
	"github.com/KrepkiyOrex/acquiring/internal/service"
)

func main() {
	log.Println("[Acquiring] ========== Docker acquiring started ... ===========")

	db1, db2 := postgres.SetupDataBases()

	defer db1.Close()
	defer db2.Close()

	if err := db1.AutoMigrate(&repository.Transactions{}); err != nil {
		log.Fatalf("[Acquiring] Could not migrate the database: %v", err)
	}

	if err := db2.AutoMigrate(&repository.CardData{}); err != nil {
		log.Fatalf("[Acquiring] Could not migrate the database: %v", err)
	}
	
	
	db1 = db1.Debug()
	
	// for prod (work)
	// engine := html.New("./app/internal/source", ".html") // шаблонизатор

	// for dev (CompileDaemon)
	engine := html.New("/app/internal/source", ".html") // шаблонизатор

	// подключил шаблонизатор
	app := fiber.New(fiber.Config{Views: engine})

	// for prod (work)
	// app.Static("/", "./app/internal/source") // for static files

	// for dev (CompileDaemon)
	app.Static("/", "/app/internal/source") // for static files

	transactionRego := repository.NewTransRepos(db1.DB) // DB *gorm.DB
	bankRepos := repository.NewBankRepos(db2.DB)        // DB *gorm.DB

	svc := service.NewService(bankRepos, transactionRego)

	handlers.SetupRoutes(app, svc)

	go func() {
		if err := app.Listen(":8081"); err != nil {
			log.Fatalf("[Acquiring] Error starting server: %v", err)
		}
	}()

	StartKafkaProcessing()
}


func StartKafkaProcessing() {
	Brokers := "kafka:9092"

	consumerConfig := kafka.ReaderConfig{
		Brokers: []string{Brokers},
		GroupID: "acquiring_group",
		Topic:   "order_created",
	}

	consumer, err := kafkaclient.NewConsumer(consumerConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	producer, err := kafkaclient.NewProducer(Brokers)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	go consumer.ConsumerOrders(producer)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan

	log.Println("[Acquiring] Shutting down Kafka processing...")
}
