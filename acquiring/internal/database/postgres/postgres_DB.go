package postgres

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBWraper struct {
	*gorm.DB
}

func (dbw *DBWraper) Close() error {
	sqlDB, err := dbw.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (dbw *DBWraper) Debug() *DBWraper {
	db := dbw.DB.Debug()
	return &DBWraper{DB: db}
}

func SetupDataBases() (*DBWraper, *DBWraper) {
	var wg sync.WaitGroup
	wg.Add(2)

	var db1, db2 *DBWraper
	var err1, err2 error

	go func() {
		defer wg.Done()
		db1, err1 = ConnectToDB("acquiring")
		if err1 != nil {
			log.Fatalf("[Acquiring] Failed to connect to the database acquiring: %v", err1)
		}
	}()

	go func() {
		defer wg.Done()
		db2, err2 = ConnectToDB("bank")
		if err2 != nil {
			log.Fatalf("[Acquiring] Failed to connect to the database bank: %v", err2)
		}
	}()

	wg.Wait()

	return db1, db2
}

// ConnectToDB функция для подключения к базе данных PostgreSQL
func ConnectToDB(dbName string) (*DBWraper, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbName, sslmode)

	// fmt.Println("PostgreSQL:", dsn) // for debug

	db, err := connectingWithRetry(dsn, 5, 5*time.Second)
	if err != nil {
		return nil, err
	}

	log.Printf("[Acquiring] === Database %v connection established ===", dbName)
	return db, nil
}

func connectingWithRetry(dsn string, maxRetries int, retryInterval time.Duration) (*DBWraper, error) {
	var db *gorm.DB
	var err error

	for i := 1; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			return &DBWraper{db}, nil
		}

		log.Printf("[Acquiring] Failed to connect to database (attempt %d/%d): %v\n", i, maxRetries, err)

		time.Sleep(retryInterval)
	}
	return nil, err
}
