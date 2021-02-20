package util

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func CompareBcrypt(wanted, actual, pepper string) bool {
	plainPassword := []byte(strings.TrimSpace(actual) + pepper)
	err := bcrypt.CompareHashAndPassword([]byte(wanted), plainPassword)
	return err == nil
}

func HashBcrypt(plain, pepper string) string {
	if plain == "" {
		return ""
	}
	plainPepperedPassword := []byte(strings.TrimSpace(plain) + pepper)
	bytes, _ := bcrypt.GenerateFromPassword(plainPepperedPassword, bcrypt.DefaultCost)
	return string(bytes)
}
