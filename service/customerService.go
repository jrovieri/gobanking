package service

import (
	"github.com/jrovieri/gobanking/domain"
	"github.com/jrovieri/gobanking/errs"
)

type CustomerService interface {
	GetAllCustomers(string) ([]domain.Customer, *errs.AppError)
	GetCustomer(string) (*domain.Customer, *errs.AppError)
}

type DefaultCustomerService struct {
	repository domain.CustomerRepository
}

func NewCustomerService(repository domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repository: repository}
}

func (s DefaultCustomerService) GetAllCustomers(status string) ([]domain.Customer, *errs.AppError) {

	switch status {
	case "active":
		status = "1"
	case "inactive":
		status = "0"
	default:
		status = ""
	}
	return s.repository.FindAll(status)
}

func (s DefaultCustomerService) GetCustomer(id string) (*domain.Customer, *errs.AppError) {
	return s.repository.ById(id)
}
