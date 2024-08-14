package model

import (
	"fmt"
	"net/mail"
	"strings"

	"gorm.io/gorm"
)

const (
	ROLE_MANAGER = "manager"
	ROLE_ADMIN   = "admin"
)

type User struct {
	gorm.Model

	UserName string `json:"username"`
	Email    string `json:"email"`
	Age      uint   `json:"age"`
	Role     string `json:"role"`
}

func (user User) ValidateUser() error {
	if strings.TrimSpace(user.UserName) == "" {
		return fmt.Errorf("username can't be empty")
	}

	if strings.TrimSpace(user.Email) == "" {
		return fmt.Errorf("email address can't be empty")
	}

	if strings.TrimSpace(user.Role) == "" {
		return fmt.Errorf("role can't be empty")
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return fmt.Errorf("invalid email address")
	}

	if user.Role != ROLE_ADMIN && user.Role != ROLE_MANAGER {
		return fmt.Errorf("invalid user role")
	}

	return nil
}

func (user User) CleanData() User {
	// Cleanup data. Should also be done in frontend
	user.UserName = strings.ToLower(strings.TrimSpace(user.UserName))
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	return user
}
