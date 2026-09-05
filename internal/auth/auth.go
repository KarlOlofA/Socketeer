package auth

import (
	"errors"
	"os"
)

type User struct {
	IpAddress string `json:"ipAdress"`
	IsAuthed  bool   `json:"isAuthed"`
}

func (u *User) ValidateApiKey(apiKey string) (string, error) {
	ak := os.Getenv("APIKEY")
	if ak != apiKey {
		return "", errors.New("Failed to validate api key.")
	}

	return apiKey, nil
}
