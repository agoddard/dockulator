package models

import (
	"dockulator/db"
	"labix.org/v2/mgo/bson"
	"log"
)

type Calculations []Calculation

func GetRecent(n int) Calculations {
	session := db.GetSession()
	defer session.Close()

	c := session.DB("").C(db.Collection)

	result := Calculations{}
	err := c.Find(nil).Sort("-time").Limit(n).All(&result)
	if err != nil {
		log.Printf("Got an error finding recent calculations: %v", err.Error())
	}
	return result
}

func UnsolvedCalculations() (result Calculations) {
	session := db.GetSession()
	defer session.Close()

	col := session.DB("").C(db.Collection)
	col.Find(bson.M{"instance": ""}).All(&result)
	return result
}
