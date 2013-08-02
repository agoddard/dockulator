package models

import (
	"dockulator/db"
	"labix.org/v2/mgo/bson"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"bytes"
	"regexp"
)

const (
	dockerPath   = "/usr/local/bin/docker"
	calcPath     = "/opt/dockulator/calculators/"
	osScriptPath = "TODO: JEROME FIXME!!!"
	calcRe = `\d+(\.\d+)? [\+\-\/\*] \d+(\.\d+)?`
)

var (
	calculationRe = regexp.MustCompile(calcRe)
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
	}
}

func CleanCalculation(calc string) string {
	noSpaces := bytes.TrimSpace([]byte(calc))
	var clean bytes.Buffer
	needLeft := true
	needRight := false
	needDecimal := true
	for _, c := range noSpaces {
		if needLeft {
			if needDecimal {
				if c == '.' {
					clean.WriteByte(c)
					needDecimal = false
				}
			}
			// If c is a digit or a decmial keep save it and keep going
			if c >= '0' && c <= '9' {
				clean.WriteByte(c)
			} else if (c == '+' || c == '-' || c == '*' || c == '/') {
				clean.WriteRune(' ')
				clean.WriteByte(c)
				clean.WriteRune(' ')
				needLeft = false
				needRight = true
				needDecimal = true
			}
		}

		if needRight {
			if needDecimal {
				if c == '.' {
					clean.WriteByte(c)
					needDecimal = false
				}
			}
			if c >= '0' && c <= '9' {
				clean.WriteByte(c)
			}
		}
	}
	if calculationRe.Match(clean.Bytes()) {
		return clean.String()
	}
	return "error"
}

// Get a calculation from mongo by _id
func Get(id string) (c *Calculation, err error) {
	session := db.GetSession()
	defer session.Close()

	col := session.DB("").C(db.Collection)

	err = col.FindId(bson.ObjectIdHex(id)).One(&c)
	return c, err
}

// Insert a calculation to mongo
func (c *Calculation) Insert() (err error) {
	session := db.GetSession()
	defer session.Close()

	col := session.DB("").C(db.Collection)

	err = col.Insert(c)
	return err
}

// Update the calculation in Mongo
func (c *Calculation) Save() (err error) {
	session := db.GetSession()
	defer session.Close()

	col := session.DB("").C(db.Collection)
	col.Update(bson.M{"_id": c.Id}, bson.M{"$set": bson.M{
		"instance": c.Instance,
		"answer":   c.Answer,
		"language": c.Language,
		"os":       c.OS,
	}})
	return err
}

func (c *Calculation) calculator() string {
	return calcPath + "calc." + c.Language
}

// GetOS will set the OS attribute of the calculation
func (c *Calculation) GetOS() error {
	// example command: `docker run 12345 getos.sh`
	cmd := exec.Command(dockerPath, "run", c.Instance, osScriptPath)
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	os := strings.TrimSpace(string(out))
	c.OS = os
	return nil
}

// Calculate will set the Answer attribute of the calculation
func (c *Calculation) Calculate() error {
	// example command: `docker run 12345 calc.rb 4 + 2`
	cmd := exec.Command(dockerPath, "run", c.Instance, c.calculator(), c.Calculation)
	out, err := cmd.Output()
	if err != nil {
		// TODO: just run the calculation in go
		return err
	}
	floatVal := strings.TrimSpace(string(out))
	answer, err := strconv.ParseFloat(string(floatVal), 64)
	if err != nil {
		// Something definitely went bad.
		return err
	}
	c.Answer = answer
	return nil
}
