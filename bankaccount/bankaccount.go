package bankaccount

import (
	"fmt"

	"github.com/jackc/pgtype"
)

type BankAccount struct {
	ID                int              `json:"bankAccountId"`
	UserID            int              `json:"-"`
	BankName          string           `json:"bankName"`
	BankAccountName   string           `json:"bankAccountName"`
	BankAccountNumber string           `json:"bankAccountNumber"`
	CreatedAt         pgtype.Timestamp `json:"-"`
	UpdatedAt         pgtype.Timestamp `json:"-"`
	DeletedAt         pgtype.Timestamp `json:"-"`
}

var (
	ErrBankAccountNotFound = fmt.Errorf("product not found")
	ErrForbidden           = fmt.Errorf("access forbidden")
)

func NewBankAccount(bankName, bankAccountName, bankAccountNumber string) BankAccount {
	return BankAccount{
		BankName:          bankName,
		BankAccountName:   bankAccountName,
		BankAccountNumber: bankAccountNumber,
	}
}

func (ba *BankAccount) patchWith(patch map[string]string) {
	for k := range patch {
		if k == "bankName" {
			ba.BankName = patch[k]
		} else if k == "bankAccountName" {
			ba.BankAccountName = patch[k]
		} else if k == "bankAccountNumber" {
			ba.BankAccountNumber = patch[k]
		}
	}
}
