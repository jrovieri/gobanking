package service

import (
	"time"

	"github.com/jrovieri/gobanking/domain"
	"github.com/jrovieri/gobanking/dto"
	"github.com/jrovieri/gobanking/errs"
)

const DB_TIME_LAYOUT = "2006-01-02T15:04:05"

type AccountService interface {
	NewAccount(dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError)
	MakeTransaction(dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppError)
}

type DefaultAccountService struct {
	repository domain.AccountRepositoryDb
}

func NewAccountService(r domain.AccountRepositoryDb) DefaultAccountService {
	return DefaultAccountService{repository: r}
}

func (s DefaultAccountService) NewAccount(req dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	acc := domain.Account{
		AccountId:   "",
		CustomerId:  req.CustomerId,
		OpeningDate: time.Now().Format(DB_TIME_LAYOUT),
		AccountType: req.AccountType,
		Amount:      req.Amount,
		Status:      "1",
	}

	newAccount, err := s.repository.Save(acc)
	if err != nil {
		return nil, err
	}

	response := newAccount.ToNewAccountResponseDto()
	return &response, nil
}

func (s DefaultAccountService) MakeTransaction(req dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppError) {

	err := req.Validate()
	if err != nil {
		return nil, err
	}

	if req.IsTransactonTypeWithdrawal() {
		account, err := s.repository.FindById(req.AccountId)
		if err != nil {
			return nil, err
		}
		if !account.CanWithdraw(req.Amount) {
			return nil, errs.NewValidationError("Insufficient balance in the account")
		}
	}

	t := domain.Transaction{
		AccountId:       req.AccountId,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format(DB_TIME_LAYOUT),
	}

	transaction, err := s.repository.SaveTransaction(t)
	if err != nil {
		return nil, err
	}
	response := transaction.ToDto()
	return &response, nil
}
