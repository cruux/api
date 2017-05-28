package main

import (
	"fmt"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/cruux/api/config"
	"github.com/cruux/api/handlers"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func rootHandler(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("{\"message\": \"Welcome on Cruux API\", \"status\": \"OK\"}"))
}

func main() {
	// Auth
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return config.JwtPublicKey, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	// Mongo
	config.MongoSession()
	defer config.MongoSession().Close()

	// Router
	globalRouter := mux.NewRouter()
	globalRouter.HandleFunc("/", rootHandler)

	publicRouter := mux.NewRouter()
	publicRouter.HandleFunc("/public/login", handlers.LoginHandler).Methods("POST")
	globalRouter.PathPrefix("/public").Handler(negroni.New(
		negroni.NewLogger(),
		negroni.NewRecovery(),
		negroni.Wrap(publicRouter),
	))

	authRouter := mux.NewRouter()
	authRouter.HandleFunc("/user/item", handlers.PostItemHandler).Methods("POST")
	authRouter.HandleFunc("/user/item", handlers.GetItemsHandler).Methods("GET")
	globalRouter.PathPrefix("/user").Handler(negroni.New(
		negroni.NewLogger(),
		negroni.NewRecovery(),
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(authRouter),
	))

	// Http
	fmt.Println("Listening on port " + config.Env["PORT"])
	http.ListenAndServe(":"+config.Env["PORT"], globalRouter)
}
