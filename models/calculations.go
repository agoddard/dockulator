package models

import (
	"dockulator/db"
	"log"
)

type Calculations []Calculation

// Get recent calcuations where error is empty
func GetRecent(n int) Calculations {
	session := db.GetSession()
	defer session.Close()

	c := session.DB("").C(db.Complete)

	result := Calculations{}
	err := c.Find(nil).Sort("-time").Limit(n).All(&result)
	if err != nil {
		log.Printf("Got an error finding recent calculations: %v", err.Error())
	}
	return result
}
