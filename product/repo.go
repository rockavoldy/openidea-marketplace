package product

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func getProductTagsByProductId(ctx context.Context, tx pgx.Tx, productId int) ([]string, error) {
	var tagsCount int
	queryCount := `SELECT COUNT(id) AS count FROM product_tags WHERE product_id = $1`
	row := tx.QueryRow(ctx, queryCount, productId)
	err := row.Scan(&tagsCount)
	if err != nil {
		return nil, err
	}
	if tagsCount == 0 {
		return nil, nil
	}

	tags := make([]string, 0)

	query := `SELECT name FROM product_tags WHERE product_id = $1`
	rows, err := tx.Query(ctx, query, productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string

		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}

		tags = append(tags, name)
	}

	return tags, nil
}

func saveProductTags(ctx context.Context, tx pgx.Tx, productId int, names []string) ([]string, error) {
	query := `DELETE FROM product_tags WHERE product_id = $1`
	_, err := tx.Exec(ctx, query, productId)
	if err != nil {
		return nil, err
	}

	log.Println(names)

	rows := make([][]any, 0)

	for _, item := range names {
		rows = append(rows, []any{productId, item})
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"product_tags"},
		[]string{"product_id", "name"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return nil, err
	}

	return getProductTagsByProductId(ctx, tx, productId)
}

func findProductById(ctx context.Context, tx pgx.Tx, id int) (Product, error) {
	query := `SELECT id, user_id, name, price, image_url, stock, condition, is_purchasable, created_at, updated_at, deleted_at FROM products WHERE id = $1 AND deleted_at IS NULL`
	row := tx.QueryRow(ctx, query, id)

	var product Product
	if err := row.Scan(&product.ID, &product.UserID, &product.Name, &product.Price, &product.ImageURL, &product.Stock, &product.Condition, &product.IsPurchasable, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt); err != nil {
		if err == pgx.ErrNoRows {
			return Product{}, ErrProductNotFound
		}

		return Product{}, err
	}

	// put tags
	tags, err := getProductTagsByProductId(ctx, tx, product.ID)
	if err != nil {
		return Product{}, err
	}
	product.Tags = tags

	return product, nil
}

func saveProduct(ctx context.Context, tx pgx.Tx, product Product) (int, error) {
	lastInsertedId := 0
	query := `INSERT INTO products(user_id, name, price, image_url, stock, condition, is_purchasable) 
	VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := tx.QueryRow(ctx, query, product.UserID, product.Name, product.Price,
		product.ImageURL, product.Stock, product.Condition, product.IsPurchasable).Scan(&lastInsertedId)
	if err != nil {
		return 0, err
	}

	return lastInsertedId, nil
}

func updateProduct(ctx context.Context, tx pgx.Tx, product Product) error {
	query := `UPDATE products SET name=$1, price=$2, image_url=$3, 
 stock=$4, condition=$5, is_purchasable=$6, updated_at=$7 WHERE id = $8 AND deleted_at IS NULL`

	_, err := tx.Exec(ctx, query, product.Name, product.Price, product.ImageURL,
		product.Stock, product.Condition, product.IsPurchasable, time.Now(), product.ID)
	if err != nil {
		return err
	}

	_, err = saveProductTags(ctx, tx, product.ID, product.Tags)
	if err != nil {
		return err
	}

	return nil
}

func softDeleteProduct(ctx context.Context, tx pgx.Tx, productId int) error {
	query := `UPDATE products SET deleted_at=$1 WHERE id = $2`

	_, err := tx.Exec(ctx, query, time.Now(), productId)
	if err != nil {
		return err
	}

	return nil
}
