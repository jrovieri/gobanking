package domain

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jrovieri/gobanking/errs"
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

func (d CustomerRepositoryDb) FindAll() ([]Customer, *errs.AppError) {

	sql := "select customer_id, name, city, zipcode, date_of_birth, status from customers"
	rows, err := d.client.Query(sql)
	if err != nil {
		log.Println("Error while querying the customer table:\n" + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	customers := make([]Customer, 0)

	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.Id, &c.Name, &c.City, &c.Zipcode, &c.Birthdate, &c.Status)
		if err != nil {
			fmt.Println("Error while scanning customers: " + err.Error())
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
			log.Println("Error while scanning customer " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}
	return &c, nil
}
