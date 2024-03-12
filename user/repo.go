package user

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func findUserById(ctx context.Context, tx pgx.Tx, id int) (User, error) {
	query := `SELECT id, username, name FROM users WHERE id = $1;`
	row := tx.QueryRow(ctx, query, id)

	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Name); err != nil {
		if err == pgx.ErrNoRows {
			return User{}, ErrUserNotFound
		}

		return User{}, err
	}

	return user, nil
}

func findUserByUsername(ctx context.Context, tx pgx.Tx, username string) (User, error) {
	query := `SELECT id, username, name, password FROM users WHERE username = $1;`
	row := tx.QueryRow(ctx, query, username)

	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Name, &user.Password); err != nil {
		if err == pgx.ErrNoRows {
			return User{}, ErrUserNotFound
		}

		return User{}, err
	}

	return user, nil
}

func saveUser(ctx context.Context, tx pgx.Tx, user User) error {
	query := `INSERT INTO users(username, name, password) VALUES($1, $2, $3);`

	_, err := tx.Exec(ctx, query, user.Username, user.Name, user.Password)
	if err != nil {
		return err
	}
	tx.Commit(ctx)

	return nil
}
