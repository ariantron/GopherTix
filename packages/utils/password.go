package utils

import (
	"golang.org/x/crypto/bcrypt"
	errs "gopher_tix/packages/common/errors"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errs.NewInternalServerError("Failed to password hashing")
	}
	return string(hashedPassword), nil
}
