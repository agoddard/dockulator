package models

import (
	"dockulator/db"
	"labix.org/v2/mgo/bson"
	"log"
)

type Calculations []Calculation

// Get recent calcuations where error is empty
func GetRecent(n int) Calculations {
	session := db.GetSession()
	defer session.Close()

	c := session.DB("").C(db.Collection)

	result := Calculations{}
	err := c.Find(bson.M{"error": ""}).Sort("-time").Limit(n).All(&result)
	if err != nil {
		log.Printf("Got an error finding recent calculations: %v", err.Error())
	}
	return result
}

// Get all unsolved calculations.
// An unsolved calculation has no instance nor an error
func UnsolvedCalculations() (result Calculations) {
	session := db.GetSession()
	defer session.Close()

	col := session.DB("").C(db.Collection)
	col.Find(bson.M{"instance": "", "error": ""}).All(&result)
	return result
}
