package config

import (
	"fmt"

	"github.com/cruux/api/utils"
	mgo "gopkg.in/mgo.v2"
)

var (
	mongoSession *mgo.Session
)

func MongoSession() *mgo.Session {
	var err error
	if mongoSession == nil {
		fmt.Println("Connecting on mongo -> " + Env["MONGO_HOST"])
		mongoSession, err = mgo.Dial(Env["MONGO_HOST"])
		utils.Fatal(err)
	}

	return mongoSession.Copy()
}
