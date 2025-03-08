package domain

import "github.com/gofiber/fiber/v2"

type TransactionRepository interface {
	GetTransByID(ctx *fiber.Ctx, transactionID int64) error
	GetTransactions(ctx *fiber.Ctx) error
	DeleteTransaction(ctx *fiber.Ctx, transactionID int64) error
}
