package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/cruux/api/config"
	"github.com/cruux/api/utils"

	"gopkg.in/mgo.v2/bson"
)

type item struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func PostItemHandler(res http.ResponseWriter, req *http.Request) {
	s := config.MongoSession()
	defer s.Close()

	var i item
	err := json.NewDecoder(req.Body).Decode(&i)
	if err != nil {
		panic(err)
	}

	c := s.DB("store").C("items")
	err = c.Insert(i)
	utils.Fatal(err)

	res.Write([]byte("Wroten"))
}

func GetItemsHandler(res http.ResponseWriter, req *http.Request) {
	s := config.MongoSession()
	defer s.Close()

	var items []item

	c := s.DB("store").C("items")
	err := c.Find(bson.M{}).All(&items)
	utils.Fatal(err)

	respBody, err := json.Marshal(items)
	utils.Fatal(err)

	res.Write(respBody)
}
