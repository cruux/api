package config

import (
	"crypto/rsa"
	"encoding/base64"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	JwtPublicKey  *rsa.PublicKey
	JwtPrivateKey *rsa.PrivateKey
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initAuth() {
	privateBytes, err := base64.StdEncoding.DecodeString(Env["JWT_PRIVATE_KEY"])
	fatal(err)

	JwtPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	fatal(err)

	publicBytes, err := base64.StdEncoding.DecodeString(Env["JWT_PUBLIC_KEY"])
	fatal(err)

	JwtPublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	fatal(err)
}
