package domain

import "github.com/jrovieri/gobanking/errs"

type Customer struct {
	Id        string
	Name      string
	City      string
	Zipcode   string
	Birthdate string
	Status    string
}

type CustomerRepository interface {
	FindAll(status string) ([]Customer, *errs.AppError)
	ById(string) (*Customer, *errs.AppError)
}
