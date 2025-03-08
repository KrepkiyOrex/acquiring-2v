package service

import (
	"github.com/KrepkiyOrex/acquiring/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type Service struct {
	BankRepo        domain.BankRepository
	TransactionRepo domain.TransactionRepository
}

func NewService(bankRepo domain.BankRepository, transactionRepo domain.TransactionRepository) *Service {
	return &Service{
		BankRepo:        bankRepo,
		TransactionRepo: transactionRepo,
	}
}

// ================== Bank Repository ==================

func (s *Service) CreateUserCard(ctx *fiber.Ctx) error {
	return s.BankRepo.CreateUserCard(ctx)
}

func (s *Service) AddFunds(ctx *fiber.Ctx) error {
	return s.BankRepo.AddFunds(ctx)
}

func (s *Service) DeductFromAccount(ctx *fiber.Ctx) error {
	return s.BankRepo.DeductFromAccount(ctx)
}

func (s *Service) GetAllCardDetails(ctx *fiber.Ctx) error {
	return s.BankRepo.GetAllCardDetails(ctx)
}

// ================== Transaction Repository ==================

func (s *Service) DeleteTransaction(ctx *fiber.Ctx, txID int64) error {
	return s.TransactionRepo.DeleteTransaction(ctx, txID)
}

func (s *Service) GetTransactions(ctx *fiber.Ctx) error {
	return s.TransactionRepo.GetTransactions(ctx)
}

func (s *Service) GetTransByID(ctx *fiber.Ctx, txID int64) error {
	return s.TransactionRepo.GetTransByID(ctx, txID)
}
