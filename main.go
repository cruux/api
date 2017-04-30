package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
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

	i := Item{Title: "First Item", Content: "And its content"}

	c := s.DB("store").C("items")
	err := c.Insert(i)
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
	// Env
	port := "4242"
	mongo_host := "mongo"

	api := Api{}

	// Mongo
	fmt.Println("Connecting on mongo -> " + mongo_host)
	mongo_session, err := mgo.Dial(mongo_host)
	if err != nil {
		panic(err)
	}
	api.mongo = mongo_session
	defer mongo_session.Close()

	// Router
	r := mux.NewRouter()
	r.HandleFunc("/", api.rootHandler)
	r.HandleFunc("/user/item", api.userPostItemHandler).Methods("POST")
	r.HandleFunc("/user/item", api.userGetItemsHandler).Methods("GET")

	// Http
	fmt.Println("Listening on port " + port)
	http.ListenAndServe(":"+port, r)
}
