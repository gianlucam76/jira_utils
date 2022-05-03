package jira

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/go-logr/logr"
)

func GetUsername(logger logr.Logger) string {
	user, ok := os.LookupEnv(username)
	if !ok {
		logger.Info(fmt.Sprintf("Env variable %s supposed to contain username not found", username))
		panic(1)
	}

	if user == "" {
		logger.Info("Username cannot be emty")
		panic(1)
	}

	return user
}

func GetPassword(logger logr.Logger) string {
	base64Password, ok := os.LookupEnv(password)
	if !ok {
		logger.Info(fmt.Sprintf("Env variable %s supposed to contain password, base64 encoed, not found", password))
		panic(1)
	}

	if base64Password == "" {
		logger.Info("Password cannot be emty")
		panic(1)
	}

	password, err := base64.StdEncoding.DecodeString(base64Password)
	if err != nil {
		logger.Info(fmt.Sprintf("Failed to decode password: %v", err))
		panic(err)
	}
	return string(password)
}
