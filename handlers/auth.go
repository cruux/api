package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cruux/api/config"
	jwt "github.com/dgrijalva/jwt-go"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

func JsonResponse(response interface{}, w http.ResponseWriter) {
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user UserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	if strings.ToLower(user.Username) != "someone" || user.Password != "p@ssword" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Println("Error logging in")
		fmt.Fprint(w, "Invalid credentials")
		return
	}

	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	tokenString, err := token.SignedString(config.JwtPrivateKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		panic(err)
	}

	response := Token{tokenString}
	JsonResponse(response, w)
}
