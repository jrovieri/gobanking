package domain

import (
	"database/sql"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/jrovieri/gobanking/errs"
	"github.com/jrovieri/gobanking/logger"
)

type AccountRepositoryDb struct {
	client *sqlx.DB
}

func NewAccountRepositoryDb(dbClient *sqlx.DB) AccountRepositoryDb {
	return AccountRepositoryDb{client: dbClient}
}

func (d AccountRepositoryDb) Save(a Account) (*Account, *errs.AppError) {
	sqlStr := "INSERT INTO accounts (customer_id, opening_date, account_type, amount, status) VALUES (?, ?, ?, ?, ?)"

	res, err := d.client.Exec(sqlStr, a.CustomerId, a.OpeningDate, a.AccountType, a.Amount, a.Status)
	if err != nil {
		logger.Error("Error while creating new account: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	id, err := res.LastInsertId()
	if err != nil {
		logger.Error("Error while getting last insert id for new account: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	a.AccountId = strconv.FormatInt(id, 10)
	return &a, nil
}

func (d AccountRepositoryDb) SaveTransaction(t Transaction) (*Transaction, *errs.AppError) {
	tx, err := d.client.Begin()
	if err != nil {
		logger.Error("Error while starting a new transaction for bank account: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	var updateErr error
	res, _ := tx.Exec(`INSERT INTO transactions (account_id, amount, transaction_type, transaction_date) 
		VALUES (?, ?, ?, ?)`, t.AccountId, t.Amount, t.TransactionType, t.TransactionDate)

	if t.IsWithdraw() {
		_, updateErr = tx.Exec("UPDATE accounts SET amount = amount - ? WHERE account_id = ?", t.Amount, t.AccountId)
	} else {
		_, updateErr = tx.Exec("UPDATE accounts SET amount = amount + ? WHERE account_id = ?", t.Amount, t.AccountId)
	}

	if updateErr != nil {
		tx.Rollback()
		logger.Error("Error while saving transaction: " + updateErr.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Error("Error while commiting transaction for bank account: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	transactionId, err := res.LastInsertId()
	if err != nil {
		logger.Error("Error while getting the last transaction id: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	account, appErr := d.FindById(t.AccountId)
	if appErr != nil {
		return nil, appErr
	}
	t.TransactionId = strconv.FormatInt(transactionId, 10)

	t.Amount = account.Amount
	return &t, nil
}

func (d AccountRepositoryDb) FindById(accountId string) (*Account, *errs.AppError) {

	var account Account
	sqlStr := "SELECT account_id, customer_id, opening_date, account_type, amount, status FROM accounts WHERE account_id = ?"

	err := d.client.Get(&account, sqlStr, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Account not found")
		} else {
			logger.Error("Error while scanning accounts: " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}
	return &account, nil
}
