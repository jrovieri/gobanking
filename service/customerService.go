package service

import (
	"github.com/jrovieri/gobanking/domain"
	"github.com/jrovieri/gobanking/errs"
)

type CustomerService interface {
	GetAllCustomers() ([]domain.Customer, *errs.AppError)
	GetCustomer(string) (*domain.Customer, *errs.AppError)
}

type DefaultCustomerService struct {
	repository domain.CustomerRepository
}

func NewCustomerService(repository domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repository: repository}
}

func (s DefaultCustomerService) GetAllCustomers() ([]domain.Customer, *errs.AppError) {
	return s.repository.FindAll()
}

func (s DefaultCustomerService) GetCustomer(id string) (*domain.Customer, *errs.AppError) {
	return s.repository.ById(id)
}
