package models

import (
	"bytes"
	"dockulator/db"
	"encoding/json"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"regexp"
	"strconv"
	"time"
	"log"
)

const (
	dockerPath  = "/usr/local/bin/docker"
	calcPath    = "/opt/dockulator/calculators/"
	calcRe      = `(-)?\d+(\.\d+)? [\+\-\/\*] (-)?\d+(\.\d+)?`
)

var (
	calculationRe = regexp.MustCompile(calcRe)
)

type Calculation struct {
	Calculation string        `json:"calculation"`
	OS          string        `json:"os"`
	Language    string        `json:"language"`
	Id          bson.ObjectId `json:"id" bson:"_id"`
	Answer      float64       `json:"answer"`
	Instance    string        `json:"-"`
	Time        time.Time     `json:"timestamp"`
	Error       string        `json:"-"`
	Processing  bool          `json:"-"`
}

// Return an empty calculation
func NewCalculation(calculation string) *Calculation {
	return &Calculation{
		Calculation: calculation,
		OS:          "",
		Language:    "",
		Answer:      0.0,
		Instance:    "",
		Id:          bson.NewObjectId(),
		Time:        bson.Now(),
		Error:       "",
		Processing:  false,
	}
}

func GetLanguage(lang string) string {
	switch lang {
	case "rb":
		return "Ruby"
	default:
		return "Unknown Language"
	}
}

// Holy fuck. What the fuck happened. I blacked out.
func CleanCalculation(calc string) string {
	noSpaces := bytes.TrimSpace([]byte(calc))
	var clean bytes.Buffer
	foundOp := false
	needDecimal := true
	needNeg := true
	for _, c := range noSpaces {
		// We start without having found the operator
		if !foundOp {
			// We start with needing a negative
			if needNeg {
				if c == '-' {
					clean.WriteByte(c)
					// Now we don't need another negative
					needNeg = false
					continue
				}
			}

			if needDecimal {
				if c == '.' {
					clean.WriteByte(c)
					// We don't need another decimal
					needDecimal = false
					// If we got a decimal we don't need a negative
					needNeg = false
				}
			}

			// We always like numbers
			if c >= '0' && c <= '9' {
				clean.WriteByte(c)
				needNeg = false
			}
			// If we come across an operator, we use it and go to the next thing
			if c == '+' || c == '-' || c == '*' || c == '/' {
				clean.WriteRune(' ')
				clean.WriteByte(c)
				clean.WriteRune(' ')
				foundOp = true
				needDecimal = true
				needNeg = true
			}
		} else {
			if needNeg {
				if c == '-' {
					clean.WriteByte(c)
					needNeg = false
				}
			}

			if needDecimal {
				if c == '.' {
					clean.WriteByte(c)
					needDecimal = false
					needNeg = false
				}
			}
			if c >= '0' && c <= '9' {
				clean.WriteByte(c)
				needNeg = false
			}
		}
	}
	// after we have parsed out the garbage, we just do a quick regexp
	// to make sure we are sane.
	if calculationRe.Match(clean.Bytes()) {
		return clean.String()
	}
	return "error"
}

// Get a calculation from mongo by _id
func Get(id string) (c *Calculation, err error) {
	session := db.GetSession()
	defer session.Close()

	col := session.DB("").C(db.Complete)

	err = col.FindId(bson.ObjectIdHex(id)).One(&c)
	return c, err
}

// Insert a calculation to a specified collection
func (c *Calculation) Insert(collection string) error {
	session := db.GetSession()
	defer session.Close()

	col := session.DB("").C(collection)
	log.Printf("Inserting into %v\n%v\n", collection, c)
	err := col.Insert(c)
	return err
}

// Update the calculation in Mongo
func (c *Calculation) Save() {
	c.Remove(db.Queue)
	if c.Error != "" {
		log.Printf("Inserting %v into Error collection", c)
		c.Insert(db.Error)
		return
	}
	log.Printf("Completed!\n %v\n", c)
	c.Insert(db.Complete)
}

func (c *Calculation) Remove(collection string) error {
	session := db.GetSession()
	defer session.Close()
	col := session.DB("").C(collection)
	log.Printf("Removing this guy from %v:\n%v\n", collection, c)
	return col.RemoveId(c.Id)
}

func (c *Calculation) calculator() string {
	return calcPath + "calc." + c.Language
}

// Get the operating system this command is running on
func (c *Calculation) OperatingSystem() (string, error) {
	cmd := NewCommand(dockerPath, "run", c.Instance, "/bin/cat", "/etc/issue")
	cmd.RegisterCleaner(FirstLine)
	out, err := cmd.Output()
	return string(out), err
}

// Get the answer to the calculation
func (c *Calculation) Calc() (string, error) {
	cmd := NewCommand(dockerPath, "run", c.Instance, c.calculator(), c.Calculation)
	out, err := cmd.Output()
	return string(out), err
}

// GetOS will set the OS attribute of the calculation
func (c *Calculation) GetOS() error {
	os, err := c.OperatingSystem()
	c.OS = os
	return err
}

// Calculate will set the Answer attribute of the calculation
func (c *Calculation) GetCalc() error {
	calculation, err := c.Calc()
	answer, err := strconv.ParseFloat(calculation, 64)
	c.Answer = answer
	if err != nil {
		c.Error = err.Error()
	}
	return err
}

func (c *Calculation) String() string {
	return fmt.Sprintf("Calculation: %v\nOS: %v\nLanguage: %v\nAnswer: %v\nInstance: %v\nError: %v\nProcessing: %v", c.Calculation, c.OS, c.Language, c.Answer, c.Instance, c.Error, c.Processing)
}

// Not a pointer method because we don't want to modify the instance.
// We want to modify a copy.
func (c Calculation) AsJson() ([]byte, error) {
	c.Language = GetLanguage(c.Language)
	return json.Marshal(c)
}

func GetNext() *Calculation {
	session := db.GetSession()
	defer session.Close()
	col := session.DB("").C(db.Queue)

	change := mgo.Change{
		Update: bson.M{"$set": bson.M{"processing": true}},
		ReturnNew: true,
	}
	var result Calculation
	info, err := col.Find(bson.M{"processing": false}).Apply(change, &result)
	log.Printf("ChangeInfo:\n%v\n", info)
	if err != nil {
		log.Printf("Error in GetNext: %v", err)
	}
	log.Printf("Found: %v\n", result)
	// If we found anything
	if result.Processing {
		return &result
	}
	return &Calculation{}
}
