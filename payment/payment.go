package payment

import (
	"fmt"

	"github.com/jackc/pgtype"
)

type Payment struct {
	ID                   int              `json:"id"`
	UserID               int              `json:"userId"`
	ProductID            int              `json:"productId"`
	BankAccountID        int              `json:"bankAccountId"`
	PaymentProofImageUrl string           `json:"paymentProofImageUrl"`
	Quantity             uint             `json:"quantity"`
	CreatedAt            pgtype.Timestamp `json:"createdAt"`
	UpdatedAt            pgtype.Timestamp `json:"updatedAt"`
	DeletedAt            pgtype.Timestamp `json:"deletedAt"`
}

var (
	ErrPaymentNotFound        = fmt.Errorf("payment not found")
	ErrInsufficientQty        = fmt.Errorf("insufficient quantity")
	ErrIncorrectPaymentDetail = fmt.Errorf("incorrect payment detail")
)

func NewPayment(userId, productId, bankAccountId int) Payment {
	return Payment{
		UserID:        userId,
		ProductID:     productId,
		BankAccountID: bankAccountId,
	}
}
