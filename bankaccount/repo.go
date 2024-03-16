package bankaccount

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

func fetchBankAccounts(ctx context.Context, tx pgx.Tx, userId int) ([]BankAccount, error) {
	var bankAccountCount int
	queryCount := `SELECT COUNT(id) AS count FROM bank_accounts WHERE user_id = $1 AND deleted_at IS NULL`
	row := tx.QueryRow(ctx, queryCount, userId)
	err := row.Scan(&bankAccountCount)
	if err != nil {
		return []BankAccount{}, err
	}
	if bankAccountCount == 0 {
		return []BankAccount{}, nil
	}

	bankAccounts := make([]BankAccount, 0)

	query := `SELECT id, bank_name, bank_account_name, bank_account_number FROM bank_accounts WHERE user_id = $1 AND deleted_at IS NULL`
	rows, err := tx.Query(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bankAccount BankAccount

		err := rows.Scan(&bankAccount.ID, &bankAccount.BankName, &bankAccount.BankAccountName, &bankAccount.BankAccountNumber)
		if err != nil {
			return nil, err
		}

		bankAccounts = append(bankAccounts, bankAccount)
	}

	return bankAccounts, nil
}

func findBankAccountById(ctx context.Context, tx pgx.Tx, id int) (BankAccount, error) {
	query := `SELECT id, user_id, bank_name, bank_account_name, bank_account_number FROM bank_accounts WHERE id = $1 AND deleted_at IS NULL`
	row := tx.QueryRow(ctx, query, id)

	var bankAccount BankAccount
	if err := row.Scan(&bankAccount.ID, &bankAccount.UserID, &bankAccount.BankName, &bankAccount.BankAccountName, &bankAccount.BankAccountNumber); err != nil {
		if err == pgx.ErrNoRows {
			return BankAccount{}, ErrBankAccountNotFound
		}

		return BankAccount{}, err
	}

	return bankAccount, nil
}

func saveBankAccount(ctx context.Context, tx pgx.Tx, bankAccount BankAccount) (int, error) {
	lastInsertedId := 0
	query := `INSERT INTO bank_accounts(user_id, bank_name, bank_account_name, bank_account_number) 
	VALUES($1, $2, $3, $4) RETURNING id`

	err := tx.QueryRow(ctx, query, bankAccount.UserID, bankAccount.BankName, bankAccount.BankAccountName, bankAccount.BankAccountNumber).Scan(&lastInsertedId)
	if err != nil {
		return 0, err
	}

	return lastInsertedId, nil
}

func updateBankAccount(ctx context.Context, tx pgx.Tx, bankAccount BankAccount) error {
	query := `UPDATE bank_accounts SET bank_name=$1, bank_account_name=$2, bank_account_number=$3, 
	updated_at=$4 WHERE id = $5 AND deleted_at IS NULL`

	_, err := tx.Exec(ctx, query, bankAccount.BankName, bankAccount.BankAccountName,
		bankAccount.BankAccountNumber, time.Now(), bankAccount.ID)
	if err != nil {
		return err
	}

	return nil
}

func softDeleteBankAccount(ctx context.Context, tx pgx.Tx, bankAccountId int) error {
	query := `UPDATE bank_accounts SET deleted_at=$1 WHERE id = $2`

	_, err := tx.Exec(ctx, query, time.Now(), bankAccountId)
	if err != nil {
		return err
	}

	return nil
}
