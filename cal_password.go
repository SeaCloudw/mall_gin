package main

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
func CheckPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err

	}
	return nil
}

// func main() {
// 	r1, _ := HashPassword("1")
// 	fmt.Println(r1)
// 	fmt.Println(CheckPassword("1", r1) == nil)

// }
