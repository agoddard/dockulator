package calculations

import (
	"dockulator/db"
	"labix.org/v2/mgo/bson"
	"time"
)

type Calculation struct {
	Calculation, OS, Language string
	Id                        bson.ObjectId "_id"
	Answer                    float64
	Instance                  string
	Time                      time.Time
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
