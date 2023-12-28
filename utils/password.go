package utils

import "golang.org/x/crypto/bcrypt"

func GenerateHashedPassword(password string) (string,error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func VerifyPassword(password, hashPassword string) (error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err
}