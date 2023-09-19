package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func main() {
	// Generate a random string of length 32
	secret, err := generateRandomString(32)
	if err != nil {
		panic(err)
	}
	fmt.Println("Generated JWT Secret:", secret)
}
