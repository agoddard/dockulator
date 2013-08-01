package main

import (
	"fmt"
	"dockulator/models"
	"dockulator/db"
//	"labix.org/v2/mgo/bson"
)

func main () {
	session := db.GetSession()
	defer session.Close()
	// implicitly uses the db in the connection
	c := session.DB("").C(db.Collection)

	/*
	// Save 10 new calculations
	for i := 0 ; i < 10; i++ {
		calc := models.NewCalculation(fmt.Sprintf("%d + 9", i))
		calc.Save()
	}
	*/

	/*
	// Save 1 new calculation and print it
	calc := models.NewCalculation("9 + 3")
	calc.Save()
	fmt.Println(calc)
	*/

	var err error
	// Get all objects
	var result models.Calculations
	err = c.Find(nil).Sort("-time").Limit(4).All(&result)
	if err != nil {
		panic(err)
	}

	// Print each object
	//fmt.Println(result)

	for i := 0; i < len(result); i++ {
		fmt.Println(string(result[i].Json()))
	}
	/*
	// Print each objects _id
	for i := 0; i < len(result); i++ {
		fmt.Println(result[i].Id.String())
	}
	for i := 0; i < len(result); i++ {
		fmt.Println(result[i].Escaped())
	}
	*/
	/*
	for i := 0; i< len(result); i++ {
		fmt.Println(result[i].Instance)
		result[i].Instance = "HELLO WORLD"
		c.Update(bson.M{"_id": result[i].Id},bson.M{"$set": bson.M{"instance": "value"}})
	}
	*/

	/*
	// Find an object by ID
	result := models.Calculation{}
	c.FindId(bson.ObjectIdHex("51c1bb61022c206ec22aca60")).One(&result)
	fmt.Println(result)
	*/

	/*
	// Get one object and update it
	result := models.Get(0)
	result.OS = "CentOS"
	result.Language = "Python"
	fmt.Println(result)
	*/

}
