package handlers

import (
	"strconv"

	"github.com/KrepkiyOrex/acquiring/internal/repository"
	"github.com/KrepkiyOrex/acquiring/internal/service"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, svc *service.Service) {
	// AuthMiddleware() добавь в будущем для app.Group

	pages := app.Group("/pages")
	pages.Get("/payment", repository.ShowPaymentPage)
	pages.Get("/top_up", repository.ShowTopUpPage)
	pages.Get("/create_card", repository.ShowCreatCard)


	// ==============================================================
	api := app.Group("/api")
	
	api.Get("/get_allcard", svc.BankRepo.GetAllCardDetails)
	api.Get("/get_transactions", svc.TransactionRepo.GetTransactions)

	api.Get("/get_transaction/:id", func(ctx *fiber.Ctx) error {
		id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid transaction ID"})
		}
		return svc.TransactionRepo.GetTransByID(ctx, id)
	})
	

	api.Delete("/delete_transaction/:id", func(ctx *fiber.Ctx) error {
		id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid transaction ID"})
		}
		return svc.TransactionRepo.DeleteTransaction(ctx, id)
	})

	api.Post("/new_card", svc.BankRepo.CreateUserCard)
	api.Post("/add_funds", svc.BankRepo.AddFunds)

	paymentProcessor := repository.NewPaymentProcessor(svc)
	api.Post("/process_payment", paymentProcessor.ProcessPayment)
}
