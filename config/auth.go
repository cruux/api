package config

import (
	"crypto/rsa"
	"encoding/base64"

	"github.com/cruux/api/utils"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	JwtPublicKey  *rsa.PublicKey
	JwtPrivateKey *rsa.PrivateKey
)

func initAuth() {
	privateBytes, err := base64.StdEncoding.DecodeString(Env["JWT_PRIVATE_KEY"])
	utils.Fatal(err)

	JwtPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	utils.Fatal(err)

	publicBytes, err := base64.StdEncoding.DecodeString(Env["JWT_PUBLIC_KEY"])
	utils.Fatal(err)

	JwtPublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	utils.Fatal(err)
}
