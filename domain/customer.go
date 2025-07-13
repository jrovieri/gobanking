package domain

import "github.com/jrovieri/gobanking/errs"

type Customer struct {
	Id        string `db:"customer_id"`
	Name      string
	City      string
	Zipcode   string
	Birthdate string `db:"date_of_birth"`
	Status    string
}

type CustomerRepository interface {
	FindAll(status string) ([]Customer, *errs.AppError)
	ById(string) (*Customer, *errs.AppError)
}
