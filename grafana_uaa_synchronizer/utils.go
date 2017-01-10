package main

import (
	"math/rand"
	"time"
)

const passwordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"


func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func RandomPassword(length int) string {
	password := make([]byte, length)
	for i := 0; i< length;i++ {
		randIndex := rand.Intn(len(passwordChars))
		password[i] = passwordChars[randIndex]
	}
	return string(password)
}

