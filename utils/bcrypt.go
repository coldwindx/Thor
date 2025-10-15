package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func BcryptMake(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if nil != err {
		panic(err)
	}
	return string(hash)
}

func BcryptMakeCheck(pwd []byte, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), pwd)
	return nil == err
}
