package service

import (
	"github.com/jrovieri/gobanking/domain"
	"github.com/jrovieri/gobanking/dto"
	"github.com/jrovieri/gobanking/errs"
)

type CustomerService interface {
	GetAllCustomers(string) ([]dto.CustomerResponse, *errs.AppError)
	GetCustomer(string) (*dto.CustomerResponse, *errs.AppError)
}

type DefaultCustomerService struct {
	repository domain.CustomerRepository
}

func NewCustomerService(repository domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repository: repository}
}

func (s DefaultCustomerService) GetAllCustomers(status string) ([]dto.CustomerResponse, *errs.AppError) {

	switch status {
	case "active":
		status = "1"
	case "inactive":
		status = "0"
	default:
		status = ""
	}

	customers, err := s.repository.FindAll(status)
	if err != nil {
		return nil, err
	}

	customersResponse := make([]dto.CustomerResponse, 0)
	for _, c := range customers {
		customersResponse = append(customersResponse, c.ToDto())
	}
	return customersResponse, err
}

func (s DefaultCustomerService) GetCustomer(id string) (*dto.CustomerResponse, *errs.AppError) {

	customer, err := s.repository.ById(id)
	if err != nil {
		return nil, err
	}
	response := customer.ToDto()
	return &response, nil
}
