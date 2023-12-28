package utils

import "math/rand"

var alphabets string = "abcdefghijklmnopqrstuvwxyz";

func RandomString(length int) string {
  bits := []rune{}
	k := len(alphabets)

	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(k)
		bits = append(bits, rune(alphabets[randomIndex]))
	}

	return string(bits)
}

func RandomEmail() string {
	return RandomString(8) + "@" + RandomString(4) + ".com"
}