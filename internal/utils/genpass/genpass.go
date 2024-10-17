package genpass

import "math/rand"

func GeneratePassword() string {
	var newPass string

	symbols := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890@!#$%_?"

	for i := 0; i < 6; i += 1 {
		num := rand.Intn(len(symbols))
		newPass += string(symbols[num])
	}

	return newPass
}
