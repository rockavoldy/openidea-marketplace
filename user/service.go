package user

import (
	"context"
	"net/http"
)

func ValidateUser(user User) error {
	if err := user.validateUsername(); err != nil {
		return err
	}
	if err := user.validateName(); err != nil {
		return err
	}
	if err := user.validatePassword(); err != nil {
		return err
	}

	return nil
}

func ValidateUserLogin(user User) error {
	if err := user.validateUsername(); err != nil {
		return err
	}
	if err := user.validatePassword(); err != nil {
		return err
	}

	return nil
}

func GetUserByUsername(ctx context.Context, username string) (User, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return User{}, err
	}
	user, err := findUserByUsername(ctx, tx, username)
	if err != nil {
		tx.Rollback(ctx)
		return User{}, err
	}

	tx.Commit(ctx)
	return user, nil
}

func registerUser(ctx context.Context, username, name, password string) (user User, statusCode int, err error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return User{}, http.StatusInternalServerError, err
	}

	user, err = GetUserByUsername(ctx, username)
	if err == nil && user.exist() {
		// when username already exist
		tx.Rollback(ctx)
		return user, http.StatusConflict, ErrUsernameExist
	}

	if err != nil {
		if err != ErrUserNotFound {
			tx.Rollback(ctx)
			return User{}, http.StatusInternalServerError, err
		}

		// when user not found, means it can be allowed to create new user
		user = NewUser(username, name, password)
		user.hashPassword()
		err := saveUser(ctx, tx, user)
		if err != nil {
			tx.Rollback(ctx)
			return User{}, http.StatusInternalServerError, err
		}
	}

	tx.Commit(ctx)
	return user, http.StatusCreated, nil
}

func loginUser(ctx context.Context, username, password string) (user User, statusCode int, err error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return User{}, http.StatusInternalServerError, err
	}
	user, err = GetUserByUsername(ctx, username)
	if err != nil {
		tx.Rollback(ctx)
		if err == ErrUserNotFound {
			return User{}, http.StatusNotFound, ErrUserNotFound
		}

		return User{}, http.StatusInternalServerError, ErrServerError
	}

	if err := user.comparePassword(password); err != nil {
		tx.Rollback(ctx)
		return User{}, http.StatusBadRequest, ErrWrongPassword
	}

	tx.Commit(ctx)
	return user, http.StatusOK, nil
}
