package models

import (
	"fmt"
	"bytes"
	"dockulator/db"
	"encoding/json"
	"labix.org/v2/mgo/bson"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	dockerPath  = "/usr/local/bin/docker"
	calcPath    = "/opt/dockulator/calculators/"
	osScriptCmd = "/bin/cat /etc/issue | head -n 1"
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
		"error":    c.Error,
	}})
	return err
}

func (c *Calculation) calculator() string {
	return calcPath + "calc." + c.Language
}

// GetOS will set the OS attribute of the calculation
func (c *Calculation) GetOS() error {
	// example command: `docker run 12345 getos.sh`
	cmd := exec.Command(dockerPath, "run", c.Instance, osScriptCmd)
	out, err := cmd.Output()
	if err != nil {
		c.Error = err.Error()
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
		c.Error = err.Error()
		// TODO: just run the calculation in go
		return err
	}
	floatVal := strings.TrimSpace(string(out))
	answer, err := strconv.ParseFloat(string(floatVal), 64)
	if err != nil {
		c.Error = err.Error()
		// Something definitely went bad.
		return err
	}
	c.Answer = answer
	return nil
}

func (c *Calculation) String() string {
	return fmt.Sprintf("Calculation: %v\nOS: %v\nLanguage: %v\nAnswer: %v\nInstance: %v\nError: %v\n", c.Calculation, c.OS, c.Language, c.Answer, c.Instance, c.Error)
}

// Don't want to modify the actual object
func (c Calculation) AsJson() ([]byte, error) {
	c.Language = GetLanguage(c.Language)
	return json.Marshal(c)
}
