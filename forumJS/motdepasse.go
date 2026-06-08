package forumjs

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	// return password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return password
	}

	return string(hashedPassword)
}

func CheckPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == nil {
		return true
	}
	return false
}

/*
azerty

c3grty

*/
