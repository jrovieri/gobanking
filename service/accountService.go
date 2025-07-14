package service

import (
	"time"

	"github.com/jrovieri/gobanking/domain"
	"github.com/jrovieri/gobanking/dto"
	"github.com/jrovieri/gobanking/errs"
)

type AccountService interface {
	NewAccount(dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError)
}

type DefaultAccountService struct {
	repository domain.AccountRepositoryDb
}

func (s DefaultAccountService) NewAccount(req dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	acc := domain.Account{
		AccountId:   "",
		CustomerId:  req.CustomerId,
		OpeningDate: time.Now().Format("2006-01-02T15:04:05"),
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

func NewAccountService(r domain.AccountRepositoryDb) DefaultAccountService {
	return DefaultAccountService{repository: r}
}
