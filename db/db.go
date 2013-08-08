package db

import (
	"fmt"
	"os"
	"labix.org/v2/mgo"
)

func init () {
	mongoUrl = os.Getenv("MONGOHQ_URL")
	if mongoUrl == "" {
		panic(fmt.Errorf("Must supply a MONGOHQ_URL env variable"))
	}
}

var (
	mongoUrl string
)

const (
	Complete = "calculations"
	Error = "errors"
	Queue = "queue"
)

func GetSession() *mgo.Session {
	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		panic(err)
	}
	return session
}
