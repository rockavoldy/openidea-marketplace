package user

import (
	"MarketplaceAPI/utils"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

var (
	ErrServerError     = fmt.Errorf("server error")
	ErrUserNotFound    = fmt.Errorf("user not found")
	ErrUsernameExist   = fmt.Errorf("username already exist")
	ErrUsernameMissing = fmt.Errorf("username is missing")
	ErrNameMissing     = fmt.Errorf("name is missing")
	ErrPasswordMissing = fmt.Errorf("password is missing")
	ErrUsernameLength  = fmt.Errorf("username is too short or too long")
	ErrNameLength      = fmt.Errorf("name is too short or too long")
	ErrPasswordLength  = fmt.Errorf("password is too short or too long")
	ErrWrongPassword   = fmt.Errorf("wrong password")
)

func NewUser(username, name, password string) User {
	return User{
		Username: username,
		Name:     name,
		Password: password,
	}
}

func (u *User) exist() bool {
	if u.Username != "" {
		return true
	}

	return false
}

func (u *User) hashPassword() error {
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword

	return nil
}

func (u *User) comparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) generateToken() string {
	token, err := utils.CreateToken(u.ID, u.Username, u.Name)
	if err != nil {
		log.Println(err)
	}

	return token
}

func (u *User) validateUsername() error {
	// Username is not null, and length between 5 to 15
	if u.Username == "" {
		return ErrUsernameMissing
	}

	if len(u.Username) < 5 || len(u.Username) > 15 {
		return ErrUsernameLength
	}

	return nil
}

func (u *User) validateName() error {
	// name is not null, and length between 5 to 50
	if u.Name == "" {
		return ErrNameMissing
	}

	if len(u.Name) < 5 || len(u.Name) > 50 {
		return ErrNameLength
	}

	return nil
}

func (u *User) validatePassword() error {
	// password is not null, and length between 5 to 15
	if u.Password == "" {
		return ErrPasswordMissing
	}

	if len(u.Password) < 5 || len(u.Password) > 15 {
		return ErrPasswordLength
	}

	return nil
}
