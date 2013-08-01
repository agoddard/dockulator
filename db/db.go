package db

import (
	"os"
	"labix.org/v2/mgo"
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
