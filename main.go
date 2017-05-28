package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/cruux/api/config"
	"github.com/cruux/api/handlers"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type Api struct {
	mongo *mgo.Session
}

type Item struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (api *Api) rootHandler(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("{\"message\": \"Welcome on Cruux API\", \"status\": \"OK\"}"))
}

func (api *Api) userPostItemHandler(res http.ResponseWriter, req *http.Request) {
	s := api.mongo.Copy()
	defer s.Close()

	var i Item
	err := json.NewDecoder(req.Body).Decode(&i)
	if err != nil {
		panic(err)
	}

	c := s.DB("store").C("items")
	err = c.Insert(i)
	if err != nil {
		panic(err)
	}
	res.Write([]byte("Wroten"))
}

func (api *Api) userGetItemsHandler(res http.ResponseWriter, req *http.Request) {
	s := api.mongo.Copy()
	defer s.Close()

	var items []Item

	c := s.DB("store").C("items")
	err := c.Find(bson.M{}).All(&items)
	if err != nil {
		panic(err)
	}

	respBody, err := json.Marshal(items)
	if err != nil {
		panic(err)
	}
	res.Write(respBody)
}

func main() {
	api := Api{}

	// Auth
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return config.JwtPublicKey, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	// Mongo
	fmt.Println("Connecting on mongo -> " + config.Env["MONGO_HOST"])
	mongoSession, err := mgo.Dial(config.Env["MONGO_HOST"])
	if err != nil {
		panic(err)
	}
	api.mongo = mongoSession
	defer mongoSession.Close()

	// Router
	globalRouter := mux.NewRouter()
	globalRouter.HandleFunc("/", api.rootHandler)

	publicRouter := mux.NewRouter()
	publicRouter.HandleFunc("/public/login", handlers.LoginHandler).Methods("POST")
	globalRouter.PathPrefix("/public").Handler(negroni.New(
		negroni.NewLogger(),
		negroni.NewRecovery(),
		negroni.Wrap(publicRouter),
	))

	authRouter := mux.NewRouter()
	authRouter.HandleFunc("/user/item", api.userPostItemHandler).Methods("POST")
	authRouter.HandleFunc("/user/item", api.userGetItemsHandler).Methods("GET")
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
