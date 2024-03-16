package payment

import (
	productObj "MarketplaceAPI/product"
	"context"
	"net/http"
)

func getPayment(ctx context.Context, id int) (Payment, int, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return Payment{}, http.StatusInternalServerError, err
	}
	defer tx.Commit(ctx)

	payment, err := findPaymentById(ctx, tx, id)
	if err != nil {
		if err == ErrPaymentNotFound {
			return Payment{}, http.StatusNotFound, err
		}
		return Payment{}, http.StatusInternalServerError, err
	}

	return payment, http.StatusOK, nil
}

func buyProduct(ctx context.Context, payment Payment) (Payment, int, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return Payment{}, http.StatusInternalServerError, err
	}
	defer tx.Commit(ctx)

	userId, ok := ctx.Value("userId").(int)
	if ok {
		payment.addUserID(userId)
	}

	product, _, err := productObj.GetProduct(ctx, payment.ProductID)
	if err != nil {
		tx.Rollback(ctx)
		return Payment{}, http.StatusInternalServerError, err
	}

	product.AdjustStock(product.Stock - payment.Quantity)

	product, _, err = productObj.PatchProduct(ctx, product)
	if err != nil {
		tx.Rollback(ctx)
		return Payment{}, http.StatusInternalServerError, err
	}

	id, err := savePayment(ctx, tx, payment)
	if err != nil {
		tx.Rollback(ctx)
		return Payment{}, http.StatusInternalServerError, err
	}
	tx.Commit(ctx)

	return getPayment(ctx, id)
}
