package product

import (
	"context"
	"net/http"
)

func listProducts(ctx context.Context) ([]Product, int, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return []Product{}, http.StatusInternalServerError, err
	}
	defer tx.Commit(ctx)

	products, err := fetchProducts(ctx, tx)
	if err != nil {
		return []Product{}, http.StatusInternalServerError, err
	}

	return products, http.StatusOK, nil
}

func GetProduct(ctx context.Context, id int) (Product, int, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return Product{}, http.StatusInternalServerError, err
	}
	defer tx.Commit(ctx)

	product, err := findProductById(ctx, tx, id)
	if err != nil {
		if err == ErrProductNotFound {
			return Product{}, http.StatusNotFound, err
		}
		return Product{}, http.StatusInternalServerError, err
	}

	return product, http.StatusOK, nil
}

func createProduct(ctx context.Context, product Product) (Product, int, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return Product{}, http.StatusInternalServerError, err
	}

	userId, ok := ctx.Value("userId").(int)
	if ok {
		product.addUserID(userId)
	}

	id, err := saveProduct(ctx, tx, product)
	if err != nil {
		tx.Rollback(ctx)
		return Product{}, http.StatusInternalServerError, err
	}

	_, err = saveProductTags(ctx, tx, id, product.Tags)
	if err != nil {
		return product, http.StatusInternalServerError, err
	}
	tx.Commit(ctx)

	return GetProduct(ctx, id)
}

func PatchProduct(ctx context.Context, product Product) (Product, int, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return Product{}, http.StatusInternalServerError, err
	}
	defer tx.Commit(ctx)

	if err := updateProduct(ctx, tx, product); err != nil {
		tx.Rollback(ctx)
		return Product{}, http.StatusInternalServerError, err
	}

	return product, http.StatusOK, nil
}

func deleteProduct(ctx context.Context, productId int) (Product, int, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return Product{}, http.StatusInternalServerError, err
	}
	defer tx.Commit(ctx)

	product, statusCode, err := GetProduct(ctx, productId)
	if err != nil {
		tx.Rollback(ctx)
		return Product{}, statusCode, err
	}

	// add deletedAt to softdelete
	err = softDeleteProduct(ctx, tx, productId)
	if err != nil {
		tx.Rollback(ctx)
		return Product{}, statusCode, err
	}

	return product, statusCode, nil

}

func StockAdjustment(ctx context.Context, qty int, product Product) (Product, int, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return Product{}, http.StatusInternalServerError, err
	}
	defer tx.Commit(ctx)

	if err := updateProduct(ctx, tx, product); err != nil {
		tx.Rollback(ctx)
		return Product{}, http.StatusInternalServerError, err
	}

	return product, http.StatusOK, nil
}
