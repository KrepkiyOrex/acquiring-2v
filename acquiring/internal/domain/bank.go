package domain

import "github.com/gofiber/fiber/v2"

type BankRepository interface {
	CreateUserCard(ctx *fiber.Ctx) error
	DeductFromAccount(ctx *fiber.Ctx) error
	AddFunds(ctx *fiber.Ctx) error
	GetAllCardDetails(ctx *fiber.Ctx) error
}
