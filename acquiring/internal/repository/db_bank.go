package repository

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/KrepkiyOrex/acquiring/internal/crypto"
	"github.com/KrepkiyOrex/acquiring/internal/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CardData struct {
	ID                  uint    `json:"id" gorm:"primaryKey"`
	Balance             float64 `json:"balance"`
	EncryptedCardNumber string  `json:"encryptedCardNumber"`
	EncryptedExpiryDate string  `json:"encryptedExpiryDate"`
	EncryptedCVV        string  `json:"encryptedCVV"`
	EncryptedCardName   string  `json:"encryptedCardName"`
}

type BankRepos struct {
	DB *gorm.DB
}

func NewBankRepos(db *gorm.DB) *BankRepos {
	return &BankRepos{DB: db}
}

func newCard() *CardData {
	return &CardData{}
}

func ShowCreatCard(ctx *fiber.Ctx) error {
	return ctx.Render("create-card", fiber.Map{
		"ErrorMessage": "",
	})
}

func (bank *BankRepos) CreateUserCard(ctx *fiber.Ctx) error {
	userCard := newCard()

	if err := ctx.BodyParser(userCard); err != nil {
		return ctx.Render("create-card", fiber.Map{
			"ErrorMessage": "Failed to parse request",
		})
	}

	if err := crypto.ProcessEncrypt((*crypto.CardData)(userCard)); err != nil {
		return fmt.Errorf("Encryption error: %v", err)
	}

	fmt.Println("User card: ", userCard)
	fmt.Println("=======================")

	if err := bank.DB.Create(&userCard).Error; err != nil {
		return ctx.Render("add-funds", fiber.Map{
			"message": "Cound not create card"})
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{"message": "Card user successfully created"})
}

func ShowPaymentPage(ctx *fiber.Ctx) error {
	return ctx.Render("payment", fiber.Map{
		"ErrorMessage": "",
	})
}

// вычитаем из баланса
func (card CardData) DecrementBalance() clause.Expr {
	return gorm.Expr("balance - ?", card.Balance)
}

// пополняем баланс
func (card CardData) IncrementBalance() clause.Expr {
	return gorm.Expr("balance + ?", card.Balance)
}

func (bank *BankRepos) DeductFromAccount(ctx *fiber.Ctx) error {
	details := ctx.Locals("cardDetails").(*CardData)

	fmt.Println("DETAILS BEFORE: ", details)

	if err := crypto.ProcessEncrypt((*crypto.CardData)(details)); err != nil {
		return fmt.Errorf("Encryption error: %v", err)
	}

	fmt.Println("DETAILS A: ", details)

	result := bank.DB.Model(newCard()).
		Where("encrypted_card_number = ?", details.EncryptedCardNumber).
		Where("encrypted_expiry_date = ?", details.EncryptedExpiryDate).
		Where("encrypted_CVV = ?", details.EncryptedCVV).
		Where("balance >= ?", details.Balance).
		Update("balance", details.DecrementBalance())

	if result.RowsAffected == 0 {
		return errors.New("Not enough money or card not found")
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Balance has been successfull deducted"})
}

func ShowTopUpPage(ctx *fiber.Ctx) error {
	return ctx.Render("add-funds", fiber.Map{
		"ErrorMessage": "",
	})
}

func (bank *BankRepos) AddFunds(ctx *fiber.Ctx) error {
	details := newCard()

	if err := ctx.BodyParser(details); err != nil {
		return ctx.Render("add-funds", fiber.Map{
			"ErrorMessage": "Failed to parse request",
		})
	}

	fmt.Printf("\nAdd funds: %v\n\n", details)

	ctx.Locals("cardDetails", details)

	if err := crypto.ProcessEncrypt((*crypto.CardData)(details)); err != nil {
		return fmt.Errorf("Encryption error: %v", err)
	}

	result := bank.DB.Model(newCard()).
		Where("encrypted_card_number = ?", details.EncryptedCardNumber).
		Update("balance", details.IncrementBalance())

	fmt.Printf("\nResult funds: %v\n\n", result)

	if result.RowsAffected == 0 {
		return ctx.Render("add-funds", fiber.Map{
			"ErrorMessage": "card not found"})
	}

	if err := SendTransacToKafka(ctx); err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "Funds added successfully"})
}

func (tr *Transactions) SetAmount(details CardData) {
	tr.Amount = details.Balance
}

type PaymentProcessor struct {
	Service *service.Service
}

func NewPaymentProcessor(service *service.Service) *PaymentProcessor {
	return &PaymentProcessor{Service: service}
}

func SendTransacToKafka(ctx *fiber.Ctx) error {
	details := ctx.Locals("cardDetails").(*CardData)

	fmt.Println("DETAILS ctx: ", details)

	// не забудь подставить сюда новую кафку
	// producer.ProducerTransac(details.Balance)

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{"message": "Transaction has been added"})
}

func (p *PaymentProcessor) ProcessPayment(ctx *fiber.Ctx) error {
	details := newCard()

	if err := ctx.BodyParser(details); err != nil {
		return ctx.Render("payment", fiber.Map{
			"ErrorMessage": "Failed to parse request",
		})
	}

	ctx.Locals("cardDetails", details)

	if err := p.Service.BankRepo.DeductFromAccount(ctx); err != nil {
		return ctx.Render("payment", fiber.Map{
			"ErrorMessage": err.Error(),
		})
	}

	if err := SendTransacToKafka(ctx); err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Payment processed successfully",
	})
}

func (bank *BankRepos) GetAllCardDetails(ctx *fiber.Ctx) error {
	balanceModels := &[]CardData{}
	if err := bank.DB.Find(&balanceModels).Error; err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get cards"})
		return err
	}

	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "cards fetched successfully",
		"data":    balanceModels,
	})
	return nil
}
