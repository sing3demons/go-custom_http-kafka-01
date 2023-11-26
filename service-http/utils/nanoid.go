package utils

import (
	"github.com/aidarkhanov/nanoid/v2"
)

const (
	alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func RandomNanoID(size uint) (string, error) {
	id, err := nanoid.GenerateString(alphabet, int(size))
	if err != nil {
		return "", err
	}
	return id, nil
}
