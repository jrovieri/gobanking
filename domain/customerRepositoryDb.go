package domain

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jrovieri/gobanking/errs"
	"github.com/jrovieri/gobanking/logger"
)

type CustomerRepositoryDb struct {
	client *sql.DB
}

func NewCustomerRepositoryDb() CustomerRepositoryDb {
	client, err := sql.Open("mysql", "gobanking:wks32x23@tcp(127.0.0.1:3306)/gobanking")
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

	var rows *sql.Rows
	var err error

	sqlStr := "select customer_id, name, city, zipcode, date_of_birth, status from customers"

	if status == "" {
		rows, err = d.client.Query(sqlStr)
	} else {
		sqlStr += " where status = ?"
		rows, err = d.client.Query(sqlStr, status)
	}

	if err != nil {
		logger.Error("Error while querying the customer table:" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	customers := make([]Customer, 0)

	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.Id, &c.Name, &c.City, &c.Zipcode, &c.Birthdate, &c.Status)
		if err != nil {
			logger.Error("Error while scanning customers: " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func (d CustomerRepositoryDb) ById(id string) (*Customer, *errs.AppError) {
	sqlStr := "select customer_id, name, city, zipcode, date_of_birth, status from customers where customer_id = ?"

	row := d.client.QueryRow(sqlStr, id)
	var c Customer

	err := row.Scan(&c.Id, &c.Name, &c.City, &c.Zipcode, &c.Birthdate, &c.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Customer not found")
		} else {
			logger.Error("Error while scanning customers: " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}
	return &c, nil
}
