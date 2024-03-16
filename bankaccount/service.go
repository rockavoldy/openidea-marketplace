package bankaccount

import (
	"context"
	"log"
	"net/http"
)

func getBankAccount(ctx context.Context, id int) (BankAccount, int, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return BankAccount{}, http.StatusInternalServerError, err
	}
	defer tx.Commit(ctx)

	bankAccount, err := findBankAccountById(ctx, tx, id)
	if err != nil {
		if err == ErrBankAccountNotFound {
			return BankAccount{}, http.StatusNotFound, err
		}
		return BankAccount{}, http.StatusInternalServerError, err
	}
	log.Println(bankAccount)

	return bankAccount, http.StatusOK, nil
}

func listBankAccounts(ctx context.Context) ([]BankAccount, int, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return []BankAccount{}, http.StatusInternalServerError, err
	}
	defer tx.Commit(ctx)

	userId, ok := ctx.Value("userId").(int)
	if !ok {
		return []BankAccount{}, http.StatusForbidden, ErrForbidden
	}

	bankAccounts, err := fetchBankAccounts(ctx, tx, userId)
	if err != nil {
		return []BankAccount{}, http.StatusInternalServerError, err
	}

	return bankAccounts, http.StatusOK, nil
}

func createBankAccount(ctx context.Context, bankAccount BankAccount) (BankAccount, int, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return BankAccount{}, http.StatusInternalServerError, err
	}

	userId, ok := ctx.Value("userId").(int)
	if ok {
		bankAccount.UserID = userId
	}

	id, err := saveBankAccount(ctx, tx, bankAccount)
	log.Println(id)
	if err != nil {
		tx.Rollback(ctx)
		return BankAccount{}, http.StatusInternalServerError, err
	}
	tx.Commit(ctx)

	return getBankAccount(ctx, id)
}

func patchBankAccount(ctx context.Context, bankAccount BankAccount) (BankAccount, int, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return BankAccount{}, http.StatusInternalServerError, err
	}
	defer tx.Commit(ctx)

	if err := updateBankAccount(ctx, tx, bankAccount); err != nil {
		tx.Rollback(ctx)
		return BankAccount{}, http.StatusInternalServerError, err
	}

	return bankAccount, http.StatusOK, nil
}

func deleteBankAccount(ctx context.Context, bankAccountId int) (BankAccount, int, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return BankAccount{}, http.StatusInternalServerError, err
	}
	defer tx.Commit(ctx)

	bankAccount, statusCode, err := getBankAccount(ctx, bankAccountId)
	if err != nil {
		tx.Rollback(ctx)
		return BankAccount{}, statusCode, err
	}

	// add deletedAt to softdelete
	err = softDeleteBankAccount(ctx, tx, bankAccountId)
	if err != nil {
		tx.Rollback(ctx)
		return BankAccount{}, statusCode, err
	}

	return bankAccount, statusCode, nil

}
