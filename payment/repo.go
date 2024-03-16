package payment

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func findPaymentById(ctx context.Context, tx pgx.Tx, id int) (Payment, error) {
	query := `SELECT id, user_id, product_id, bank_account_id, payment_proof_image_url, quantity FROM payments WHERE id = $1 AND deleted_at IS NULL`
	row := tx.QueryRow(ctx, query, id)

	var payment Payment
	if err := row.Scan(&payment.ID, &payment.UserID, &payment.ProductID, &payment.BankAccountID, &payment.PaymentProofImageUrl, &payment.Quantity); err != nil {
		if err == pgx.ErrNoRows {
			return Payment{}, ErrPaymentNotFound
		}

		return Payment{}, err
	}

	return payment, nil

}

func savePayment(ctx context.Context, tx pgx.Tx, payment Payment) (int, error) {
	lastInsertedId := 0
	query := `INSERT INTO payments(user_id, product_id, bank_account_id, payment_proof_image_url, quantity) 
	VALUES($1, $2, $3, $4, $5) RETURNING id`

	err := tx.QueryRow(ctx, query, payment.UserID, payment.ProductID, payment.BankAccountID, payment.PaymentProofImageUrl, payment.Quantity).Scan(&lastInsertedId)
	if err != nil {
		return 0, err
	}

	return lastInsertedId, nil
}
