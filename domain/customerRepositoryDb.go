package domain

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jrovieri/gobanking/errs"
	"github.com/jrovieri/gobanking/logger"
)

type CustomerRepositoryDb struct {
	client *sqlx.DB
}

func NewCustomerRepositoryDb() CustomerRepositoryDb {
	client, err := sqlx.Open("mysql", "gobanking:wks32x23@tcp(127.0.0.1:3306)/gobanking")
	if err != nil {
		panic(err)
	}

	if err := client.Ping(); err != nil {
		log.Fatalln(err)
	}

	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return CustomerRepositoryDb{client}
}

func (d CustomerRepositoryDb) FindAll(status string) ([]Customer, *errs.AppError) {

	var err error
	customers := make([]Customer, 0)

	sqlStr := "select customer_id, name, city, zipcode, date_of_birth, status from customers"

	if status == "" {
		err = d.client.Select(&customers, sqlStr)
	} else {
		sqlStr += " where status = ?"
		err = d.client.Select(&customers, sqlStr, status)
	}

	if err != nil {
		logger.Error("Error while querying the customer table:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	return customers, nil
}

func (d CustomerRepositoryDb) ById(id string) (*Customer, *errs.AppError) {

	var customer Customer
	sqlStr := "select customer_id, name, city, zipcode, date_of_birth, status from customers where customer_id = ?"

	err := d.client.Get(&customer, sqlStr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Customer not found")
		} else {
			logger.Error("Error while scanning customers: " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}
	return &customer, nil
}
