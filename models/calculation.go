package models

import (
	"dockulator/db"
	"encoding/json"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

type Calculation struct {
	Calculation string        `json:"calculation"`
	OS          string        `json:"os"`
	Language    string        `json:"language"`
	Id          bson.ObjectId `json:"-" bson:"_id"`
	Answer      float64       `json:"answer"`
	Instance    string        `json:"-"`
	Time        time.Time     `json:"timestamp"`
}

func NewCalculation(calculation string) *Calculation {
	return &Calculation{
		Calculation: calculation,
		OS:          "",
		Language:    "",
		Answer:      0.0,
		Instance:    "",
		Id:          bson.NewObjectId(),
		Time:        bson.Now(),
	}
}

func Get(id string) (c *Calculation, err error) {
	session := db.GetSession()
	defer session.Close()

	col := session.DB("").C(db.Collection)

	err = col.FindId(bson.ObjectIdHex(id)).One(&c)
	return c, err
}

func (c *Calculation) Insert() (err error) {
	session := db.GetSession()
	defer session.Close()

	col := session.DB("").C(db.Collection)

	err = col.Insert(c)
	return err
}

func (c *Calculation) Save() (err error) {
	session := db.GetSession()
	defer session.Close()

	col := session.DB("").C(db.Collection)
	col.Update(bson.M{"_id": c.Id}, bson.M{"$set": bson.M{"instance": c.Instance, "answer": c.Answer}})
	return err
}

func (c *Calculation) Json() string {
	json, err := json.Marshal(c)
	if err != nil {
		log.Printf("Got an error marshaling calculation: %v", err.Error())
		return ""
	}
	return string(json)
}
