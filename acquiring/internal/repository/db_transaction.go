package repository

import (
	"net/http"
	"time"

	// "github.com/KrepkiyOrex/acquiring/internal/producer"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Transactions struct {
	TransactionID int64     `json:"transaction_id" gorm:"primaryKey"`
	OrderID       int64     `json:"orderId"`
	UserID        int64     `json:"userId"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"paymentMethod"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type TransRepos struct {
	DB *gorm.DB
}

func NewTransRepos(db *gorm.DB) *TransRepos {
	return &TransRepos{DB: db}
}

func NewTransaction() *Transactions {
	return &Transactions{}
}

func (tr *TransRepos) GetTransByID(ctx *fiber.Ctx, txID int64) error {
	transModel := NewTransaction()

	err := tr.DB.Where("transaction_id = ?", txID).First(transModel).Error
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the transactions"})
		return err
	}
	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "transactions id fetched successfully",
		"data":    transModel,
	})
	return nil
}

func (tr *TransRepos) GetTransactions(ctx *fiber.Ctx) error {
	transModels := &[]Transactions{}
	if err := tr.DB.Find(&transModels).Error; err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get transactions"})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "transactions fetched successfully",
		"data":    transModels,
	})
	return nil
}

func (tr *TransRepos) DeleteTransaction(ctx *fiber.Ctx, txID int64) error {
	transModel := NewTransaction()

	result := tr.DB.Delete(&transModel, txID)

	if result.Error != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not delete transaction"})
	}

	if result.RowsAffected == 0 {
		return ctx.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": "transaction not found"})
	}

	return ctx.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "transaction deleted successfully"})
}
