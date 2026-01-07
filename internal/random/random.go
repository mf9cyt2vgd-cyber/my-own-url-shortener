package random

import "math/rand/v2"

func NewRandomString(length int) string {
	data := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = data[rand.IntN(len(data))]
	}

	return string(result)
}
