package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"
)

func GetSHA256Hash(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetKeyToken() string {
	keyToken := os.Getenv("KEY_TOKEN")

	return keyToken
}

func TimeNowUTC() time.Time {
	utc, _ := time.LoadLocation("UTC")
	return time.Now().In(utc)
}

func FindIntInSlice(slice []int, val int) bool {
	if len(slice) == 0 {
		return false
	}

	for _, item := range slice {
		if item == val {
			return true
		}
	}

	return false
}
