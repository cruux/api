package config

import (
	"errors"
	"fmt"
	"os"
)

var Env = map[string]string{
	"PORT":            "",
	"MONGO_HOST":      "",
	"JWT_PUBLIC_KEY":  "",
	"JWT_PRIVATE_KEY": "",
}

func initEnv() {
	for k, _ := range Env {
		Env[k] = os.Getenv(k)
		if Env[k] == "" {
			fmt.Println("Error: Env var " + k + " must be set.")
			panic(errors.New("Env var " + k + " is not set."))
		}
	}
}
