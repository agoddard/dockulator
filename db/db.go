package db

import (
	"fmt"
	"os"
	"labix.org/v2/mgo"
)

func init () {
	mongoUrl = os.Getenv("MONGO_URL")
	if mongoUrl == "" {
		panic(fmt.Errorf("Must supply a MONGO_URL env variable"))
	}
}

var (
	mongoUrl string
)

const (
	Collection = "calculations"
)

func GetSession() *mgo.Session {
	mongoUrl := os.Getenv("MONGO_URL")
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		panic(err)
	}
	return session
}
